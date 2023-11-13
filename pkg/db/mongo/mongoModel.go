package db

import (
	"Share-Wallet/pkg/common/config"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"strings"

	"time"
)

const CSysLog = "sys_log"

type MongoDB struct {
	mongo *mongo.Client
}

type Log struct {
	UserID     string `bson:"user_id"`
	Content    string `bson:"content"`
	CreateTime int64  `bson:"create_time"`
}

func (d *MongoDB) InitMongoDB() {
	// mongo init
	// "mongodb://sysop:moon@localhost/records"
	uri := "mongodb://sample.host:27017/?maxPoolSize=20&w=majority"
	if config.Config.Mongo.DBUri != "" {
		// example: mongodb://$user:$password@mongo1.mongo:27017,mongo2.mongo:27017,mongo3.mongo:27017/$DBDatabase/?replicaSet=rs0&readPreference=secondary&authSource=admin&maxPoolSize=$DBMaxPoolSize
		uri = config.Config.Mongo.DBUri
	} else {
		if config.Config.Mongo.DBPassword != "" && config.Config.Mongo.DBUserName != "" {
			uri = fmt.Sprintf("mongodb://%s:%s@%s/%s?maxPoolSize=%d", config.Config.Mongo.DBUserName, config.Config.Mongo.DBPassword, config.Config.Mongo.DBAddress,
				config.Config.Mongo.DBDatabase, config.Config.Mongo.DBMaxPoolSize)
		} else {
			uri = fmt.Sprintf("mongodb://%s/%s/?maxPoolSize=%d",
				config.Config.Mongo.DBAddress, config.Config.Mongo.DBDatabase,
				config.Config.Mongo.DBMaxPoolSize)
		}
	}
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println(" mongo.Connect  failed, try ", err.Error(), uri)
		time.Sleep(time.Duration(30) * time.Second)
		mongoClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
		if err != nil {
			fmt.Println(" mongo.Connect retry failed, panic", err.Error(), uri)
			panic(err.Error())
		}
	}
	fmt.Println("0", "mongo driver client init success: ", uri)
	// mongodb create index
	if err := createMongoIndex(mongoClient, CSysLog, false, "send_id", "-send_time"); err != nil {
		fmt.Println("send_id", "-send_time", "index create failed", err.Error())
	}
	fmt.Println("create index success")
	d.mongo = mongoClient
}

func (d *MongoDB) SaveLog(tagSendLog *Log) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongo.Database(config.Config.Mongo.DBDatabase).Collection(CSysLog)
	_, err := c.InsertOne(ctx, tagSendLog)
	return err
}

func (d *MongoDB) GetTagSendLogs(userID string, showNumber, pageNumber int32) ([]Log, error) {
	var tagSendLogs []Log
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongo.Database(config.Config.Mongo.DBDatabase).Collection(CSysLog)
	findOpts := options.Find().SetLimit(int64(showNumber)).SetSkip(int64(showNumber) * (int64(pageNumber) - 1)).SetSort(bson.M{"send_time": -1})
	cursor, err := c.Find(ctx, bson.M{"send_id": userID}, findOpts)
	if err != nil {
		return tagSendLogs, err
	}
	err = cursor.All(ctx, &tagSendLogs)
	if err != nil {
		return tagSendLogs, err
	}
	return tagSendLogs, nil
}

func createMongoIndex(client *mongo.Client, collection string, isUnique bool, keys ...string) error {
	db := client.Database(config.Config.Mongo.DBDatabase).Collection(collection)
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	indexView := db.Indexes()
	keysDoc := bsonx.Doc{}

	// 复合索引
	for _, key := range keys {
		if strings.HasPrefix(key, "-") {
			keysDoc = keysDoc.Append(strings.TrimLeft(key, "-"), bsonx.Int32(-1))
		} else {
			keysDoc = keysDoc.Append(key, bsonx.Int32(1))
		}
	}

	// 创建索引
	index := mongo.IndexModel{
		Keys: keysDoc,
	}
	if isUnique == true {
		index.Options = options.Index().SetUnique(true)
	}
	_, err := indexView.CreateOne(
		context.Background(),
		index,
		opts,
	)
	if err != nil {
		return err
	}
	return nil
}
