package database

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"net/url"
	"os"
)

type MongoDb struct{}

var client *mongo.Client
var db *mongo.Database

var initialized = false

func Initialize() {

	cluster := os.Getenv("MONGO_URL") // "localhost:27017"
	username := os.Getenv("MONGO_USER")
	password := os.Getenv("MONGO_PASS")

	if cluster == "" {
		fmt.Println("mongodb uri not set. run without db connection")
		return
	}

	uri := "mongodb://" + url.QueryEscape(username) + ":" + url.QueryEscape(password) + "@" + cluster + "/admin"

	fmt.Println(uri)
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	db = client.Database("slot_data_1")
	//userCollection = client.Database("slot_data_2").Collection("user")
	//spinCollection = client.Database("slot_data_2").Collection("spin")
	//collectCollection = client.Database("slot_data_2").Collection("collect")
	initialized = true
}

func DisConnect() {
	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func InsertSpinData(symbol string, spinRawData map[string]interface{}) {
	if initialized == false {
		return
	}
	spinCollection := db.Collection(symbol + "_spin")
	_, err := spinCollection.InsertOne(context.TODO(), spinRawData)
	if err != nil {
		panic(err)
	}
}

//func InsertCollectionData(collectRawData bson.M) {
//	_, err := collectCollection.InsertOne(context.TODO(), collectRawData)
//	if err != nil {
//		panic(err)
//	}
//}
//

func GetUserList(symbol string) []User {
	if initialized == false {
		return nil
	}
	spinCollection := db.Collection(symbol + "_spin")
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$uuid"},
			{"count", bson.D{{"$sum", 1}}},
		}}}

	cursor, err := spinCollection.Aggregate(context.TODO(), mongo.Pipeline{groupStage})
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.Background())
	//
	var users []User
	if err = cursor.All(context.TODO(), &users); err != nil {
		panic(err)
	}
	return users
}

func GetSpinDataByUser(userId string, symbol string) ([]bson.M, error) {
	if initialized == false {
		return nil, nil
	}
	spinCollection := db.Collection(symbol + "_spin")
	cursor, err := spinCollection.Find(context.TODO(), bson.D{
		{"uuid", userId},
	})
	if err != nil {
		return nil, err
	}

	var spinResult []bson.M
	if err = cursor.All(context.TODO(), &spinResult); err != nil {
		panic(err)
	}
	return spinResult, nil
}
