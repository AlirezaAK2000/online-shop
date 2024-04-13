package initializers

import (
	"context"
	"log"
	"os"

	"github.com/AlirezaAK2000/online-shop/constants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func InitializeMongoConnection() {

	if Client == nil {
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		opts := options.Client().ApplyURI(os.Getenv("MONGODB_URL")).SetServerAPIOptions(serverAPI)
		// Create a new client and connect to the server
		var err error
		Client, err = mongo.Connect(context.Background(), opts)

		if err != nil {
			panic(err)
		}

		var result bson.M
		if err := Client.Database("admin").RunCommand(context.Background(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
			panic(err)
		}
		log.Println("Pinged your deployment. You successfully connected to MongoDB!")

		userInitializer()
		productInitializer()
	}
}

func DisconnectMongoDBClient() {
	if err := Client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func userInitializer() {

	coll := Client.Database(os.Getenv("MONGODB_DB_NAME")).Collection(constants.UserCollectionName)
	ctx := context.Background()
	// Define the index options
	indexOptions := options.Index().SetUnique(true)

	// Define the index model
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"Email", 1}}, // Change "email" to the field you want to enforce uniqueness on
		Options: indexOptions,
	}

	// Create the index
	_, err := coll.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		panic(err)
	}
}

func productInitializer() {
	Client.Database(os.Getenv("MONGODB_DB_NAME")).Collection(constants.ProductCollectionName)
}
