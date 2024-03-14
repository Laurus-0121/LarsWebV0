package dao

import (
	"context"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const username = "haizhi"
const password = "haizhi666"
const ip = "127.0.0.1"
const port = "27017"
const dbname = "LarsWebv0"

var (
	database *mongo.Database
)

func MongoDBSetup() {
	mongoUrl := fmt.Sprintf("mongodb://%v:%v@%v:%v/%v", username, password, ip, port, dbname)
	clientOptions := options.Client().ApplyURI(mongoUrl)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logger.Errorf("fail set up mongo: %v", err)
	}
	database = client.Database(dbname)
}
