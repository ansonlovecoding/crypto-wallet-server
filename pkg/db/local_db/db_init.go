package db

import (
	"Share-Wallet/pkg/common/constant"
	"Share-Wallet/pkg/common/log"
	"Share-Wallet/pkg/db/local_db/model_struct"
	"Share-Wallet/pkg/utils"
	"errors"
	"os"
	"sync"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

const (
	LocalWalletTableName = "local_wallet_types"
)

var UserDBMap map[string]*DataBase

type DataBase struct {
	loginUserID string
	dbDir       string
	conn        *gorm.DB
	mRWMutex    sync.RWMutex
}

var UserDBLock sync.RWMutex

func init() {
	UserDBMap = make(map[string]*DataBase, 0)
}

func NewDataBase(loginUserID string, dbDir string) (*DataBase, error) {

	UserDBLock.Lock()
	defer UserDBLock.Unlock()
	dataBase, ok := UserDBMap[loginUserID]
	if !ok {
		dataBase = &DataBase{loginUserID: loginUserID, dbDir: dbDir}
		err := dataBase.initDB()
		if err != nil {
			return dataBase, utils.Wrap(err, "initDB failed")
		}
		UserDBMap[loginUserID] = dataBase
		log.Info("", "open db", loginUserID)
	} else {
		log.Info("", "db in map", loginUserID)
	}

	return dataBase, nil
}
func (d *DataBase) initDB() error {

	if d.loginUserID == "" {
		return errors.New("no uid")
	}
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()

	dbFileName := d.dbDir + "/Wallet_" + constant.BigVersion + "_" + d.loginUserID + ".db"
	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{Logger: log.GetSqlLogger(constant.SQLiteLogFileName)})
	db.Logger.LogMode(gormlogger.Silent)

	log.Info("open db:", dbFileName)
	if err != nil {
		return utils.Wrap(err, "open db failed")
	}
	sqlDB, err := db.DB()

	if err != nil {
		return utils.Wrap(err, "get sql db failed")
	}
	sqlDB.SetConnMaxLifetime(time.Hour * 1)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(2)
	d.conn = db

	db.AutoMigrate(&model_struct.LocalUser{},
		&model_struct.LocalWallet{},
		&model_struct.LocalWalletType{},
		&model_struct.LocalUserAddressBook{})

	if db.Migrator().HasTable(&model_struct.LocalWalletType{}) {
		// local_wallet_type table is initialized with coins mentioned in the prototype.
		// Status of the respective coins can get it from MYSQL server later.
		err := d.InitLocalWalletCoinType()
		if err != nil {
			return err
		}
	}
	return nil
}

func RemoveAllLocalDatabases(dbDir string) error {
	err := os.RemoveAll(dbDir + "/")
	return err
}
