package repo

import (
	"log"
	"os"

	"github.com/AlirezaAK2000/online-shop/constants"
	"github.com/AlirezaAK2000/online-shop/initializers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type Product struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Name          string             `bson:"name,omitempty"`
	Description   string             `bson:"description,omitempty"`
	Price         float64            `bson:"price,omitempty"`
	StockQuantity int32              `bson:"stock_quantity,omitempty"`
	Category      string             `bson:"category,omitempty" validate:"gte=0"`
}

func (p *Product) InsertProduct() (*mongo.InsertOneResult, error) {

	coll := initializers.Client.Database(os.Getenv("MONGODB_DB_NAME")).Collection(constants.ProductCollectionName)
	result, err := coll.InsertOne(context.Background(), p)
	return result, err

}

func FindProductByID(id string) (*Product, error) {

	coll := initializers.Client.Database(os.Getenv("MONGODB_DB_NAME")).Collection(constants.ProductCollectionName)
	filter := bson.D{{"_id", id}}
	var p Product
	err := coll.FindOne(context.Background(), filter).Decode(&p)
	return &p, err
}

func DeleteProductByID(id string) (*mongo.DeleteResult, error) {

	coll := initializers.Client.Database(os.Getenv("MONGODB_DB_NAME")).Collection(constants.ProductCollectionName)
	filter := bson.D{{"_id", id}}
	result, err := coll.DeleteOne(context.Background(), filter)
	return result, err

}

func FindAllProducts() (*[]Product, error) {

	coll := initializers.Client.Database(os.Getenv("MONGODB_DB_NAME")).Collection(constants.ProductCollectionName)
	filter := bson.D{{}}
	cursor, err := coll.Find(context.Background(), filter)
	var results []Product
	if err1 := cursor.All(context.Background(), &results); err1 != nil {
		log.Fatal(err1)
		err = err1
	}
	return &results, err

}
