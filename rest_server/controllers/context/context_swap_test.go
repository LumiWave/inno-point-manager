package context

import (
	"testing"

	"github.com/LumiWave/inno-point-manager/rest_server/controllers/auth"
)

func p2c(pointAdj int64, coinAdj float64) *ReqSwapInfo {
	return &ReqSwapInfo{
		AUID:          1,
		TxType:        EventID_P2C,
		SwapFromPoint: SwapPoint{AdjustPointQuantity: pointAdj},
		SwapToCoin:    SwapCoin{AdjustCoinQuantity: coinAdj},
	}
}

func c2p(coinAdj float64, pointAdj int64) *ReqSwapInfo {
	return &ReqSwapInfo{
		AUID:         1,
		TxType:       EventID_C2P,
		SwapFromCoin: SwapCoin{AdjustCoinQuantity: coinAdj},
		SwapToPoint:  SwapPoint{AdjustPointQuantity: pointAdj},
	}
}

func TestReqSwapInfoCheckValidate(t *testing.T) {
	cases := []struct {
		name string
		req  *ReqSwapInfo
		ok   bool
	}{
		{"포인트->코인 정상", p2c(-100, 1.5), true},
		{"포인트->코인 포인트가 양수 (포인트가 늘어나는 요청)", p2c(100, 1.5), false},
		{"포인트->코인 포인트 0", p2c(0, 1.5), false},
		{"포인트->코인 코인이 음수", p2c(-100, -1.5), false},
		{"코인->포인트 정상", c2p(-1.5, 100), true},
		{"코인->포인트 코인이 양수", c2p(1.5, 100), false},
		{"코인->포인트 포인트가 음수", c2p(-1.5, -100), false},
		{"프리세일 코인->포인트 정상", func() *ReqSwapInfo { r := c2p(-1.5, 100); r.TxType = EventID_Server_C2P; return r }(), true},
		{"프리세일 코인->포인트 부호 오류", func() *ReqSwapInfo { r := c2p(1.5, 100); r.TxType = EventID_Server_C2P; return r }(), false},
		{"알 수 없는 타입", &ReqSwapInfo{AUID: 1, TxType: 99}, false},
		{"타입 없음", &ReqSwapInfo{AUID: 1}, false},
		{"계정 없음", func() *ReqSwapInfo { r := p2c(-100, 1.5); r.AUID = 0; return r }(), false},
	}
	for _, c := range cases {
		resp := c.req.CheckValidate(nil)
		if (resp == nil) != c.ok {
			t.Errorf("%s : 통과 여부가 기대와 다름 (기대 %v, 응답 %v)", c.name, c.ok, resp)
		}
	}
}

func TestReqSwapStatusCheckValidate(t *testing.T) {
	allowed := map[int64]bool{
		SWAP_status_fee_transfer_start:           true,
		SWAP_status_fee_transfer_success:         true,
		SWAP_status_fee_transfer_fail:            true,
		SWAP_status_token_transfer_deposit_start: true,
		SWAP_status_token_transfer_deposit_fail:  true,
	}
	for status := int64(-1); status <= 12; status++ {
		req := &ReqSwapStatus{TxStatus: status, FromWalletAddress: "0xabc"}
		resp := req.CheckValidate(nil)
		if (resp == nil) != allowed[status] {
			t.Errorf("상태 %d : 통과 여부가 기대와 다름 (기대 %v)", status, allowed[status])
		}
	}

	// 지갑 주소가 없으면 거절한다.
	if resp := (&ReqSwapStatus{TxStatus: SWAP_status_fee_transfer_start}).CheckValidate(nil); resp == nil {
		t.Error("지갑 주소 없는 요청이 통과함")
	}
}

func TestCanCancelByUser(t *testing.T) {
	cases := []struct {
		name string
		info *ReqSwapInfo
		ok   bool
	}{
		{"포인트->코인 시작 직후", &ReqSwapInfo{TxType: EventID_P2C, TxStatus: SWAP_status_init}, true},
		{"포인트->코인 수수료 전송 중", &ReqSwapInfo{TxType: EventID_P2C, TxStatus: SWAP_status_fee_transfer_start}, true},
		{"포인트->코인 수수료 확인됨", &ReqSwapInfo{TxType: EventID_P2C, TxStatus: SWAP_status_fee_transfer_success}, false},
		{"포인트->코인 수수료 확인 플래그", &ReqSwapInfo{TxType: EventID_P2C, TxStatus: SWAP_status_fee_transfer_start, IsFeeComplete: true}, false},
		{"포인트->코인 출금 시작", &ReqSwapInfo{TxType: EventID_P2C, TxStatus: SWAP_status_token_transfer_withdrawal_start, IsFeeComplete: true}, false},
		{"포인트->코인 출금 해시가 있음", &ReqSwapInfo{TxType: EventID_P2C, TxStatus: SWAP_status_init, SwapToCoin: SwapCoin{TokenTxHash: "0x1"}}, false},
		{"포인트->코인 출금 성공", &ReqSwapInfo{TxType: EventID_P2C, TxStatus: SWAP_status_token_transfer_withdrawal_success}, false},
		{"코인->포인트 시작 직후", &ReqSwapInfo{TxType: EventID_C2P, TxStatus: SWAP_status_init}, true},
		{"코인->포인트 전송 중", &ReqSwapInfo{TxType: EventID_C2P, TxStatus: SWAP_status_token_transfer_deposit_start}, true},
		{"코인->포인트 입금 확인됨", &ReqSwapInfo{TxType: EventID_C2P, TxStatus: SWAP_status_token_transfer_deposit_success}, false},
		{"코인->코인 수수료만 확인됨", &ReqSwapInfo{TxType: EventID_C2C, TxStatus: SWAP_status_fee_transfer_success, IsFeeComplete: true}, true},
		{"코인->코인 둘 다 확인됨", &ReqSwapInfo{TxType: EventID_C2C, TxStatus: SWAP_status_token_transfer_deposit_success, IsFeeComplete: true, SwapFromCoin: SwapCoin{IsComplete: true}}, false},
		{"알 수 없는 타입", &ReqSwapInfo{TxType: 99, TxStatus: SWAP_status_init}, false},
	}
	for _, c := range cases {
		if got := c.info.CanCancelByUser(); got != c.ok {
			t.Errorf("%s : 취소 가능 여부 %v, 기대 %v", c.name, got, c.ok)
		}
	}
}

func TestCanReportDeposit(t *testing.T) {
	if (&ReqSwapInfo{TxType: EventID_P2C, TxStatus: SWAP_status_init}).CanReportDeposit() {
		t.Error("포인트->코인 스왑에 코인 전송 시작 신고가 통과함")
	}
	if !(&ReqSwapInfo{TxType: EventID_C2P, TxStatus: SWAP_status_init}).CanReportDeposit() {
		t.Error("코인->포인트 스왑의 코인 전송 시작 신고가 거절됨")
	}
	if (&ReqSwapInfo{TxType: EventID_C2P, TxStatus: SWAP_status_token_transfer_deposit_success}).CanReportDeposit() {
		t.Error("입금이 확인된 스왑에 코인 전송 시작 신고가 통과함")
	}
}

func TestCanPayout(t *testing.T) {
	cases := []struct {
		name string
		info *ReqSwapInfo
		ok   bool
	}{
		{"수수료 전송 중", &ReqSwapInfo{TxType: EventID_P2C, TxStatus: SWAP_status_fee_transfer_start}, true},
		{"콜백이 먼저 온 경우", &ReqSwapInfo{TxType: EventID_P2C, TxStatus: SWAP_status_init}, true},
		{"프론트가 성공을 먼저 알린 경우", &ReqSwapInfo{TxType: EventID_P2C, TxStatus: SWAP_status_fee_transfer_success}, true},
		{"출금 시작 (같은 콜백이 다시 옴)", &ReqSwapInfo{TxType: EventID_P2C, TxStatus: SWAP_status_token_transfer_withdrawal_start, SwapToCoin: SwapCoin{TokenTxHash: "0x1"}}, false},
		{"출금 해시만 남은 경우", &ReqSwapInfo{TxType: EventID_P2C, TxStatus: SWAP_status_fee_transfer_success, SwapToCoin: SwapCoin{TokenTxHash: "0x1"}}, false},
		{"출금 성공", &ReqSwapInfo{TxType: EventID_P2C, TxStatus: SWAP_status_token_transfer_withdrawal_success}, false},
		{"출금 실패", &ReqSwapInfo{TxType: EventID_P2C, TxStatus: SWAP_status_token_transfer_withdrawal_fail}, false},
		{"코인->포인트는 지급 없음", &ReqSwapInfo{TxType: EventID_C2P, TxStatus: SWAP_status_init}, false},
		{"코인->코인 입금 확인", &ReqSwapInfo{TxType: EventID_C2C, TxStatus: SWAP_status_token_transfer_deposit_success}, true},
	}
	for _, c := range cases {
		if got := c.info.CanPayout(); got != c.ok {
			t.Errorf("%s : 지급 가능 여부 %v, 기대 %v", c.name, got, c.ok)
		}
	}
}

func TestTokenKind(t *testing.T) {
	ctx := func(v *auth.VerifyAuthToken) *PointManagerContext { return &PointManagerContext{VerifyValue: v} }

	cases := []struct {
		name   string
		value  *auth.VerifyAuthToken
		app    bool
		member bool
	}{
		{"토큰 없음", nil, false, false},
		{"빈 토큰", &auth.VerifyAuthToken{}, false, false},
		{"앱 로그인", &auth.VerifyAuthToken{LoginType: AuthLoginType_App, AppID: 4}, true, false},
		{"앱 로그인인데 AppID 없음", &auth.VerifyAuthToken{LoginType: AuthLoginType_App}, false, false},
		{"앱 회원 (AppID 가 있어도 앱 토큰이 아님)", &auth.VerifyAuthToken{LoginType: AuthLoginType_AppAccount, AppID: 4, AUID: 10}, false, true},
		{"웹 회원", &auth.VerifyAuthToken{LoginType: AuthLoginType_WebAccount, AUID: 10}, false, true},
		{"웹 회원인데 계정 없음", &auth.VerifyAuthToken{LoginType: AuthLoginType_WebAccount}, false, false},
		{"알 수 없는 로그인 종류", &auth.VerifyAuthToken{LoginType: 9, AppID: 4, AUID: 10}, false, false},
	}
	for _, c := range cases {
		if got := ctx(c.value).IsAppToken(); got != c.app {
			t.Errorf("%s : IsAppToken %v, 기대 %v", c.name, got, c.app)
		}
		if got := ctx(c.value).IsMemberToken(); got != c.member {
			t.Errorf("%s : IsMemberToken %v, 기대 %v", c.name, got, c.member)
		}
	}
}
