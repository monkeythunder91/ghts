/* This file is part of GHTS.

GHTS is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

GHTS is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with GHTS.  If not, see <http://www.gnu.org/licenses/>.

@author: UnHa Kim <unha.kim.ghts@gmail.com> */

package shared

import (
	dec "github.com/wayn3h0/go-decimal"

	"bytes"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"
)

type S비어있는_구조체 struct{}

// 안전한 bool
type s안전한_bool struct {
	sync.RWMutex
	값 bool
}

func (this *s안전한_bool) G값() bool {
	this.RLock() // Go언어의 Embedded Lock
	defer this.RUnlock()

	return this.값
}

func (this *s안전한_bool) S값(값 bool) error {
	this.Lock()
	defer this.Unlock()

	if this.값 == 값 {
		return F에러_생성("이미 %v임.", 값)
	} else {
		this.값 = 값
		return nil
	}
}

// 기본 메시지
type s기본_메시지 struct {
	구분 string
	내용 []string
}

func (this s기본_메시지) G구분() string {
	return this.구분
}

func (this s기본_메시지) G내용(인덱스 int) string {
	if 인덱스 >= len(this.내용) {
		F에러_출력("인덱스 입력값은 '길이'보다 작아야 함 : 길이 %v, 입력값 %v", len(this.내용), 인덱스)
		panic("무효한 인덱스")
	}

	return this.내용[인덱스]
}

func (this s기본_메시지) G내용_전체() []string {
	return this.내용
}

func (this s기본_메시지) G길이() int {
	return len(this.내용)
}

func (this s기본_메시지) String() string {
	var 버퍼 bytes.Buffer

	버퍼.WriteString("구분 : " + this.구분 + "\n")

	if len(this.내용) == 0 {
		버퍼.WriteString("내용 없음. len(내용) == 0. \n")
	} else {
		버퍼.WriteString("내용\n")

		for i := 0; i < len(this.내용); i++ {
			버퍼.WriteString(strconv.Itoa(i) + " : " + this.내용[i] + "\n")
		}
	}

	return 버퍼.String()
}

// 질의 메시지
type s질의_메시지 struct {
	s기본_메시지 // Go언어 구조체 embedding(임베딩) 기능. 상속 비스무리함.
	회신_채널   chan I회신
}

func (this s질의_메시지) G회신_채널() chan I회신 {
	return this.회신_채널
}

func (this s질의_메시지) G검사(타이틀 string, 질의_길이 int) error {
	if this.G구분() == P메시지_GET &&
		this.G길이() == 질의_길이 {
		return nil
	}

	에러 := F에러_생성("잘못된 %s 질의 메시지. 구분 '%v', 길이 %v, 내용 '%v'",
		타이틀, this.G구분(), this.G길이(), this.G내용_전체())

	this.G회신_채널() <- New회신(에러, P메시지_에러)

	return 에러
}

// 회신 메시지
type s회신_메시지 struct {
	s기본_메시지 // Go언어 구조체 embedding(임베딩)
	에러      error
}

func (this s회신_메시지) G에러() error {
	return this.에러
}

// 종목
type s종목 struct {
	코드 string
	이름 string
}

func (this *s종목) G코드() string {
	return this.코드
}

func (this *s종목) G이름() string {
	return this.이름
}

func (this *s종목) String() string {
	return this.코드 + " " + this.이름
}

// 통화
type s통화 struct {
	단위   T통화단위
	금액   *dec.Decimal
	변경불가 bool
}

func (this *s통화) G단위() T통화단위    { return this.단위 }
func (this *s통화) G실수값() float64 { return this.금액.Float() }
func (this *s통화) G정밀값() *dec.Decimal {
	// 참조형이므로 그대로 주지 않고, 복사본을 준다.
	if this.금액 == nil {
		return nil
	}

	정밀값, _ := dec.Parse(this.금액.String())

	return 정밀값
}
func (this *s통화) G실수_문자열(소숫점_이하_자릿수 int) string {
	return this.금액.FloatString(소숫점_이하_자릿수)
}

func (this *s통화) G비교(다른_통화 I통화) T비교결과 {
	switch {
	case this.단위 != 다른_통화.G단위():
		return P비교불가
	default:
		return T비교결과(this.금액.Cmp(다른_통화.G정밀값()))
	}
}

func (this *s통화) G부호() T부호 {
	return T부호(this.금액.Sign())
}

func (this *s통화) G복사본() I통화 {
	s := new(s통화)
	s.단위 = this.G단위()
	s.금액 = this.G정밀값()
	s.변경불가 = false

	return s
}

func (this *s통화) G변경불가() bool {
	return this.변경불가
}

func (this *s통화) S동결() {
	this.변경불가 = true
}

func (this *s통화) S더하기(다른_통화 I통화) I통화 {
	if this.변경불가 {
		panic("변경불가능한 값입니다.")
	}

	다른_통화_금액 := 다른_통화.G정밀값()

	if this.단위 != 다른_통화.G단위() ||
		this.금액 == nil ||
		다른_통화_금액 == nil {
		this.금액 = nil
	} else {
		this.금액 = this.금액.Add(다른_통화_금액)
	}

	return this
}

func (this *s통화) S빼기(다른_통화 I통화) I통화 {
	if this.변경불가 {
		panic("변경불가능한 값입니다.")
	}

	다른_통화_금액 := 다른_통화.G정밀값()

	if this.단위 != 다른_통화.G단위() ||
		this.금액 == nil ||
		다른_통화_금액 == nil {
		this.금액 = nil
	} else {
		this.금액 = this.금액.Sub(다른_통화_금액)
	}

	return this
}

func (this *s통화) S곱하기(다른_통화 I통화) I통화 {
	if this.변경불가 {
		panic("변경불가능한 값입니다.")
	}

	다른_통화_금액 := 다른_통화.G정밀값()

	if this.단위 != 다른_통화.G단위() ||
		this.금액 == nil ||
		다른_통화_금액 == nil {
		this.금액 = nil
	} else {
		this.금액 = this.금액.Mul(다른_통화_금액)
	}

	return this
}

func (this *s통화) S나누기(다른_통화 I통화) I통화 {
	if this.변경불가 {
		panic("변경불가능한 값입니다.")
	}

	다른_통화_금액 := 다른_통화.G정밀값()

	if this.단위 != 다른_통화.G단위() ||
		this.금액 == nil ||
		다른_통화_금액 == nil {
		this.금액 = nil

		return this
	}

	분자, 변환성공1 := new(big.Rat).SetString(this.G정밀값().String())
	분모, 변환성공2 := new(big.Rat).SetString(다른_통화_금액.String())

	// 변환에 실패하거나, 분모가 0이 되면 안 됨.
	if !변환성공1 || !변환성공2 || 분모.Cmp(big.NewRat(0, 1)) == 0 {
		this.금액 = nil

		return this
	}

	결과값 := new(big.Rat).Quo(분자, 분모)

	// 소숫점 이하 1000자리 정도면 충분히 정밀하지 않을까?
	문자열 := 결과값.FloatString(1000)

	for strings.HasSuffix(문자열, "0") {
		문자열 = strings.TrimSuffix(문자열, "0")
	}

	if strings.HasSuffix(문자열, ".") {
		문자열 = strings.TrimSuffix(문자열, ".")
	}

	this.금액, _ = dec.Parse(문자열)

	return this
}

func (this *s통화) S금액(금액 string) I통화 {
	if this.변경불가 {
		panic("변경불가능한 값입니다.")
	}

	정밀값, 에러 := dec.Parse(금액)

	if 에러 != nil {
		this.금액 = nil
	} else {
		this.금액 = 정밀값
	}

	return this
}

func (this *s통화) String() string {
	return string(this.단위) + " " + this.금액.String()
}

// 가격정보
type s가격정보 struct {
	종목 I종목
	가격 I통화
	시점 time.Time
}

func (this *s가격정보) G종목() I종목       { return this.종목 }
func (this *s가격정보) G가격() I통화       { return this.가격.G복사본() }
func (this *s가격정보) G시점() time.Time { return this.시점 }
