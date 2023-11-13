package db

import (
	"Share-Wallet/pkg/common/config"
	"Share-Wallet/pkg/common/constant"
	"Share-Wallet/pkg/common/log"
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type MysqlDB struct {
	sync.RWMutex
	dbMap map[string]*gorm.DB
}

func key(dbAddress, dbName string) string {
	return dbAddress + "_" + dbName
}

func InitMysqlDB() {
	// When there is no open IM database, connect to the mysql built-in database to create wallet database
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.Config.Mysql.DBUserName, config.Config.Mysql.DBPassword, config.Config.Mysql.DBAddress[0], "mysql")
	var db *gorm.DB
	var err1 error
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: log.GetSqlLogger(constant.MySQLLogFileName), NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		fmt.Println("0", "Open failed ", err.Error(), dsn)
		time.Sleep(time.Duration(30) * time.Second)
		db, err1 = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: log.GetSqlLogger(constant.MySQLLogFileName), NamingStrategy: schema.NamingStrategy{SingularTable: true}})
		if err1 != nil {
			fmt.Println("0", "Open failed ", err1.Error(), dsn)
			panic(err1.Error())
		}
	}

	// Check the database and table during initialization
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s default charset utf8mb4 COLLATE utf8mb4_unicode_ci;", config.Config.Mysql.DBDatabaseName)
	err = db.Exec(sql).Error
	if err != nil {
		fmt.Println("0", "Exec failed ", err.Error(), sql)
		panic(err.Error())
	}

	sqlDB, _ := db.DB()
	sqlDB.Close()

	dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.Config.Mysql.DBUserName, config.Config.Mysql.DBPassword, config.Config.Mysql.DBAddress[0], config.Config.Mysql.DBDatabaseName)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: log.GetSqlLogger(constant.MySQLLogFileName), NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		fmt.Println("0", "Open failed ", err.Error(), dsn)
		panic(err.Error())
	}

	sqlDB, _ = db.DB()
	fmt.Println("open db ok ", dsn)
	db.AutoMigrate(
		&UserAccounts{},
		&UserBalances{},
		&LocalCurrencyRates{},
		&WalletCurrency{},
		&Statistics{},
		&AdminUser{},
		&AdminRole{},
		&AdminActions{},
		&EthDetailTX{},
		&TronDetailTX{},
		&AccountInformation{},
		&FundsLog{},
		&CoinCurrencyValues{},
	)

	db.Set("gorm:table_options", "CHARSET=utf8mb4")
	db.Set("gorm:table_options", "collation=utf8mb4_unicode_ci")

	if !db.Migrator().HasTable(&UserAccounts{}) {
		fmt.Println("CreateTable UserAccounts")
		db.Migrator().CreateTable(&UserAccounts{})
	}
	if !db.Migrator().HasTable(&UserBalances{}) {
		fmt.Println("CreateTable UserBalances")
		db.Migrator().CreateTable(&UserBalances{})
	}
	if !db.Migrator().HasTable(&LocalCurrencyRates{}) {
		fmt.Println("CreateTable LocalCurrencyRates")
		db.Migrator().CreateTable(&LocalCurrencyRates{})
	}
	if !db.Migrator().HasTable(&WalletCurrency{}) {
		fmt.Println("CreateTable WalletCurrency")
		db.Migrator().CreateTable(&WalletCurrency{})
	}
	if !db.Migrator().HasTable(&Statistics{}) {
		fmt.Println("CreateTable Statistics")
		db.Migrator().CreateTable(&Statistics{})
	}
	if !db.Migrator().HasTable(&AdminUser{}) {
		fmt.Println("CreateTable AdminUser")
		db.Migrator().CreateTable(&AdminUser{})
	}
	if !db.Migrator().HasTable(&AdminRole{}) {
		fmt.Println("CreateTable AdminRole")
		db.Migrator().CreateTable(&AdminRole{})
	}
	if !db.Migrator().HasTable(&AdminActions{}) {
		fmt.Println("CreateTable AdminActions")
		db.Migrator().CreateTable(&AdminActions{})
	}
	if !db.Migrator().HasTable(&EthDetailTX{}) {
		fmt.Println("CreateTable EthDetailTX")
		db.Migrator().CreateTable(&EthDetailTX{})
	}
	if !db.Migrator().HasTable(&TronDetailTX{}) {
		fmt.Println("CreateTable TronDetailTX")
		db.Migrator().CreateTable(&TronDetailTX{})
	}
	if !db.Migrator().HasTable(&AccountInformation{}) {
		fmt.Println("CreateTable AccountInformation")
		db.Migrator().CreateTable(&AccountInformation{})
	}
	if !db.Migrator().HasTable(&FundsLog{}) {
		fmt.Println("CreateTable FundsLog")
		db.Migrator().CreateTable(&FundsLog{})
	}
	if !db.Migrator().HasTable(&CoinCurrencyValues{}) {
		fmt.Println("CreateTable CoinCurrencyValues")
		db.Migrator().CreateTable(&CoinCurrencyValues{})
	}

	sqlDB.Close()
	return
}

func (m *MysqlDB) DefaultGormDB() (*gorm.DB, error) {
	return m.GormDB(config.Config.Mysql.DBAddress[0], config.Config.Mysql.DBDatabaseName)
}

func (m *MysqlDB) GormDB(dbAddress, dbName string) (*gorm.DB, error) {
	m.Lock()
	defer m.Unlock()

	k := key(dbAddress, dbName)
	if _, ok := m.dbMap[k]; !ok {
		if err := m.open(dbAddress, dbName); err != nil {
			return nil, err
		}
	}
	return m.dbMap[k], nil
}

func (m *MysqlDB) open(dbAddress, dbName string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.Config.Mysql.DBUserName, config.Config.Mysql.DBPassword, dbAddress, dbName)
	// db, err := gorm.Open("mysql", dsn, &gorm.Config{Logger: log.GetNewLogger(constant.SQLiteLogFileName)})
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: log.GetSqlLogger(constant.MySQLLogFileName), NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		return err
	}

	sqlDB, _ := database.DB()
	sqlDB.SetMaxOpenConns(config.Config.Mysql.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(config.Config.Mysql.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.Config.Mysql.DBMaxLifeTime) * time.Second)

	if m.dbMap == nil {
		m.dbMap = make(map[string]*gorm.DB)
	}
	k := key(dbAddress, dbName)
	m.dbMap[k] = database
	return nil
}
