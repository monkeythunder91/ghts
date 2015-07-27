package price_data

import (
	공용 "github.com/ghts/ghts/shared"
	
	"math"
	"testing"
	"time"
)

func TestF가격정보_Go루틴_수신_및_캐시(테스트 *testing.T) {
	공용.F메모("증권사 최신 가격정보 입수 기능 개발할 것.")
	
	const 값_수량 = 10
	
	에러 := f가격정보_Go루틴_테스트_초기화()
	공용.F테스트_에러없음(테스트, 에러)
	
	종목_모음 := 공용.F샘플_종목_모음()
	통화값_모음 := 공용.F샘플_통화값_모음(값_수량)
	
	값_모음 := make(map[string]공용.I가격정보)
	r := 공용.F임의값_생성기()
	
	for i:=0 ; i < 값_수량 ; i++ {
		종목코드 := 종목_모음[r.Intn(len(종목_모음))].G코드()
		통화 := 통화값_모음[i]
		가격정보 := 공용.New가격정보(종목코드, 통화, time.Now())
				
		값_모음[종목코드] = 가격정보
		
		// 값 쓰기
		회신 := 공용.New질의(공용.P메시지_SET, 
					종목코드, 통화.G단위(), 통화.G실수값(), 
					가격정보.G시점().Format(공용.P시간_형식)).G회신(Ch가격정보)
		공용.F테스트_에러없음(테스트, 회신.G에러())
		공용.F테스트_같음(테스트, 회신.G구분(), 공용.P메시지_OK)
		공용.F테스트_같음(테스트, 회신.G길이(), 0)
	}
	
	// 값 읽기
	현재_시점 := time.Now()
	
	for 종목코드, 가격정보 := range 값_모음 {
		통화단위 := 가격정보.G가격().G단위()
		금액 := 가격정보.G가격().G실수값()
		
		회신 := 공용.New질의(공용.P메시지_GET, 종목코드, 1).G회신(Ch가격정보)
		공용.F테스트_에러없음(테스트, 회신.G에러())
		공용.F테스트_같음(테스트, 회신.G구분(), 공용.P메시지_OK)
		공용.F테스트_같음(테스트, 회신.G길이(), 3)
		공용.F테스트_같음(테스트, 회신.G내용(0), 통화단위)
		공용.F테스트_같음(테스트, 회신.G내용(1), 공용.F2문자열(금액))
	
		_, 에러 := time.Parse(공용.P시간_형식, 회신.G내용(2))
		공용.F테스트_에러없음(테스트, 에러)
		차이 := 현재_시점.Sub(time.Now()).Seconds()
		공용.F테스트_참임(테스트, 차이 < 5.0)
	}
	
	// 값 수정
	수정된_값_모음 := make(map[string]공용.I가격정보)
	통화값_모음 = 공용.F샘플_통화값_모음(len(값_모음))
	
	i := 0
	
	for 종목코드, _ := range 값_모음 {
		가격정보 := 공용.New가격정보(종목코드, 통화값_모음[i], time.Now())
		수정된_값_모음[종목코드] = 가격정보
		i++
		
		// 값 쓰기
		회신 := 공용.New질의(공용.P메시지_SET, 
					종목코드, 가격정보.G가격().G단위(), 가격정보.G가격().G실수값(), 
					가격정보.G시점().Format(공용.P시간_형식)).G회신(Ch가격정보)
		공용.F테스트_에러없음(테스트, 회신.G에러())
		공용.F테스트_같음(테스트, 회신.G구분(), 공용.P메시지_OK)
		공용.F테스트_같음(테스트, 회신.G길이(), 0)
	}
	
	현재_시점 = time.Now()
	
	// 수정된 값 읽기
	for 종목코드, 가격정보 := range 수정된_값_모음 {
		통화단위 := 가격정보.G가격().G단위()
		금액 := 가격정보.G가격().G실수값()
		
		회신 := 공용.New질의(공용.P메시지_GET, 종목코드, 1).G회신(Ch가격정보)
		공용.F테스트_에러없음(테스트, 회신.G에러())
		공용.F테스트_같음(테스트, 회신.G구분(), 공용.P메시지_OK)
		공용.F테스트_같음(테스트, 회신.G길이(), 3)
		공용.F테스트_같음(테스트, 회신.G내용(0), 통화단위)
		공용.F테스트_같음(테스트, 회신.G내용(1), 공용.F2문자열(금액))
	
		시점, 에러 := time.Parse(공용.P시간_형식, 회신.G내용(2))
		공용.F테스트_에러없음(테스트, 에러)
		차이 := 현재_시점.Sub(시점).Seconds()
		공용.F테스트_참임(테스트, 차이 < 5.0)
	}
}

func TestF가격정보_Go루틴_구독채널_등록(테스트 *testing.T) {
	const 구독채널_수량 = 10
	const 가격정보_수량 = 20
	
	f가격정보_Go루틴_테스트_초기화()
	
	구독채널_모음 := make([]chan 공용.I가격정보, 구독채널_수량)
	
	for i:=0 ; i < 구독채널_수량 ; i++ {
		구독채널_모음[i] = make(chan 공용.I가격정보, 가격정보_수량 + 10)
		Ch가격정보_구독채널_등록 <- 구독채널_모음[i]
	}
	
	종목정보_모음 := 공용.F샘플_종목_모음()
	통화단위_모음 := 공용.F샘플_통화단위_모음()
	r := 공용.F임의값_생성기()
	가격정보_모음 := make([]공용.I가격정보, 가격정보_수량)
	
	for i:=0 ; i < 가격정보_수량 ; i++ {
		종목코드 := 종목정보_모음[r.Intn(len(종목정보_모음))].G코드()
		통화단위 := 통화단위_모음[r.Intn(len(통화단위_모음))]
		금액 := math.Trunc(r.Float64() * math.Pow10(r.Intn(10)) * 100) / 100
		가격정보 := 공용.New가격정보(종목코드, 공용.New통화(통화단위, 금액), time.Now())
		가격정보_모음[i] = 가격정보
		
		회신 := 공용.New질의(공용.P메시지_SET, 종목코드, 
					통화단위, 금액, time.Now()).G회신(Ch가격정보)
		공용.F테스트_에러없음(테스트, 회신.G에러())
		공용.F테스트_같음(테스트, 회신.G길이(), 0)
	}
	
	현재_시점 := time.Now()
	
	for _, 구독채널 := range 구독채널_모음 {
		for i:=0 ; i < 가격정보_수량 ; i++ {
			가격정보 := <-구독채널
			
			공용.F테스트_같음(테스트, 가격정보.G종목코드(), 가격정보_모음[i].G종목코드())
			공용.F테스트_같음(테스트, 가격정보.G가격().G비교(가격정보_모음[i].G가격()), 공용.P같음)
			공용.F테스트_참임(테스트, 현재_시점.Sub(가격정보.G시점()).Seconds() < 10.0)
		}
	}
}

func f가격정보_Go루틴_테스트_초기화() error {
	F가격정보_모듈_실행()
	회신 := 공용.New질의(공용.P메시지_초기화).G회신(ch제어_가격정보_Go루틴)
	
	return 회신.G에러()
}