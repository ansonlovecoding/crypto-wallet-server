package model_struct

import "time"

type LocalUser struct {
	ID           uint64    `gorm:"column:id;primary_key;" json:"id"`
	UserID       string    `gorm:"column:user_id;unique;type:varchar(255);" json:"user_id"`
	WalletID     string    `gorm:"column:walletID;type:varchar(255)" json:"wallet_id"`
	UserName     string    `gorm:"column:user_name" json:"user_name"`
	Password     string    `gorm:"column:password" json:"password"`
	EntropyLevel uint32    `gorm:"column:entropy_level;type:varchar(255)" json:"entropy_level"`
	PublicKey    string    `gorm:"column:public_key;type:varchar(255)" json:"public_key"`
	SeedPhrase   string    `gorm:"column:seed_phrase;type:varchar(1024)" json:"seed_phrase"`
	AccountType  uint8     `gorm:"column:account_type" json:"account_type"`
	Status       uint8     `gorm:"column:status" json:"status"`
	CreateTime   time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time" json:"update_time"`
	DeleteTime   time.Time `gorm:"column:delete_time" json:"delete_time"`
}

type LocalWallet struct {
	ID                 uint64 `gorm:"column:id;primary_key;" json:"id"`
	UserID             string `gorm:"column:user_id;type:varchar(255)" json:"user_id"`
	PublicKey          string `gorm:"column:public_key" json:"public_key"`                           // Depreciated
	WalletAddress      string `gorm:"column:wallet_address;type:varchar(255)" json:"wallet_address"` //Depreciated
	Status             uint8  `gorm:"column:status" json:"status"`
	CoinType           uint8  `gorm:"column:coin_type" json:"coin_type"`
	Account            string `gorm:"column:account" json:"account"`
	P2PKHAddress       string `gorm:"column:p2pkh_address;comment:coin address for addr_type 1" json:"p2pkh_address"`
	P2SHSegwitAddress  string `gorm:"column:p2sh_segwit_address" json:"p2sh_segwit_address"`
	ContractAddress    string `gorm:"column:contract_address;comment:coin address for addr_type 1" json:"contract_address"`
	Bech32Address      string `gorm:"column:bech32_address" json:"bech32_address"`
	FullPublicKey      string `gorm:"column:full_public_key" json:"full_public_key"`
	MultisigAddress    string `gorm:"column:multisig_address" json:"multisig_address"`
	RedeemScript       string `gorm:"column:redeem_script" json:"redeem_script"`
	WalletImportFormat string `gorm:"column:wallet_import_format"`
	Idx                int64  `gorm:"column:idx" json:"idx"`
	AddrStatus         string `gorm:"column:addr_status" json:"addr_status"`
	AddrType           uint8  `gorm:"column:addr_type;comment:1 coin address 2 token address" json:"addr_type"`

	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
	DeleteTime time.Time `gorm:"column:delete_time" json:"delete_time"`
}

type LocalWalletType struct {
	ID          uint64    `gorm:"column:id;primary_key;" json:"id"`
	CoinType    uint8     `gorm:"column:coin_type;unique;" json:"coin_type"`
	CoinName    string    `gorm:"column:coin_name;unique;type:varchar(255)" json:"coin_name"`
	Description string    `gorm:"column:description" json:"description"`
	Balance     float64   `gorm:"column:balance" json:"balance"`
	Status      uint8     `gorm:"column:status" json:"status"`
	CreateTime  time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:update_time" json:"update_time"`
	DeleteTime  time.Time `gorm:"column:delete_time" json:"delete_time"`
}
type LocalUserAddressBook struct {
	ID         uint64    `gorm:"column:id;primary_key;" json:"id"`
	UserID     string    `gorm:"column:user_id;type:varchar(255)" json:"user_id"`
	Name       string    `gorm:"column:name" json:"name"`
	Address    string    `gorm:"column:address" json:"address"`
	CoinType   uint8     `gorm:"column:coin_type" json:"coin_type"`
	Status     uint8     `gorm:"column:status" json:"status"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
	DeleteTime time.Time `gorm:"column:delete_time" json:"delete_time"`
}
