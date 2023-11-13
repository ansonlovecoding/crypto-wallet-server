package seed_pharse

import (
	"Share-Wallet/internal/wallet-sdk-core/bip32"
	"Share-Wallet/internal/wallet-sdk-core/bip39"
	ep "Share-Wallet/internal/wallet-sdk-core/eth"
	"Share-Wallet/internal/wallet-sdk-core/register"
	"Share-Wallet/internal/wallet-sdk-core/synchronization"
	"Share-Wallet/internal/wallet-sdk-core/transfer"
	"Share-Wallet/pkg/common/constant"
	"Share-Wallet/pkg/common/http"
	db "Share-Wallet/pkg/db/local_db"
	"Share-Wallet/pkg/db/local_db/model_struct"
	sdkstruct "Share-Wallet/pkg/sdk_struct"
	"Share-Wallet/pkg/struct/wallet_api"
	"Share-Wallet/pkg/utils"
	"Share-Wallet/pkg/wallet/account"
	"Share-Wallet/pkg/wallet/coin"
	"Share-Wallet/pkg/wallet/key"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	"gorm.io/gorm"
)

type WalletMgr struct {
	Db                    *db.DataBase
	LoginUserID           string
	LoginTime             int64
	UserRandomSecret      []byte
	SupportTokenAddresses []*wallet_api.SupportTokenAddress
}

var Wg *WalletMgr

// GenerateSeedPharse function return seed phrase and add generated fields to local DB.
// Functionality part is considered for the time being.
// Encryption and hashing is not applied now.

func (w *WalletMgr) GenerateSeedPharse(lang, key string) (string, error) {
	userID := w.LoginUserID

	// Check userID is valid
	if len(userID) == 0 {
		return "", errors.New("userID is empty")
	}

	userExist, err := w.Db.GetUserByUserID(userID)
	if err != nil {
		return "", fmt.Errorf("error in GetUserByUserID() %w", err)
	}
	if userExist.Status != 1 {
		return "", errors.New("seed phrase has been generated before")
	}

	// Set wordlist based on language. Default one is English
	/* Supported languages
	EN    = "en"
	CHSim = "ch-sim"
	CHTra = "ch-tra"
	FR    = "fr"
	IT    = "it"
	JA    = "ja"
	KO    = "ko"
	SP    = "sp"
	*/
	bip39.SetWordListLanguage(lang)

	// Fetching entropy level using 128; will give 12 words
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		log.Println("error in NewEntropy()")
	}
	// Fetching seed phrase from BIP39
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		log.Println("Error in NewMnemonic()")
	}

	encSeedPhrase, err := utils.EncryptAES(mnemonic, key)
	if err != nil {
		log.Println("Error in EncryptAES()")
	}

	user := model_struct.LocalUser{
		UserID:       userID,
		EntropyLevel: 128,
		SeedPhrase:   string(encSeedPhrase),
		Status:       2, // Status 2: Seed phrase created. but not confirmed by user
	}
	// decy, _ := utils.DecryptAES(encSeedPhrase, key)

	// wallet := model_struct.LocalWallet{
	// 	UserID:        userID,
	// 	PublicKey:     "",
	// 	WalletAddress: "",
	// 	Status:        128,
	// 	CoinType:      1, // Assumping the first coin, can link with the coin_type table in next phase.
	// 	CreateTime:    time.Now(),
	// }

	err = w.Db.UpdateLocalUser(&user)
	if err != nil {
		return "", fmt.Errorf("error in UpdateLocalUser() %w", err)
	}
	// err = w.Db.InsertLocalWallet(&wallet)
	// if err != nil {
	// 	return "", fmt.Errorf("error in InsertLocalUser() %w", err)
	// 	// Todo: Rollback required.

	// }

	return mnemonic, nil

}

func (w *WalletMgr) SetPasscode(passcode string, accountType int) (bool, error) {
	userID := w.LoginUserID
	if passcode == "" {
		return false, errors.New("passcode is empty")
	}
	userExist, err := w.Db.GetUserByUserID(userID)
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return false, fmt.Errorf("error in GetUserByUserID() %w", err)
	}
	if userExist.UserID == userID {
		return false, errors.New("user exists with UserID")
	}
	if len(passcode) != 6 {
		return false, errors.New("passcode length restricted to six characters")
	}

	// password hashing
	has := md5.Sum([]byte(passcode))
	password := fmt.Sprintf("%x", has)

	// user_model

	user := model_struct.LocalUser{
		UserID:      userID,
		Password:    password,
		AccountType: uint8(accountType), // 1-> Account generated using our platform
		Status:      1,                  // 1-> Passcode is generated and pending seed phrase creation
	}
	// Adding user to local DB.
	// Server synchronisation may require in next phase.

	err = w.Db.InsertLocalUser(&user)
	if err != nil {
		return false, fmt.Errorf("error in InsertLocalUser() %w", err)
	}
	return true, nil
}

func (w *WalletMgr) VerifyPasscode(passcodeFromUser string) bool {

	userID := w.LoginUserID
	if passcodeFromUser == "" {
		log.Println(errors.New("passcode is empty"))
		return false
	}
	if userID == "" {
		log.Println(errors.New("userID is empty"))
		return false
	}
	if len(passcodeFromUser) < 6 || len(passcodeFromUser) > 20 {
		log.Println(errors.New("passcode length should be between 6 and 20 character"))
		return false
	}

	user, err := w.Db.GetUserByUserID(userID)
	if err != nil {
		log.Println(fmt.Errorf("error in GetUserByUserID() %w", err))
		return false
	}
	has := md5.Sum([]byte(passcodeFromUser))
	password := fmt.Sprintf("%x", has)
	if user.Password != password {
		log.Println(errors.New("passcode mismatch"))
		return false
	}
	// Password maximum login attempt can consider in next phase
	// Also token.
	return true
}
func (w *WalletMgr) VerifySeedPhrase(userID, seedFromRequest string) (bool, error) {

	if seedFromRequest == "" {
		return false, errors.New("seedphrase is empty")
	}
	if userID == "" {
		return false, errors.New("userID is empty")
	}

	user, err := w.Db.GetUserByUserID(userID)
	if err != nil {
		return false, fmt.Errorf("error in GetUserByUserID() %w", err)
	}

	// Hashing can consider in the next phase
	// seedHash := md5.Sum([]byte(seed))
	// seedHashFormat := fmt.Sprintf("%x", seedHash)

	if user.SeedPhrase != seedFromRequest {
		return false, errors.New("seed mismatch")
	}

	userModel := model_struct.LocalUser{
		UserID: userID,
		Status: 3, // Status 3: Seed phrase verified by user.
	}
	err = w.Db.UpdateLocalUserStatus(&userModel)
	if err != nil {
		return false, fmt.Errorf("error in UpdateLocalUser() %w", err)
	}

	// To Do: Restrict user with too many incorrect attempt.For eg: If consecutive 3 attempt failed we can block/ suggest the user to generate the seed phrase again
	return true, nil
}

func (w *WalletMgr) FetchSeedPhraseRandom() (string, error) {
	userID := w.LoginUserID
	var (
		seedPhraseList []string
		// seedPhraseListCommaSeperated go mobile not supporting array
		seedPhraseListCommaSeperated string
	)

	if userID == "" {
		return seedPhraseListCommaSeperated, errors.New("userID is empty")
	}

	user, err := w.Db.GetUserByUserID(userID)
	if err != nil {
		return seedPhraseListCommaSeperated, fmt.Errorf("error in GetUserByUserID() %w", err)
	}

	seedPhraseList = strings.Split(user.SeedPhrase, " ")

	if len(seedPhraseList) != 12 {
		return seedPhraseListCommaSeperated, errors.New("seedPhraseList is invalid")
	}

	// To do: Can improve randomness
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(seedPhraseList), func(i, j int) { seedPhraseList[i], seedPhraseList[j] = seedPhraseList[j], seedPhraseList[i] })

	// To do: Restrict user with too many incorrect attempt.For eg: If consecutive 3 attempt failed we can block/ suggest the user to generate the seed phrase again
	seedPhraseListCommaSeperated = strings.Join(seedPhraseList, "")
	return seedPhraseListCommaSeperated, nil
}

func (w *WalletMgr) CheckUserExist(userID string) bool {

	if userID == "" {
		log.Println("userID is empty")
		return false
	}

	user, err := w.Db.GetUserByUserID(userID)
	if err != nil {
		log.Println(fmt.Errorf("error in GetUserByUserID() %w", err))
		return false
	}
	if user != nil {
		return true
	}

	return false
}
func (w *WalletMgr) FetchSeedPhraseWord(str string) string {

	var (
		wordlist []string
		counter  int
	)
	if str == "" {
		log.Println(errors.New("word is empty"))
		return ""
	}

	str = fmt.Sprintf("%s", strings.Trim(str, ""))
	str = fmt.Sprintf("%s", strings.ToLower(str))
	bip39.SetWordListLanguage("en")
	list := bip39.GetWordList()
	for _, v := range list {
		if strings.HasPrefix(v, str) {
			wordlist = append(wordlist, v)
			counter++
		}
		if counter == 4 {
			return strings.Join(wordlist, " ")

		}
	}
	return strings.Join(wordlist, " ")
}

func (w *WalletMgr) FetchLocalCoins() string {
	walletCoinType, err := w.Db.GetLocalWalletCoinType()
	if err != nil {
		log.Println(fmt.Errorf("error in FetchLocalCoins() %w", err))
		return ""
	}
	var list []sdkstruct.WalletCoinType
	for _, v := range walletCoinType {
		list = append(list, sdkstruct.WalletCoinType{
			CoinType:    int(v.CoinType),
			Name:        v.CoinName,
			Description: v.Description,
			Balance:     v.Balance,
		})
	}
	return utils.StructToJsonString(list)
}

func NewWalletMgr(dataBase *db.DataBase, loginUserID string, loginTime int32, userRandomSecret []byte) (w *WalletMgr) {
	return &WalletMgr{Db: dataBase, LoginUserID: loginUserID, LoginTime: int64(loginTime), UserRandomSecret: userRandomSecret}
}

func InitWallet(userID string, config string) bool {

	var (
		err error
	)
	if err = json.Unmarshal([]byte(config), &sdkstruct.SvrConf); err != nil {
		return false
	}
	if sdkstruct.SvrConf.DataDir == "" {
		return false
	}

	db, err := db.NewDataBase(userID, sdkstruct.SvrConf.DataDir)
	if err != nil {
		log.Println("Database initialization failed")
		return false
	}
	randomUserSecret, ok := configureSecretFile(sdkstruct.SvrConf.DataDir, userID)
	if !ok {
		return false
	}

	token := ""
	// Move to config/db later.
	// Add HMAC Token for API authentication
	postAPI := http.NewPostApi(token, sdkstruct.SvrConf.ApiAddr)
	Wg = NewWalletMgr(db, userID, int32(time.Now().Unix()), randomUserSecret)

	//request token addresses
	var supportTokenAddresses []*wallet_api.SupportTokenAddress
	req := wallet_api.GetSupportTokenAddressesRequest{
		OperationID: utils.OperationIDGenerator(),
	}
	resp, err := postAPI.PostWalletAPI(constant.GetSupportTokenAddressURL, req, constant.APITimeout)
	if err == nil && resp != nil {
		var respObj wallet_api.GetSupportTokenAddressesBaseResp
		err = utils.JsonStringToStruct(string(resp), &respObj)
		if err == nil {
			log.Println("GetSupportTokenAddressURL AddressList:", respObj.Data.AddressList)
			supportTokenAddresses = respObj.Data.AddressList
		}

	}

	Wg.SupportTokenAddresses = supportTokenAddresses
	register.RegistryMgr = register.NewRegistry(db, userID, int32(time.Now().Unix()), postAPI, randomUserSecret, supportTokenAddresses)
	synchronization.Sg = synchronization.NewWalletMgr(db, userID, int32(time.Now().Unix()), randomUserSecret, postAPI)
	transfer.TransferMgr = transfer.NewTransfer(db, userID, int32(time.Now().Unix()), postAPI, randomUserSecret, supportTokenAddresses)

	//if the second time call initwallet, need to call the synchronization function
	user, err := db.GetUserByUserID(userID)
	if err == nil && user != nil {
		synchronization.Sg.Synchronize(userID, user.PublicKey, sdkstruct.SvrConf.Platform)
	}

	log.Println("SDK is successfully initialized")
	return true
}

func configureSecretFile(dirPath, userID string) ([]byte, bool) {
	// Remove Secret from logs, added only for testing
	secretFilePath := fmt.Sprintf("%s%s%s%s", dirPath, "/Wallet_Secret_", userID, ".txt")
	var secretByte []byte
	if _, err := os.Stat(secretFilePath); errors.Is(err, os.ErrNotExist) {
		secret := utils.RandomSecretString(26)
		err := ioutil.WriteFile(secretFilePath, []byte(secret), 0777)
		if err != nil {
			log.Fatal(err)
			return secretByte, false
		}
		secretByte = []byte(secret)
		log.Println("File Created:", string(secretByte))
		return []byte(secret), true
	} else {
		secretByte, err = ioutil.ReadFile(secretFilePath)
		if err != nil {
			log.Fatal(err)
			return secretByte, false
		}
	}
	log.Println("File Read:", string(secretByte))
	return secretByte, true
}

// GenerateUserAccount fn calls set passcode, seedphrase generation and account generation
func (w *WalletMgr) GenerateUserAccount(lang, secret string) (string, error) {
	// check if w.SupportTokenAddresses exist, if not request it again
	if w.SupportTokenAddresses == nil || len(w.SupportTokenAddresses) == 0 {
		//request token addresses
		token := ""
		// Move to config/db later.
		// Add HMAC Token for API authentication
		postAPI := http.NewPostApi(token, sdkstruct.SvrConf.ApiAddr)
		req := wallet_api.GetSupportTokenAddressesRequest{
			OperationID: utils.OperationIDGenerator(),
		}
		resp, err := postAPI.PostWalletAPI(constant.GetSupportTokenAddressURL, req, constant.APITimeout)
		if err != nil {
			return "", errors.New("Network error, please try later")
		}
		var respObj wallet_api.GetSupportTokenAddressesBaseResp
		err = utils.JsonStringToStruct(string(resp), &respObj)
		if err != nil {
			return "", errors.New("Network error, please try later")
		}
		w.SupportTokenAddresses = respObj.Data.AddressList
	}

	// To do: Refactor DB fns to use transactions to roll back db to previous state
	ok, err := w.SetPasscode(secret, constant.UserAccountCreate)
	if err != nil {
		log.Println(fmt.Errorf("error in SetPasscode() %w", err))
		return "", err
	}
	if !ok {
		log.Println(errors.New("SetPasscode failed"))
		return "", errors.New("SetPasscode failed")
	}

	log.Printf("user passcode generated")
	key := fmt.Sprintf("%s%s", secret, w.UserRandomSecret)
	seedPhrase, err := w.GenerateSeedPharse(lang, key)
	if err != nil {
		log.Println(fmt.Errorf("error in GenerateSeedPharse() %w", err))
		return "", err
	}
	log.Printf("SeedPharse generated")
	err = w.CreateAccount(key)
	if err != nil {
		log.Println(fmt.Errorf("error in GenerateSeedPharse() %w", err))
		return "", err
	}
	return seedPhrase, nil
}

func (w *WalletMgr) CreateAccount(keyEnc string) error {

	userID := w.LoginUserID
	if userID == "" {
		return errors.New("UserID/Seed is empty")
	}
	user, err := w.Db.GetUserByUserID(userID)
	if err != nil {
		return fmt.Errorf("error in GetUserByUserID() %w", err)
	}
	if user.Password == "" {
		return errors.New("user password is empty")
	}
	if user.Status != 2 {
		return errors.New("seed phrase generation failed")
	}

	seedEncrpted := user.SeedPhrase
	seed, err := utils.DecryptAES(seedEncrpted, keyEnc)
	if err != nil {
		return fmt.Errorf("error in DecryptAES() %w", err)
	}
	if seed == "" {
		return errors.New("seed is nil")
	}

	//To Do: Research and refactor below 2 variable initialisations
	idxFrom := 0
	count := 1

	seedByte := bip39.NewSeed(seed, "")
	masterKey, _ := bip32.NewMasterKey(seedByte)
	publicKey := masterKey.PublicKey().String()

	//remember add the implement to the new coin to HDkeyGenerator function
	// Eth Flow
	if true {
		hdKey := key.HDkeyGenerator(coin.ETH)
		walletKeys, err := hdKey.CreateKey(seedByte, account.AccountTypeClient, uint32(idxFrom), uint32(count))
		if err != nil {
			return fmt.Errorf("error in CreateKey() %w", err)
		}

		log.Println("walletKeys", walletKeys)
		privKey := ep.NewPrivKey(w.Db)
		err = privKey.Import(account.AccountTypeClient, walletKeys, userID, keyEnc, w.SupportTokenAddresses, constant.ETHCoin)
		if err != nil {
			return fmt.Errorf("error in Import() %w", err)
		}
	}

	// TRX Flow
	if true {
		hdKey := key.HDkeyGenerator(coin.TRX)
		walletKeys, err := hdKey.CreateKey(seedByte, account.AccountTypeClient, uint32(idxFrom), uint32(count))
		if err != nil {
			return fmt.Errorf("error in CreateKey() %w", err)
		}

		log.Println("walletKeys", walletKeys)
		privKey := ep.NewPrivKey(w.Db)
		err = privKey.Import(account.AccountTypeClient, walletKeys, userID, keyEnc, w.SupportTokenAddresses, constant.TRX)
		if err != nil {
			return fmt.Errorf("error in Import() %w", err)
		}
	}

	//BTC Flow
	if true {
		hdKey := key.HDkeyGenerator(coin.BTC)
		walletKeys, err := hdKey.CreateKey(seedByte, account.AccountTypeClient, uint32(idxFrom), uint32(count))
		if err != nil {
			return fmt.Errorf("error in CreateKey() %w", err)
		}

		log.Println("walletKeys", walletKeys)
		privKey := ep.NewPrivKey(w.Db)
		err = privKey.Import(account.AccountTypeClient, walletKeys, userID, keyEnc, w.SupportTokenAddresses, constant.BTCCoin)
		if err != nil {
			return fmt.Errorf("error in Import() %w", err)
		}
	}

	userModel := model_struct.LocalUser{
		UserID:    userID,
		PublicKey: publicKey,
		Status:    4, // Status 4: Private key generated.
	}
	err = w.Db.UpdateLocalUserPublicKeyAndStatus(&userModel)
	if err != nil {
		return fmt.Errorf("error in UpdateLocalUser() %w", err)
	}
	synchronization.Sg.Synchronize(w.LoginUserID, publicKey, sdkstruct.SvrConf.Platform)
	return nil
}

func (w *WalletMgr) GetUserStatus() int {

	userID := w.LoginUserID
	if userID == "" {
		log.Println(errors.New("userID is empty"))
		return 0
	}
	user, err := w.Db.GetUserByUserID(userID)
	if err != nil {
		log.Println(fmt.Errorf("error in GetUserByUserID() %w", err))
		return 0
	}
	return int(user.Status)
}

func (w *WalletMgr) GetUser(secret string) string {

	userID := w.LoginUserID
	if userID == "" {
		log.Println(errors.New("userID is empty"))
		return ""
	}
	user, err := w.Db.GetUserByUserID(userID)
	if err != nil {
		log.Println(fmt.Errorf("error in GetUserByUserID() %w", err))
		return ""
	}

	key := fmt.Sprintf("%s%s", secret, w.UserRandomSecret)

	seed, err := utils.DecryptAES(user.SeedPhrase, key)
	if err != nil {
		log.Println(fmt.Errorf("error in DecryptAES() %w", err))
		return ""
	}
	return seed
}

// RecoverUserAccount fn calls set passcode, register account with the seed phrase
func (w *WalletMgr) RecoverUserAccount(secret, seedPhrase string) (bool, error) {
	// check if w.SupportTokenAddresses exist, if not return error
	if w.SupportTokenAddresses == nil || len(w.SupportTokenAddresses) == 0 {
		//request token addresses
		token := ""
		// Move to config/db later.
		// Add HMAC Token for API authentication
		postAPI := http.NewPostApi(token, sdkstruct.SvrConf.ApiAddr)
		req := wallet_api.GetSupportTokenAddressesRequest{
			OperationID: utils.OperationIDGenerator(),
		}
		resp, err := postAPI.PostWalletAPI(constant.GetSupportTokenAddressURL, req, constant.APITimeout)
		if err != nil {
			return false, errors.New("Network error, please try later")
		}
		var respObj wallet_api.GetSupportTokenAddressesBaseResp
		err = utils.JsonStringToStruct(string(resp), &respObj)
		if err != nil {
			return false, errors.New("Network error, please try later")
		}
		w.SupportTokenAddresses = respObj.Data.AddressList
	}

	// To do: Refactor DB fns to use transactions to roll back db to previous state
	ok, err := w.SetPasscode(secret, constant.UserAccountRecovered)
	if err != nil {
		log.Println(fmt.Errorf("error in SetPasscode() %w", err))
		return false, err
	}
	if !ok {
		log.Println(errors.New("SetPasscode failed"))
		return false, errors.New("SetPasscode failed")
	}

	log.Printf("user passcode generated")
	key := fmt.Sprintf("%s%s", secret, w.UserRandomSecret)
	err = w.EncrptSeedPhrase(key, seedPhrase)
	if err != nil {
		log.Println(fmt.Errorf("error in GenerateSeedPharse() %w", err))
		return false, err
	}
	log.Printf("SeedPharse encrypted")
	err = w.CreateAccount(key)
	if err != nil {
		log.Println(fmt.Errorf("error in CreateAccount() %w", err))
		return false, err
	}
	//synchronization.Sg.Synchronize(w.LoginUserID, sdkstruct.SvrConf.Platform)
	return true, nil
}

func (w *WalletMgr) EncrptSeedPhrase(key, seedPhrase string) error {
	userID := w.LoginUserID

	// Check userID is valid
	if len(userID) == 0 {
		return errors.New("userID is empty")
	}

	userExist, err := w.Db.GetUserByUserID(userID)
	if err != nil {
		return fmt.Errorf("error in GetUserByUserID() %w", err)
	}

	if userExist.Status != 1 {
		return errors.New("seed phrase has been generated before")
	}

	encSeedPhrase, err := utils.EncryptAES(seedPhrase, key)
	if err != nil {
		log.Println("Error in EncryptAES()")
	}

	user := model_struct.LocalUser{
		UserID:       userID,
		EntropyLevel: 128,
		SeedPhrase:   string(encSeedPhrase),
		Status:       2, // Status 2: Seed phrase created. but not confirmed by user
	}

	err = w.Db.UpdateLocalUser(&user)
	if err != nil {
		return fmt.Errorf("error in UpdateLocalUser() %w", err)
	}
	return nil

}
func (w *WalletMgr) GetLocalWallet2(coinType int) string {

	j, _ := w.Db.GetLocalWalletByUserID("08412", coinType)
	fmt.Println("Je", j.WalletImportFormat)
	return j.WalletImportFormat
}

func (w *WalletMgr) AddAddressBook(coinType int, address string, name string) (bool, error) {

	userID := w.LoginUserID
	//USDT-ERC20 and ETH use same address book, Trx and USDT-TRC20 use same address book
	if coinType == constant.USDTERC20 {
		coinType = constant.ETHCoin
	} else if coinType == constant.USDTTRC20 {
		coinType = constant.TRX
	}

	book := model_struct.LocalUserAddressBook{
		UserID:     userID,
		Name:       name,
		Address:    address,
		CoinType:   uint8(coinType),
		Status:     1,
		CreateTime: time.Now(),
	}

	model, err := w.Db.GetLocalAddressByAddress(userID, coinType, address)
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		log.Println(fmt.Errorf("error in GetLocalAddressByAddress() %w", err))
		return false, err
	}

	if model.UserID != "" {
		log.Println(fmt.Errorf("address already exist"))
		return false, errors.Wrap(constant.ErrAddressAlreadyExists, fmt.Sprintf("%v", constant.ErrAddressAlreadyExists.ErrCode))
	}

	if coinType == constant.ETHCoin || coinType == constant.USDTERC20 {
		if !utils.IsValidAddress(address) || utils.IsZeroAddress(address) {
			log.Println(fmt.Errorf("The address is incorrect, please input again"))
			return false, errors.Wrap(constant.ErrIncorrectAddress, fmt.Sprintf("%v", constant.ErrIncorrectAddress.ErrCode))
		}
	} else if coinType == constant.TRX || coinType == constant.USDTTRC20 {
		if !utils.IsValidTRXAddress(address) {
			log.Println(fmt.Errorf("The address is incorrect, please input again"))
			return false, errors.Wrap(constant.ErrIncorrectAddress, fmt.Sprintf("%v", constant.ErrIncorrectAddress.ErrCode))
		}
	}

	err = w.Db.InsertLocalUserAddressBook(&book)
	if err != nil {
		log.Println(fmt.Errorf("error in InsertLocalUserAddressBook() %w", err))
		return false, err
	}
	return true, nil

}
func (w *WalletMgr) GetAddressBook(coinType int) string {

	userID := w.LoginUserID
	if userID == "" {
		log.Println(errors.New("userID is empty"))
		return ""
	}
	book, err := w.Db.GetLocalUserAddressBook(userID, coinType)
	if err != nil {
		log.Println(fmt.Errorf("error in GetLocalUserAddressBook() %w", err))
		return ""
	}

	res := []sdkstruct.GetAddressBook{}
	for _, v := range book {
		book := sdkstruct.GetAddressBook{}
		book.Address = v.Address
		book.Name = v.Name
		book.CoinType = int(v.CoinType)
		res = append(res, book)
	}
	return utils.StructToJsonString(res)
}
func (w *WalletMgr) GetUserPrivateKey(secret string, coinType int, keyType int) string {

	userID := w.LoginUserID
	if userID == "" {
		log.Println(errors.New("userID is empty"))
		return ""
	}

	key := fmt.Sprintf("%s%s", secret, w.UserRandomSecret)

	wallet, err := w.Db.GetLocalWalletByUserID(userID, coinType)
	if err != nil {
		log.Println(fmt.Errorf("error in GetUserByUserID() %w", err))
		return ""
	}
	if wallet.WalletImportFormat == "" {
		log.Println("wallet.WalletImportFormat is nil")
		return ""
	}
	privatekey, err := utils.DecryptAES(wallet.WalletImportFormat, key)
	if err != nil {
		log.Println(fmt.Errorf("error in DecryptAES() %w", err))
		return ""
	}
	// ToDo: Hex/Com/UnCompressed
	if keyType == 1 {
		if coinType == constant.ETHCoin || coinType == constant.USDTERC20 {
			return privatekey
		} else if coinType == constant.TRX || coinType == constant.USDTTRC20 {
			return strings.TrimPrefix(privatekey, "0x")
		} else {
			return privatekey
		}

	} else if keyType == 2 {
		return privatekey
	} else if keyType == 3 {
		return privatekey
	} else {
		return privatekey
	}
}
func (w *WalletMgr) DeleteAddressBook(coinType int, address string) bool {
	userID := w.LoginUserID
	err := w.Db.DeleteAddressBook(coinType, address, userID)
	if err != nil {
		log.Println(fmt.Errorf("error in DeleteAddressBook() %w", err))
		return false
	}
	return true
}

func (w *WalletMgr) GetAddressBookByAddress(coinType int, address string) string {

	userID := w.LoginUserID
	if userID == "" {
		log.Println(errors.New("userID is empty"))
		return ""
	}
	book, err := w.Db.GetLocalUserAddressBookbyAddress(userID, coinType, address)
	if err != nil {
		log.Println(fmt.Errorf("error in GetLocalUserAddressBook() %w", err))
		return ""
	}
	res := sdkstruct.GetAddressBook{}
	res.Address = book.Address
	res.Name = book.Name
	return utils.StructToJsonString(res)
}
