package config

import (
	"Share-Wallet/pkg/wallet/address"
	"Share-Wallet/pkg/wallet/coin"
)

// TODO:
// - use https://github.com/spf13/viper

// WalletRoot wallet root config
type WalletRoot struct {
	AddressType  address.AddrType  `toml:"address_type"`
	CoinTypeCode coin.CoinTypeCode `toml:"coin_type"`
	Bitcoin      Bitcoin           `toml:"bitcoin"`
	Ethereum     Ethereum          `toml:"ethereum"`
	USDTERC20    ERC20             `toml:"usdt_erc20"`
	Tron         Tron              `toml:"tron"`
	USDTTRC20    TRC20             `toml:"usdt_trc20"`
	Ripple       Ripple            `toml:"ripple"`
	Tracer       Tracer            `toml:"tracer"`
	FilePath     FilePath          `toml:"file_path"`
}

// Bitcoin Bitcoin information
type Bitcoin struct {
	Host        string `toml:"host"`
	User        string `toml:"user"`
	Pass        string `toml:"pass"`
	PostMode    bool   `toml:"http_post_mode"`
	DisableTLS  bool   `toml:"disable_tls"`
	NetworkType string `toml:"network_type""`

	Block BitcoinBlock `toml:"block"`
	Fee   BitcoinFee   `toml:"fee"`
}

// BitcoinBlock block information of Bitcoin
// FIXME: keygen/signature wallet doesn't have this value
//
//	so validation can not be used
type BitcoinBlock struct {
	ConfirmationNum uint64 `toml:"confirmation_num"`
}

// BitcoinFee range of adjustment calculated fee when sending coin
type BitcoinFee struct {
	AdjustmentMin float64 `toml:"adjustment_min"`
	AdjustmentMax float64 `toml:"adjustment_max"`
}

// Ethereum information
type Ethereum struct {
	Host        string `toml:"host"`
	IPCPath     string `toml:"ipc_path"`
	WSPath      string `toml:"ws_path"`
	NetworkType string `toml:"network_type"`
	KeyDirName  string `toml:"keydir"`
}

// ERC20 information
type ERC20 struct {
	Symbol          string `toml:"symbol"`
	Name            string `toml:"name"`
	ContractAddress string `toml:"contract_address"`
	Decimals        int    `toml:"decimals"`
	AbiJson         string `toml:"abi_json"`
}

// Tron information
type Tron struct {
	Host        string `toml:"host"`
	ApiKey      string `toml:"api_key"`
	Timeout     int    `toml:"timeout"`
	NetworkType string `toml:"network_type"`
}

// TRC20 information
type TRC20 struct {
	Symbol          string `toml:"symbol"`
	Name            string `toml:"name"`
	ContractAddress string `toml:"contract_address"`
	Decimals        int    `toml:"decimals"`
	AbiJson         string `toml:"abi_json"`
}

// Ripple information
type Ripple struct {
	WebsocketPublicURL string    `toml:"websocket_public_url"`
	WebsocketAdminURL  string    `toml:"websocket_admin_url"`
	NetworkType        string    `toml:"network_type"`
	API                RippleAPI `toml:"api"`
}

// RippleAPI is ripple-lib server info
type RippleAPI struct {
	URL      string       `toml:"url"`
	IsSecure bool         `toml:"is_secure"`
	TxData   RippleTxData `toml:"transaction"`
}

// RippleTxData is used for api command to send coin
type RippleTxData struct {
	Account string `toml:"sender_account"`
	Secret  string `toml:"sender_secret"`
}

// Logger logger info
type Logger struct {
	Service  string `toml:"service"`
	Env      string `toml:"env"`
	Level    string `toml:"level"`
	IsLogger bool   `toml:"is_logger"`
}

// Tracer is open tracing
type Tracer struct {
	Type    string       `toml:"type"`
	Jaeger  TracerDetail `toml:"jaeger"`
	Datadog TracerDetail `toml:"datadog"`
}

// TracerDetail includes specific service config
type TracerDetail struct {
	ServiceName         string  `toml:"service_name"`
	CollectorEndpoint   string  `toml:"collector_endpoint"`
	SamplingProbability float64 `toml:"sampling_probability"`
	IsDebug             bool    `toml:"is_debug"`
}

// MySQL MySQL info
type MySQL struct {
	Host  string `toml:"host"`
	DB    string `toml:"dbname"`
	User  string `toml:"user"`
	Pass  string `toml:"pass"`
	Debug bool   `toml:"debug"`
}

// FilePath if file path group
type FilePath struct {
	Tx         string `toml:"tx"`
	Address    string `toml:"address"`
	FullPubKey string `toml:"full_pubkey"`
}

// PubKeyFile saved pubKey file path which is used when import/export file
type PubKeyFile struct {
	BasePath string `toml:"base_path"`
}

// AddressFile saved address file path which is used when import/export file
type AddressFile struct {
	BasePath string `toml:"base_path"`
}
