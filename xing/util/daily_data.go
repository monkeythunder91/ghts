/* Copyright (C) 2015-2020 김운하(UnHa Kim)  < unha.kim.ghts at gmail dot com >

이 파일은 GHTS의 일부입니다.

이 프로그램은 자유 소프트웨어입니다.
소프트웨어의 피양도자는 자유 소프트웨어 재단이 공표한 GNU LGPL 2.1판
규정에 따라 프로그램을 개작하거나 재배포할 수 있습니다.

이 프로그램은 유용하게 사용될 수 있으리라는 희망에서 배포되고 있지만,
특정한 목적에 적합하다거나, 이익을 안겨줄 수 있다는 묵시적인 보증을 포함한
어떠한 형태의 보증도 제공하지 않습니다.
보다 자세한 사항에 대해서는 GNU LGPL 2.1판을 참고하시기 바랍니다.
GNU LGPL 2.1판은 이 프로그램과 함께 제공됩니다.
만약, 이 문서가 누락되어 있다면 자유 소프트웨어 재단으로 문의하시기 바랍니다.
(자유 소프트웨어 재단 : Free Software Foundation, Inc.,
59 Temple Place - Suite 330, Boston, MA 02111-1307, USA)

Copyright (C) 2015-2020년 UnHa Kim (< unha.kim.ghts at gmail dot com >)

This file is part of GHTS.

GHTS is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, version 2.1 of the License.

GHTS is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with GHTS.  If not, see <http://www.gnu.org/licenses/>. */

package util

import (
	"github.com/ghts/ghts/lib"
	xt "github.com/ghts/ghts/xing/base"
	xing "github.com/ghts/ghts/xing/go"
	bolt "go.etcd.io/bbolt"
	"time"
)

const P일일_가격정보_버킷 = "daily_price_data"

func DB(db파일경로 string) (db *bolt.DB, 에러 error) {
	if db, 에러 = bolt.Open(db파일경로, 0600, nil); 에러 == nil {
		// 일일 가격정보 버킷
		db.Update(func(tx *bolt.Tx) error {
			_, 에러 := tx.CreateBucketIfNotExists([]byte(P일일_가격정보_버킷))

			return 에러
		})
	}

	return
}

func F종목별_일일_가격정보_읽기(종목코드 string, db *bolt.DB) (
	종목별_일일_가격정보_모음 *lib.S종목별_일일_가격정보_모음, 에러 error) {
	에러 = db.View(func(tx *bolt.Tx) error {
		일일_가격정보_버킷 := tx.Bucket([]byte(P일일_가격정보_버킷))
		바이트_모음 := 일일_가격정보_버킷.Get([]byte(종목코드))

		if 에러 := lib.F디코딩(lib.MsgPack, 바이트_모음, &종목별_일일_가격정보_모음); 에러 != nil {
			종목별_일일_가격정보_모음 = nil
			return 에러
		}

		return nil
	})

	종목별_일일_가격정보_모음.S정렬_및_인덱스_설정()

	return
}

func F종목별_일일_가격정보_저장(값 *lib.S종목별_일일_가격정보_모음, db *bolt.DB) error {
	바이트_모음, 에러 := lib.F인코딩(lib.MsgPack, 값)
	if 에러 != nil {
		return 에러
	}

	return db.Update(func(tx *bolt.Tx) error {
		일일_가격정보_버킷 := tx.Bucket([]byte(P일일_가격정보_버킷))

		return 일일_가격정보_버킷.Put([]byte(값.G종목코드()), 바이트_모음)
	})
}

func F일일_가격정보_수집(db파일경로 string) (에러 error) {
	defer xing.F리소스_정리()
	xing.F초기화(xt.P서버_모의투자)

	db, 에러 := DB(db파일경로)
	lib.F확인(에러)
	defer db.Close()

	var 시작일, 종료일, 마지막_저장일 time.Time
	var 종목별_일일_가격정보_모음 *lib.S종목별_일일_가격정보_모음

	for i, 종목코드 := range xing.F종목코드_모음_전체() {
		종목별_일일_가격정보_모음, 에러 = F종목별_일일_가격정보_읽기(종목코드, db)

		if 에러 == nil {
			// lib.S종목별_일일_가격정보_모음 는 일자 순서로 정렬되어 있음.
			마지막_저장일 = 종목별_일일_가격정보_모음.M저장소[len(종목별_일일_가격정보_모음.M저장소)-1].G일자2()

			시작일 = 마지막_저장일.AddDate(0, 0, 1)
		} else {
			시작일 = lib.F지금().AddDate(-30, 0, 0)
		}

		if 시작일.Equal(xing.F당일()) || 시작일.After(xing.F당일()) {
			lib.F체크포인트(i, 종목코드, "이미 최신 데이터로 업데이트 되어 있음.")
			continue // 이미 최신 데이터로 업데이트 되어 있음.
		}

		if lib.F금일().After(xing.F당일()) {
			종료일 = lib.F금일()
		} else if lib.F금일().Equal(xing.F당일()) &&
			lib.F지금().Hour() >= 16 {
			종료일 = lib.F금일()
		} else {
			종료일 = xing.F전일()
		}

		값_모음, 에러 := xing.TrT8413_현물_차트_일주월(종목코드, 시작일, 종료일, xt.P일주월_일)
		if 에러 != nil {
			lib.F에러_출력(에러)
			continue
		} else if len(값_모음) == 0 {
			lib.F체크포인트(i, 종목코드, "추가 저장할 데이터가 없음.")
			continue // 추가 저장할 데이터가 없음.
		}

		일일_가격정보_슬라이스 := make([]*lib.S일일_가격정보, len(값_모음))

		for i, 일일_데이터 := range 값_모음 {
			일일_가격정보_슬라이스[i] = lib.New일일_가격정보(
				일일_데이터.M종목코드,
				일일_데이터.M일자,
				일일_데이터.M시가,
				일일_데이터.M고가,
				일일_데이터.M저가,
				일일_데이터.M종가,
				일일_데이터.M거래량)
		}

		lib.F체크포인트(i, 종목코드, 시작일, 종료일, len(값_모음), len(일일_가격정보_슬라이스))

		if 종목별_일일_가격정보_모음 != nil && len(종목별_일일_가격정보_모음.M저장소) > 0 {
			일일_가격정보_슬라이스 = append(일일_가격정보_슬라이스, 종목별_일일_가격정보_모음.M저장소...)
		}

		종목별_일일_가격정보_모음, 에러 = lib.New종목별_일일_가격정보_모음(일일_가격정보_슬라이스)
		if 에러 != nil {
			lib.F에러_출력(에러)
			continue
		}

		에러 = F종목별_일일_가격정보_저장(종목별_일일_가격정보_모음, db)
		if 에러 != nil {
			lib.F에러_출력(에러)
			continue
		}
	}

	return nil
}
