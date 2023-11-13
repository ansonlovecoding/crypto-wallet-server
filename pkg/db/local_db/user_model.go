package db

import (
	"Share-Wallet/pkg/common/constant"
	"Share-Wallet/pkg/db/local_db/model_struct"
	"Share-Wallet/pkg/utils"
	"errors"
	"time"
)

func (d *DataBase) InsertLocalUser(user *model_struct.LocalUser) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Create(user).Error, "InsertLoginUser failed")
}

func (d *DataBase) InsertLocalWallet(wallet *model_struct.LocalWallet) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Create(wallet).Error, "InsertLoginUser failed")
}
func (d *DataBase) GetUser() (*model_struct.LocalUser, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var user model_struct.LocalUser
	return &user, utils.Wrap(d.conn.First(&user).Error, "GetLoginUserInfo failed")
}

func (d *DataBase) UpdateLocalUser(user *model_struct.LocalUser) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.Model(&user).Where("user_id = ?", user.UserID).Select("EntropyLevel", "SeedPhrase", "Status").Updates(model_struct.LocalUser{EntropyLevel: user.EntropyLevel, SeedPhrase: user.SeedPhrase, Status: user.Status})
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateLocalUser failed")
}
func (d *DataBase) UpdateLocalUserStatus(user *model_struct.LocalUser) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.Model(&user).Where("user_id = ?", user.UserID).Select("Status").Updates(model_struct.LocalUser{Status: user.Status})
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateLocalUserStatus failed")
}
func (d *DataBase) UpdateLocalUserPublicKeyAndStatus(user *model_struct.LocalUser) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.Model(&user).Where("user_id = ?", user.UserID).Select("Status", "PublicKey").Updates(user)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateLocalUserStatus failed")
}

func (d *DataBase) GetUserByUserID(userID string) (*model_struct.LocalUser, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var user model_struct.LocalUser
	return &user, utils.Wrap(d.conn.First(&user).Where("user_id = ?", userID).Error, "GetUserByUserID failed")
}

func (d *DataBase) UpdateLocalWallet(wallet *model_struct.LocalWallet) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.Model(&wallet).
		Where("user_id = ? and coin_type = ?", wallet.UserID, wallet.CoinType).
		Select("WalletAddress", "Status", "UpdateTime").
		Updates(model_struct.LocalWallet{WalletAddress: wallet.WalletAddress, Status: wallet.Status, UpdateTime: wallet.UpdateTime})
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateLocalWallet failed")
}

func (d *DataBase) UpdateLocalWalletAddress(wallet *model_struct.LocalWallet) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.Model(&wallet).
		Where("user_id = ? and coin_type = ?", wallet.UserID, wallet.CoinType).
		Select("Account", "P2PKHAddress", "P2SHSegwitAddress", "Bech32Address", "FullPublicKey", "WalletImportFormat", "Idx", "AddrStatus", "UpdateTime").
		Updates(model_struct.LocalWallet{Account: wallet.Account,
			P2PKHAddress:       wallet.P2PKHAddress,
			P2SHSegwitAddress:  wallet.P2SHSegwitAddress,
			Bech32Address:      wallet.Bech32Address,
			FullPublicKey:      wallet.FullPublicKey,
			WalletImportFormat: wallet.WalletImportFormat,
			Idx:                wallet.Idx,
			AddrStatus:         wallet.AddrStatus,
			UpdateTime:         time.Now(),
		})
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateLocalWallet failed")
}
func (d *DataBase) GetLocalWalletByUserID(userID string, coinType int) (*model_struct.LocalWallet, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var wallet model_struct.LocalWallet
	err := d.conn.Model(&wallet).Where("user_id = ? and coin_type = ?", userID, coinType).Find(&wallet).Error
	return &wallet, utils.Wrap(err, "GetLocalWalletByUserID failed")
}
func (d *DataBase) InsertLocalUserAddressBook(book *model_struct.LocalUserAddressBook) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Create(book).Error, "InsertLocalUserAddressBook failed")
}
func (d *DataBase) GetLocalUserAddressBook(userID string, coinType int) ([]*model_struct.LocalUserAddressBook, error) {

	var book []model_struct.LocalUserAddressBook
	var err error
	if coinType > 0 {
		//USDT-ERC20 and ETH use same address book, Trx and USDT-TRC20 use same address book
		if coinType == constant.USDTERC20 {
			coinType = constant.ETHCoin
		} else if coinType == constant.USDTTRC20 {
			coinType = constant.TRX
		}
		err = utils.Wrap(d.conn.Order("id desc").Where("status = ? and user_id = ? and coin_type = ?", 1, userID, coinType).Find(&book).Error, "GetLocalUserAddressBook() failed")
	} else {
		err = utils.Wrap(d.conn.Order("id desc").Where("status = ? and user_id = ?", 1, userID).Find(&book).Error, "GetLocalUserAddressBook() failed")
	}

	if err != nil {
		return nil, utils.Wrap(err, "GetLocalUserAddressBook() failed")
	}
	var list []*model_struct.LocalUserAddressBook
	for _, v := range book {
		v1 := v
		list = append(list, &v1)
	}
	return list, nil
}
func (d *DataBase) GetLocalAddressByAddress(userID string, coinType int, address string) (*model_struct.LocalUserAddressBook, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var book model_struct.LocalUserAddressBook
	//USDT-ERC20 and ETH use same address book, Trx and USDT-TRC20 use same address book
	if coinType == constant.USDTERC20 {
		coinType = constant.ETHCoin
	} else if coinType == constant.USDTTRC20 {
		coinType = constant.TRX
	}
	return &book, d.conn.Order("id asc").Where("user_id = ? and coin_type = ? and address = ?", userID, coinType, address).Find(&book).Error
}
func (d *DataBase) DeleteAddressBook(coinType int, address string, userID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var book model_struct.LocalUserAddressBook
	return utils.Wrap(d.conn.Where("coin_type = ? and address = ? and user_id = ?", coinType, address, userID).Delete(&book).Error, "DeleteAddressBook failed")
}
func (d *DataBase) GetLocalUserAddressBookbyAddress(userID string, coinType int, address string) (*model_struct.LocalUserAddressBook, error) {

	var book model_struct.LocalUserAddressBook
	//USDT-ERC20 and ETH use same address book, Trx and USDT-TRC20 use same address book
	if coinType == constant.USDTERC20 {
		coinType = constant.ETHCoin
	} else if coinType == constant.USDTTRC20 {
		coinType = constant.TRX
	}
	err := utils.Wrap(d.conn.Order("id asc").Where("status = ? and user_id = ? and coin_type = ? and address = ?", 1, userID, coinType, address).Find(&book).Error, "GetLocalUserAddressBook() failed")
	if err != nil {
		return nil, utils.Wrap(err, "GetLocalUserAddressBookbyAddress() failed")
	}

	return &book, nil
}

func (d *DataBase) GetLocalUserAddresses(userID string) ([]model_struct.LocalWallet, error) {
	var wallet []model_struct.LocalWallet
	err := utils.Wrap(d.conn.Table("local_wallets").Debug().Where("user_id = ?", userID).Select("coin_type, p2pkh_address").Scan(&wallet).Error, "GetLocalUserAddresses() failed")
	if err != nil {
		return nil, utils.Wrap(err, "GetLocalUserAddresses() failed")
	}
	return wallet, err
}
func (d *DataBase) UpdateLocalUserWalletID(user *model_struct.LocalUser) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.Model(&user).Where("user_id = ?", user.UserID).Select("walletID").Updates(model_struct.LocalUser{WalletID: user.WalletID})
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateLocalUserWalletID failed")
}

func (d *DataBase) IsTheAddressExist(userID string, coinType int, address string) (bool, error) {
	var book model_struct.LocalUserAddressBook
	err := d.conn.Order("id asc").Where("user_id = ? and coin_type = ? and address = ?", userID, coinType, address).Find(&book).Error
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (d *DataBase) UpdateLocalWalletTypes(wallet *model_struct.LocalWalletType, coinType string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.Model(&wallet).Where("coin_name = ?", coinType).Select("status").Updates(model_struct.LocalWalletType{Status: wallet.Status})
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateLocalUserWalletID failed")
}
