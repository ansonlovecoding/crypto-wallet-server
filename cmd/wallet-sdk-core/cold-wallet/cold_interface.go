package cold_wallet

import (
	seed_phrase "Share-Wallet/internal/wallet-sdk-core/seed_phrase"
	syncr "Share-Wallet/internal/wallet-sdk-core/synchronization"
	"Share-Wallet/pkg/common/config"
	"Share-Wallet/pkg/common/constant"
	sdkstruct "Share-Wallet/pkg/sdk_struct"
	"fmt"
)

func TestCold() string {
	fmt.Println(config.Config.SDKVersion)
	return "I'm cold interface!"
}

// SDKVersion func returns SDK version
func SDKVersion() string {
	return fmt.Sprintf("%s%s%s", constant.SdkVersion, constant.BigVersion, constant.UpdateVersion)
}

// SetPasscode func set passcode for login userID and returns true if success.
// accountType: 1 create, 2 recover
func SetPasscode(Password string, accountType int) (bool, error) {
	return seed_phrase.Wg.SetPasscode(Password, accountType)
}

// VerifyPasscode func verify passcode for userID and returns true if success.
func VerifyPasscode(Password string) bool {
	return seed_phrase.Wg.VerifyPasscode(Password)
}

// GenerateSeedPhrase func generates seed phrase for userID.
func GenerateSeedPhrase(lang, secret string) (string, error) {
	return seed_phrase.Wg.GenerateSeedPharse(lang, secret)
}

// VerifySeedPharse func verify seed phrase for userID.
func VerifySeedPharse(userID, seed string) (bool, error) {
	return seed_phrase.Wg.VerifySeedPhrase(userID, seed)
}

// FetchSeedPhraseRandom func fetches seed phrase in random order.
func FetchSeedPhraseRandom() (string, error) {
	return seed_phrase.Wg.FetchSeedPhraseRandom()
}

// CheckUserExist func fetches seed phrase in random order.
func CheckUserExist(userID string) bool {
	return seed_phrase.Wg.CheckUserExist(userID)
}

// FetchSeedPhraseWord func returns seed phrase word suggestions.
func FetchSeedPhraseWord(word string) string {
	return seed_phrase.Wg.FetchSeedPhraseWord(word)
}

// FetchLocalCoins func fetches coin list
func FetchLocalCoins() string {
	return seed_phrase.Wg.FetchLocalCoins()
}

// InitWalletSDK func initialise the sdk with database connection.
func InitWalletSDK(userID string, config string) bool {
	return seed_phrase.InitWallet(userID, config)
}

// GenerateUserAccount func calls setpasscode,generateseedphrase and register account then return seed phrase
/* Supported languages (lang)
//	EN    = "en"
//	CHSim = "ch-sim"
//	CHTra = "ch-tra"
//	FR    = "fr"
//	IT    = "it"
//	JA    = "ja"
//	KO    = "ko"
//	SP    = "sp"
//	*/
func GenerateUserAccount(lang, secret string) (string, error) {
	return seed_phrase.Wg.GenerateUserAccount(lang, secret)
}

// GetUserStatus func fetches the status of user.
func GetUserStatus() int {
	return seed_phrase.Wg.GetUserStatus()
}

// GetUser func fetches user sensitive data.
func GetUserSeedPhrase(secret string) string {
	return seed_phrase.Wg.GetUser(secret)
}

// RecoverUserAccount func calls setpasscode,register account.
func RecoverUserAccount(secret, seedPhrase string) (bool, error) {
	return seed_phrase.Wg.RecoverUserAccount(secret, seedPhrase)
}

// AddAddressBook func add address and  name to local db
func AddAddressBook(coinType int, address, name string) (result bool, err error) {
	return seed_phrase.Wg.AddAddressBook(coinType, address, name)
}

// GetLocalUserAddressBook func returns addressbook
func GetLocalUserAddressBook(coinType int) string {
	return seed_phrase.Wg.GetAddressBook(coinType)
}

// GetUser func fetches user sensitive data.
func GetUserPrivateKey(secret string, coinType int, keyType int) string {
	return seed_phrase.Wg.GetUserPrivateKey(secret, coinType, keyType)
}

// DeleteAddressBook func deletes address book
func DeleteAddressBook(coinType int, address string) bool {
	return seed_phrase.Wg.DeleteAddressBook(coinType, address)
}

// GetLocalUserAddressBook func returns addressbook
func GetLocalUserAddressBookByAddress(coinType int, address string) string {
	return seed_phrase.Wg.GetAddressBookByAddress(coinType, address)
}

//SynchronizeUserAccount func creates account information in server
func SynchronizeUserAccount(userID, publicKey string) bool {
	return syncr.Sg.Synchronize(userID, publicKey, sdkstruct.SvrConf.Platform)
}
