package db_test

import (
	"Share-Wallet/pkg/db"
	"strings"
	"testing"
)

func TestRedisDB_InsertUnconfirmedOrder(t *testing.T) {
	db.DB.RedisDB.InitRedisDB()
	err := db.DB.RedisDB.InsertUnconfirmedOrder("0xd8ae905d75852516992cee29663607cc5c4a43e886b8a7a362e585ca280dfab8", 2)
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log("insert unconfirmedOrder success!")
}

func TestRedisDB_ReadUnconfirmedOrder(t *testing.T) {
	db.DB.RedisDB.InitRedisDB()
	txList, err := db.DB.RedisDB.ReadUnconfirmedOrder(2, 100)
	if err != nil {
		t.Log(err.Error())
		return
	}

	t.Log("read unconfirmedOrder success!", len(txList))

	for _, v := range txList {
		strList := strings.Split(v, "|")
		if len(strList) == 2 {
			err := db.DB.RedisDB.RemoveUnconfirmedOrder(strList[1], 2)
			if err != nil {
				t.Log(err.Error())
				return
			}
		}
	}
}
func TestRedisDB_InsertAddressToUpdate(t *testing.T) {
	db.DB.RedisDB.InitRedisDB()
	err := db.DB.RedisDB.InsertAddressToUpdate("0xd0609aa34a6b78d47540a566899694966d4c0489")
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log("insert address to update success!")
}

func TestRedisDB_ReadUpdateAddresses(t *testing.T) {
	db.DB.RedisDB.InitRedisDB()
	addresses, err := db.DB.RedisDB.ReadUpdateAddresses(1)
	if err != nil {
		t.Log(err.Error())
		return
	}

	t.Log("read addressseslist success!", len(addresses))

	for _, v := range addresses {
		t.Log("Address is ", v)
		strList := strings.Split(v, "|")
		if len(strList) == 2 {
			err := db.DB.RedisDB.RemoveAddressFromList(strList[1])
			if err != nil {
				t.Log(err.Error())
				return
			}
			// t.Log("public address is  ", strList[0])
		}
	}
}
