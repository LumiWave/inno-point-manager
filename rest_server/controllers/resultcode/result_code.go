package resultcode

const (
	Result_Success = 0

	Result_Require_PageInfo               = 10001 // 유효한 페이지 정보 필요
	Result_Require_MemberIdx              = 11002 // 맴버 index 정보 푤이
	Result_Require_ValidPrivateWalletAddr = 11003 // 유효한 private 지갑 주소 필요
	Result_Require_ValidPublicWalletAddr  = 11004 // 유요한 public 지갑 주소 필요
	Result_Require_ValidPointAmount       = 11005 // 포인트 정보 필요
	Result_Require_PointType              = 11006 // point 이벤트 타입 정보 필요
	Result_Require_LatestPointAmount      = 11007 // 마지막 app 포인트 정보 필요
	Result_Require_AdjustPointAmount      = 11008 // app 포인트 변화량 필요
	Result_Require_PrivateTokenAmount     = 11009 // private 토큰 정보 필요

	Result_Error_LatestPointAmountIsDiffrent = 11100 // 마지막 포인트 정보가 다르다.
	Result_Error_NotExistMember              = 11101 // 존재하지 않는 member
	Result_Error_LackOfTokenQuantity         = 11102 // 토큰 수량이 부족하다.

	Result_Require_AUID           = 12000 // 유효한 au_id 정보 필요
	Result_Require_MUID           = 12001 // 유효한 mu_id 정보 필요
	Result_Require_AppID          = 12002 // 유효한 app_id 정보 필요
	Result_Require_DatabaseID     = 12003 // 유효한 database_id 정보 필요
	Result_Require_PointID        = 12004 // 유효한 point_id 정보 필요
	Result_Require_AdjustQuantity = 12005 // 유효한 adjust_quantity 정보 필요

	Result_RedisError_Lock_fail = 18000 // redis lock error

	Result_DBError         = 19000 // db 에러
	Result_Invalid_DBID    = 19001 // 유효하지 못한 database index
	Result_DBError_Unknown = 19002 // 알려지지 않은 db 에러

	Result_Error_Invalid_data   = 50001 // 요청 조건값에 대한 데이터가 존재하지 않습니다.
	Result_Error_duplicate_auid = 50102 // 해당 App에 중복된 AUID가 있습니다.

	Result_Auth_RequireMessage    = 20000
	Result_Auth_RequireSign       = 20001
	Result_Auth_InvalidLoginInfo  = 20002
	Result_Auth_DontEncryptJwt    = 20003
	Result_Auth_InvalidJwt        = 20004
	Result_Auth_InvalidWalletType = 20005
)

var ResultCodeText = map[int]string{
	Result_Success: "success",

	Result_Require_PageInfo:               "require page info",
	Result_Require_MemberIdx:              "require member index",
	Result_Require_ValidPrivateWalletAddr: "require valid private wallet address",
	Result_Require_ValidPublicWalletAddr:  "require valid public wallet address",
	Result_Require_ValidPointAmount:       "require valid point amount",
	Result_Require_PointType:              "require point type",
	Result_Require_LatestPointAmount:      "require latest point amount",
	Result_Require_AdjustPointAmount:      "require adjust point amount",
	Result_Require_PrivateTokenAmount:     "require private token amount",

	Result_Error_LatestPointAmountIsDiffrent: "latest point information is different",
	Result_Error_NotExistMember:              "Not exist member",
	Result_Error_LackOfTokenQuantity:         "Lack of token quantity",

	Result_RedisError_Lock_fail: "Redis lock error.",

	Result_DBError:         "Internal DB error",
	Result_Invalid_DBID:    "Invalid DB ID",
	Result_DBError_Unknown: "Unknown DB error",

	Result_Error_Invalid_data: "	Invalid data received.",
	Result_Error_duplicate_auid: "The app has duplicate AUIDs.",

	Result_Require_AUID:           "Requires valid 'au_id' information.",
	Result_Require_MUID:           "Requires valid 'mu_id' information.",
	Result_Require_AppID:          "Requires valid 'app_id' information.",
	Result_Require_DatabaseID:     "Requires valid 'database_id' information.",
	Result_Require_PointID:        "Requires valid 'point_id' information.",
	Result_Require_AdjustQuantity: "Requires valid 'adjust_quantity' information.",

	Result_Auth_RequireMessage:    "Message is required",
	Result_Auth_RequireSign:       "Sign info is required",
	Result_Auth_InvalidLoginInfo:  "Invalid login info",
	Result_Auth_DontEncryptJwt:    "Auth token create fail",
	Result_Auth_InvalidJwt:        "Invalid jwt token",
	Result_Auth_InvalidWalletType: "Invalid wallet type",
}
