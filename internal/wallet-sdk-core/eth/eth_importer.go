package eth

import (
	"Share-Wallet/pkg/common/log"
	db "Share-Wallet/pkg/db/local_db"
	"Share-Wallet/pkg/db/local_db/model_struct"
	"Share-Wallet/pkg/struct/wallet_api"
	"Share-Wallet/pkg/utils"
	"Share-Wallet/pkg/wallet/account"
	"Share-Wallet/pkg/wallet/address"
	"Share-Wallet/pkg/wallet/key"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// PrivKey type
type PrivKey struct {
	logger *zap.Logger
	Db     *db.DataBase
	// accountKeyRepo coldrepo.AccountKeyRepositorier
	// wtype          wallet.WalletType
}

// NewPrivKey returns privKey object
func NewPrivKey(
	Db *db.DataBase,
	// logger *zap.Logger,
	// accountKeyRepo coldrepo.AccountKeyRepositorier,
	// wtype wallet.WalletType,
) *PrivKey {
	return &PrivKey{
		Db: Db,
		// logger: logger,
		// accountKeyRepo: accountKeyRepo,
		// wtype:          wtype,
	}
}

type AccountKey struct { // ID
	ID int64
	// coin type code
	Coin string
	// account type
	Account string
	// address as standard pubkey script that Pays To PubKey Hash (P2PKH)
	P2PKHAddress string
	// p2sh-segwit address
	P2SHSegwitAddress string
	// bech32 address
	Bech32Address string
	// full public key
	FullPublicKey string
	// multisig address
	MultisigAddress string
	// redeedScript after multisig address generated
	RedeemScript string
	// WIF
	WalletImportFormat string
	// index for hd wallet
	Idx int64
	// progress status for address generating
	AddrStatus int8
	// updated date
	UpdatedAt time.Time
}

// Import imports privKey for accountKey
func (p *PrivKey) Import(accountType account.AccountType, walletKeys []key.WalletKey, userID, secret string, supportTokens []*wallet_api.SupportTokenAddress, coinType uint8) error {

	for _, record := range walletKeys {

		//check generated address
		// paddress, err := p.eth.ImportRawKey(record.WIF, secret)
		// if err != nil {
		// 	return errors.Wrap(err, "failed in ImportRawKey()")
		// }
		//record.P2PKHAddr = strings.ToLower(record.P2PKHAddr)
		// if paddress != strings.ToLower(record.P2PKHAddr) {
		// 	fmt.Println("inconsistency between generated address")
		// }

		//fmt.Println("Original Pkey", record.WIF)
		pKeyEncrypted, _ := utils.EncryptAES(record.WIF, secret)
		//dd, _ := utils.DecryptAES(pKeyEncrypted, secret)
		//fmt.Println("Decrypted Pkey", dd)

		//insert coin address
		localWallet := model_struct.LocalWallet{
			Account:            accountType.String(),
			WalletImportFormat: pKeyEncrypted,
			AddrStatus:         address.AddrStatusPrivKeyImported.String(),
			UserID:             userID,
			P2PKHAddress:       record.P2PKHAddr,
			P2SHSegwitAddress:  record.P2SHSegWitAddr,
			FullPublicKey:      record.FullPubKey,
			Bech32Address:      record.Bech32Addr,
			CoinType:           coinType,
			CreateTime:         time.Now(),
		}
		//err := p.Db.UpdateLocalWalletAddress(&localWallet)
		err := p.Db.InsertLocalWallet(&localWallet)
		if err != nil {
			return errors.Wrap(err, "fail to update UpdateLocalWalletV1")
		}

		//insert token address
		fmt.Println("\nsupportTokens", supportTokens)
		if supportTokens != nil {
			for _, token := range supportTokens {
				if token.BelongCoin == coinType {
					localUSDTERC20 := model_struct.LocalWallet{
						Account:            accountType.String(),
						WalletImportFormat: pKeyEncrypted,
						AddrStatus:         address.AddrStatusPrivKeyImported.String(),
						UserID:             userID,
						P2PKHAddress:       record.P2PKHAddr,
						P2SHSegwitAddress:  record.P2SHSegWitAddr,
						FullPublicKey:      record.FullPubKey,
						Bech32Address:      record.Bech32Addr,
						ContractAddress:    token.ContractAddress,
						CoinType:           token.CoinType,
						CreateTime:         time.Now(),
					}
					err = p.Db.InsertLocalWallet(&localUSDTERC20)
					if err != nil {
						log.NewError("", "InsertLocalWallet error", err.Error(), "coin type", token.CoinType)
					}
				}
			}
		}

	}

	return nil
}
