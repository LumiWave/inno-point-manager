package token_manager_server

const (
	ERR0003 = "ERR-0003"
	ERR1103 = "ERR-1103"
)

var ResultCodeText = map[string]string{
	ERR0003: "Invalid to wallet address",
	ERR1103: "Insufficient balance",
}
