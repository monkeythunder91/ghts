/* Copyright(C) 2015-2020년 김운하(UnHa Kim)  < unha.kim.ghts at gmail dot com >

키움증권 API는 OCX규격으로 작성되어 있는 데,
OCX규격은 Go언어로 직접 사용하기에 기술적 난이도가 높아서,
손쉽게 다룰 수 있게 도와주는 Qt 라이브러리의 오픈소스 버전을 사용하였습니다.

Qt 라이브러리의 오픈소스 버전의 경우
GHTS의 대부분에서 사용하고 있는 GNU LGPL v 2.1보다
좀 더 강력하고 엄격한 소스코드 공개 의무가 있는 GNU GPL v2를 사용해야 합니다.

이는 개발 난이도 경감을 위한 개발자의 필요에 의한 것이며,
사용자에게 GPL v2의 소스코드 공개 의무를 강제하려는 의도는 아닙니다.

키움증권 API 호출 모듈에 적용된 GPL v2이 LGPL v2보다 더 엄격하긴 합니다만,
키움증권 API 호출 모듈을 애초 의도된 사용법대로 '소켓을 통해서 호출'하여 사용하는 경우에는
GPL v2에서 규정하는 '하나의 단일 소프트웨어' 규정에 포함되지 않기에
사용자가 작성한 소스코드는 GPL v2의 소스코드 공개 의무가 적용되지 않습니다.

다만, 키움증권 API 호출 모듈 그 자체를 수정하거나 타인에게 배포할 경우,
GPL v2 규정에 따른 소스코드 공개 의무가 발생할 수 있습니다.

---------------------------------------------------------

이 프로그램은 자유 소프트웨어입니다.
소프트웨어의 피양도자는 자유 소프트웨어 재단이 공표한 GNU GPL v2
규정에 따라 프로그램을 개작하거나 재배포할 수 있습니다.

이 프로그램은 유용하게 사용될 수 있으리라는 희망에서 배포되고 있지만,
특정한 목적에 적합하다거나, 이익을 안겨줄 수 있다는 묵시적인 보증을 포함한
어떠한 형태의 보증도 제공하지 않습니다.
보다 자세한 사항에 대해서는 GNU GPL v2를 참고하시기 바랍니다.
GNU GPL v2는 이 프로그램과 함께 제공됩니다.

만약, 이 문서가 누락되어 있다면 자유 소프트웨어 재단으로 문의하시기 바랍니다.
(자유 소프트웨어 재단 : Free Software Foundation, Inc.,
59 Temple Place - Suite 330, Boston, MA 02111-1307, USA) */

package kt

import "github.com/ghts/ghts/lib"

const (
	ProgID     = "KHOPENAPI.KHOpenAPICtrl.1"
	CLSID문자열   = "{A1574A0D-6BFA-4BD7-9020-DED88711818D}"
	IID문자열_메소드 = "{CF20FBB6-EDD4-4BE5-A473-FEF91977DEB6}"
	IID문자열_이벤트 = "{7335F12D-8973-4BD5-B7F0-12DF03D175B7}"
)

// TR 및 응답 종류
const (
	TR조회 lib.TR구분 = iota
	TR주문
	TR실시간_정보_구독
	TR실시간_정보_해지
	TR실시간_정보_일괄_해지
	TR접속
	TR접속_상태
	TR접속_해제
	TR초기화
	TR종료

	// Kiwoom API에서 사용되는 것들
	TR종목코드_리스트
	TR로그인_정보
	TR소켓_테스트
)

func TR구분_String(v lib.TR구분) string {
	switch v {
	case TR조회:
		return "TR조회"
	case TR주문:
		return "TR주문"
	case TR실시간_정보_구독:
		return "TR실시간_정보_구독"
	case TR실시간_정보_해지:
		return "TR실시간_정보_해지"
	case TR실시간_정보_일괄_해지:
		return "TR실시간_정보_일괄_해지"
	case TR접속:
		return "TR접속"
	case TR접속_상태:
		return "TR접속_상태"
	case TR접속_해제:
		return "TR접속_해제"
	case TR초기화:
		return "TR초기화"
	case TR종료:
		return "TR종료"
	case TR종목코드_리스트:
		return "TR종목코드_리스트"
	case TR로그인_정보:
		return "TR로그인_정보"
	case TR소켓_테스트:
		return "신호"
	default:
		return lib.F2문자열("예상하지 못한 M값 : '%v'", v)
	}
}

type T로그인_정보_구분 uint8

const (
	P전체_계좌_수량 T로그인_정보_구분 = iota
	P전체_계좌_번호
	P사용자_ID
	P사용자_이름
	P키보드_보안_상태 // 0:정상, 1:해지
	P방화벽_상태    // 0:미설정, 1:설정, 2:해지
)

func (t T로그인_정보_구분) String() string {
	switch t {
	case P전체_계좌_수량:
		return "ACCOUNT_CNT"
	case P전체_계좌_번호:
		return "ACCNO"
	case P사용자_ID:
		return "USER_ID"
	case P사용자_이름:
		return "USER_NAME"
	case P키보드_보안_상태:
		return "KEY_BSECGB"
	case P방화벽_상태:
		return "FIREW_SECGB:"
	}

	return lib.F2문자열("예상하지 못한 로그인 정보 구분값 : '%v'", int(t))
}

type T방화벽_상태 uint8

const (
	P방화벽_미설정 T방화벽_상태 = iota
	P방화벽_설정
	P방화벽_해지
)

func (t T방화벽_상태) String() string {
	switch t {
	case P방화벽_미설정:
		return "미설정"
	case P방화벽_설정:
		return "설정"
	case P방화벽_해지:
		return "해지"
	}

	return lib.F2문자열("예상하지 못한 방화벽 상태 구분값 : '%v'", int(t))
}
