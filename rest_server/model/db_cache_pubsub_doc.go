package model

var (
	gIsMaintenance            = false
	gIsExternalTransferEnable = true
	gIsSwapEnable             = true
	gIsPointUpdateEnable      = true
)

func SetMaintenance(isMaintenance bool) {
	gIsMaintenance = isMaintenance
}

func GetMaintenance() bool {
	return gIsMaintenance
}

func SetExternalTransferEnable(isEnable bool) {
	gIsExternalTransferEnable = isEnable
}

func GetExternalTransferEnable() bool {
	return gIsExternalTransferEnable
}

func SetSwapEnable(isEnable bool) {
	gIsSwapEnable = isEnable
}

func GetSwapEnable() bool {
	return gIsSwapEnable
}

func SetPointUpdateEnable(isEnable bool) {
	gIsPointUpdateEnable = isEnable
}

func GetPointUpdateEnable() bool {
	return gIsPointUpdateEnable
}
