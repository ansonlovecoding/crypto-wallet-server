package coin

import "github.com/btcsuite/btcd/chaincfg"

// CoinType creates a separate subtree for every cryptocoin
type CoinType uint32

// Uint32 is converter
func (c CoinType) Uint32() uint32 {
	return uint32(c)
}

// coin_type
// https://github.com/satoshilabs/slips/blob/master/slip-0044.md
const (
	CoinTypeBitcoin     CoinType = 0   // Bitcoin
	CoinTypeTestnet     CoinType = 1   // Testnet (all coins)
	CoinTypeLitecoin    CoinType = 2   // Litecoin
	CoinTypeEther       CoinType = 60  // Ether
	CoinTypeRipple      CoinType = 144 // Ripple
	CoinTypeBitcoinCash CoinType = 145 // Bitcoin Cash
	CoinTypeTrx         CoinType = 195 // TRX
	// Token
	CoinTypeUSDTERC20 CoinType = 9000 // TODO: temporary
	CoinTypeUSDTTRC20 CoinType = 9001
)

// CoinTypeCode coin type code
type CoinTypeCode string

// coin_type_code
const (
	BTC       CoinTypeCode = "btc"
	ETH       CoinTypeCode = "eth"
	USDTERC20 CoinTypeCode = "usdterc20"
	USDTTRC20 CoinTypeCode = "usdttrc20"
	TRX       CoinTypeCode = "trx"
	BCH       CoinTypeCode = "bch"
	LTC       CoinTypeCode = "ltc"
	XRP       CoinTypeCode = "xrp"
)

// String converter
func (c CoinTypeCode) String() string {
	return string(c)
}

// CoinType returns CoinType
func (c CoinTypeCode) CoinType(conf *chaincfg.Params) CoinType {
	if conf.Name != "mainnet" {
		return CoinTypeTestnet
	}
	if coinType, ok := CoinTypeCodeValue[c]; ok {
		return coinType
	}
	// coinType could not found
	return CoinTypeTestnet
}

// CoinTypeCodeValue value
var CoinTypeCodeValue = map[CoinTypeCode]CoinType{
	BTC:       CoinTypeBitcoin,
	BCH:       CoinTypeBitcoinCash,
	LTC:       CoinTypeLitecoin,
	ETH:       CoinTypeEther,
	XRP:       CoinTypeRipple,
	TRX:       CoinTypeTrx,
	USDTERC20: CoinTypeUSDTERC20,
	USDTTRC20: CoinTypeUSDTTRC20,
}

// IsCoinTypeCode validate
func IsCoinTypeCode(val string) bool {
	if _, ok := CoinTypeCodeValue[CoinTypeCode(val)]; ok {
		return true
	}
	return false
}

// IsBTCGroup validates bitcoin group
func IsBTCGroup(val CoinTypeCode) bool {
	if val == BTC || val == BCH {
		return true
	}
	return false
}

// IsETHGroup validates ethreum, ERC20 group
func IsETHGroup(val CoinTypeCode) bool {
	if val == ETH || val == USDTERC20 {
		return true
	}
	return false
}

// IsTronGroup validates ethreum, ERC20 group
func IsTronGroup(val CoinTypeCode) bool {
	if val == TRX || val == USDTTRC20 {
		return true
	}
	return false
}
