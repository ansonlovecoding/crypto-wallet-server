package db

import (
	"Share-Wallet/pkg/common/config"
	"Share-Wallet/pkg/common/constant"
	"Share-Wallet/pkg/common/log"
	"context"
	"fmt"
	"strconv"
	"time"

	go_redis "github.com/go-redis/redis/v8"
)

const (
	uidPidToken          = "UID_PID_TOKEN_STATUS:"
	OrderLockKey         = "ORDER_LOCK:"
	UnConfirmedOrderList = "UNCONFIRMED_ORDERLIST:"
	UpdateAddressesList  = "UPDATE_ADDRESSES:"
)

type RedisDB struct {
	redisClient go_redis.UniversalClient
}

func (d *RedisDB) InitRedisDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if config.Config.Redis.EnableCluster {
		d.redisClient = go_redis.NewClusterClient(&go_redis.ClusterOptions{
			Addrs:    config.Config.Redis.DBAddress,
			PoolSize: 50,
		})
		_, err := d.redisClient.Ping(ctx).Result()
		if err != nil {
			panic(err.Error())
		}
	} else {
		d.redisClient = go_redis.NewClient(&go_redis.Options{
			Addr:     config.Config.Redis.DBAddress[0],
			Password: config.Config.Redis.DBPassWord, // no password set
			DB:       0,                              // use default DB
			PoolSize: 100,                            // 连接池大小
		})
		_, err := d.redisClient.Ping(ctx).Result()
		if err != nil {
			panic(err.Error())
		}
	}
}

// Store userid and platform class to redis
func (d *RedisDB) AddTokenFlag(userID string, platformID int, token string, flag int) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	log.NewDebug("", "add token key is ", key)
	return d.redisClient.HSet(context.Background(), key, token, flag).Err()
}

func (d *RedisDB) GetTokenMapByUidPid(userID, platformID string) (map[string]interface{}, error) {
	key := uidPidToken + userID + ":" + platformID
	log.NewDebug("", "get token key is ", key)
	m, err := d.redisClient.HGetAll(context.Background(), key).Result()
	mm := make(map[string]interface{})
	for k, v := range m {
		if j, ok := strconv.Atoi(v); ok == nil {
			mm[k] = j
		}
	}
	return mm, err
}
func (d *RedisDB) SetTokenMapByUidPid(userID string, platformID int, m map[string]interface{}) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	//for k, v := range m {
	//	err := d.rdb.HSet(context.Background(), key, k, v).Err()
	//	if err != nil {
	//		return err
	//	}
	//}
	//return nil
	return d.redisClient.HMSet(context.Background(), key, m).Err()
}
func (d *RedisDB) DeleteTokenByUidPid(userID string, platformID int, fields []string) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	return d.redisClient.HDel(context.Background(), key, fields...).Err()
}

func (d *RedisDB) DeleteTokenByUid(userID string, platformID int) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	return d.redisClient.Del(context.Background(), key).Err()
}

func (d *RedisDB) LockOrder(txid string, coinType uint8, seconds int) error {
	key := OrderLockKey + txid + ":" + getCoinName(coinType)
	return d.redisClient.SetNX(context.Background(), key, "LOCK", time.Duration(seconds)*time.Second).Err()
}

func (d *RedisDB) UnLockOrder(txid string, coinType uint8) error {
	key := OrderLockKey + txid + ":" + getCoinName(coinType)
	return d.redisClient.Del(context.Background(), key).Err()
}

func (d *RedisDB) InsertUnconfirmedOrder(txid string, coinType uint8) error {
	key := UnConfirmedOrderList + getCoinName(coinType)
	return d.redisClient.XAdd(context.Background(), &go_redis.XAddArgs{
		Stream:     key,
		NoMkStream: false,
		MaxLen:     100000000,
		Approx:     false,
		ID:         "",
		Values:     []interface{}{"message", txid},
	}).Err()
}

func (d *RedisDB) ReadUnconfirmedOrder(coinType uint8, num int64) ([]string, error) {
	key := UnConfirmedOrderList + getCoinName(coinType)
	xStreams, err := d.redisClient.XRead(context.Background(), &go_redis.XReadArgs{
		Streams: []string{key, "0-0"},
		Count:   num,
		Block:   time.Second * 10,
	}).Result()
	if err != nil {
		fmt.Println("XRead failed", err.Error())
		return nil, err
	}
	var txList []string
	if xStreams == nil {
		fmt.Println("XRead get empty stream")
		return txList, nil
	}
	for _, stream := range xStreams {
		messages := stream.Messages
		for _, message := range messages {
			for _, v := range message.Values {
				//fmt.Println(message.ID, v.(string))
				if v != "message" {
					txList = append(txList, v.(string)+"|"+message.ID)
				}
			}
		}
	}
	return txList, nil
}

func (d *RedisDB) RemoveUnconfirmedOrder(msgID string, coinType uint8) error {
	key := UnConfirmedOrderList + getCoinName(coinType)
	return d.redisClient.XDel(context.Background(), key, msgID).Err()
}

func getCoinName(coinType uint8) string {
	switch coinType {
	case constant.BTCCoin:
		return "BTC"
	case constant.ETHCoin:
		return "ETH"
	case constant.USDTERC20:
		return "USDT-ERC20"
	case constant.TRX:
		return "TRX"
	case constant.USDTTRC20:
		return "USDT-TRC20"
	default:
		return ""
	}
}

func (d *RedisDB) InsertAddressToUpdate(address string) error {
	key := UpdateAddressesList //+ getCoinName(coinType)

	return d.redisClient.XAdd(context.Background(), &go_redis.XAddArgs{
		Stream:     key,
		NoMkStream: false,
		MaxLen:     100000000,
		Approx:     false,
		ID:         "",
		Values:     []interface{}{"message", address},
	}).Err()
}
func (d *RedisDB) ReadUpdateAddresses(num int64) ([]string, error) {
	key := UpdateAddressesList //+ getCoinName(coinType)
	xStreams, err := d.redisClient.XRead(context.Background(), &go_redis.XReadArgs{
		Streams: []string{key, "0-0"},
		Count:   num,
		Block:   time.Second * 10,
	}).Result()
	if err != nil {
		fmt.Println("Update addresses XRead failed", err.Error(), key)
		return nil, err
	}
	var addressList []string
	for _, stream := range xStreams {
		messages := stream.Messages
		for _, message := range messages {
			for _, v := range message.Values {
				if v != "message" {
					addressList = append(addressList, v.(string)+"|"+message.ID)
				}
			}
		}
	}
	return addressList, nil
}
func (d *RedisDB) RemoveAddressFromList(msgID string) error {
	key := UpdateAddressesList //+ getCoinName(coinType)
	return d.redisClient.XDel(context.Background(), key, msgID).Err()
}
