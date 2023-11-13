package db

import (
	"Share-Wallet/pkg/common/config"
	db2 "Share-Wallet/pkg/db/mongo"
	db "Share-Wallet/pkg/db/mysql"
	db3 "Share-Wallet/pkg/db/redis"

	go_redis "github.com/go-redis/redis/v8"
	"gopkg.in/mgo.v2"
)

var DB DataBases

type DataBases struct {
	MysqlDB    db.MysqlDB
	MgoSession *mgo.Session
	//redisPool   *redis.Pool
	MongoDB db2.MongoDB
	RedisDB db3.RedisDB
}

type RedisClient struct {
	client  *go_redis.Client
	cluster *go_redis.ClusterClient
	go_redis.UniversalClient
	enableCluster bool
}

func init() {
	if config.Config.Mysql.DBAddress == nil || config.Config.IsSkipDatabase {
		return
	}

	//mysql init
	db.InitMysqlDB()
	//mongo init
	// DB.MongoDB.InitMongoDB()
	//redis init
	DB.RedisDB.InitRedisDB()

}
