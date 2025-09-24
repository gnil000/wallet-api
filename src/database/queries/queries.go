package queries

import (
	_ "embed"
	"strings"
)

//go:embed update_deposit_wallet.sql
var UpdateDepositWallet string

//go:embed update_withdraw_wallet.sql
var UpdateWithdrawWallet string

//go:embed find_wallet.sql
var FindWallet string

//go:embed get_wallets.sql
var GetWallets string

func ToExpectQuery(query string) string {
	query = strings.ReplaceAll(query, "$", "[$]")
	query = strings.ReplaceAll(query, "(", "\\(")
	query = strings.ReplaceAll(query, ")", "\\)")
	return query
}
