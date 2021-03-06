package util

import (
	"github.com/ghts/ghts/lib"
	"sort"
)

func New영업일_모음(db파일경로 string) *S영업일_모음 {
	db, 에러 := DB(db파일경로)
	defer db.Close()
	lib.F확인(에러)

	일일_가격정보_모음_KODEX200, 에러 := F종목별_일일_가격정보_읽기("069500", db)
	lib.F확인(에러)

	일일_가격정보_모음_삼성전자, 에러 := F종목별_일일_가격정보_읽기("005930", db)
	lib.F확인(에러)

	개장일_맵 := make(map[uint32]int)

	for _, 일일_정보 := range 일일_가격정보_모음_KODEX200.M저장소 {
		개장일_맵[일일_정보.M일자] = -1
	}

	for _, 일일_정보 := range 일일_가격정보_모음_삼성전자.M저장소 {
		개장일_맵[일일_정보.M일자] = -1
	}

	개장일_모음 := make([]int, len(개장일_맵))

	i := 0
	for 개장일, _ := range 개장일_맵 {
		개장일_모음[i] = int(개장일)
		i++
	}

	sort.Ints(개장일_모음)

	저장소 := make([]uint32, len(개장일_맵))

	for i, 개장일 := range 개장일_모음 {
		저장소[i] = uint32(개장일)
		개장일_맵[uint32(개장일)] = i
	}

	s := new(S영업일_모음)
	s.M저장소 = 저장소
	s.인덱스_맵 = 개장일_맵

	return s
}

type S영업일_모음 struct {
	M저장소  []uint32
	인덱스_맵 map[uint32]int
}

func (s S영업일_모음) G인덱스(일자 uint32) int {
	if 인덱스, 존재함 := s.인덱스_맵[일자]; 존재함 {
		return 인덱스
	} else {
		return -1
	}
}

func (s S영업일_모음) G증분_개장일(일자 uint32, 증분 int) (uint32, error) {
	if 인덱스 := s.G인덱스(일자); 인덱스 < 0 {
		return 0, lib.New에러("존재하지 않는 일자 : '%v'", 일자)
	} else if 인덱스+증분 < 0 || 인덱스+증분 >= len(s.M저장소) {
		return 0, lib.New에러("범위를 벗어난 증분 : '%v' '%v'", 인덱스+증분, len(s.M저장소))
	} else {
		return s.M저장소[인덱스+증분], nil
	}
}

