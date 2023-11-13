package key

type WalletKey struct {
	WIF            string
	P2PKHAddr      string
	P2SHSegWitAddr string
	Bech32Addr     string
	FullPubKey     string
	RedeemScript   string
}
