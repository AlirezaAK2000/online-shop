package repo

import (
	"log"
	"os"

	"github.com/AlirezaAK2000/online-shop/constants"
	"github.com/AlirezaAK2000/online-shop/initializers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type User struct {
	Email    string `bson:"_id,omitempty" validate:"required,email"`
	Password string `bson:"password,omitempty" validate:"required"`
	Address  string `bson:"address,omitempty"`
	Phone    string `bson:"phone,omitempty"`
	Type     string `bson:"type,omitempty"`
}

func (u *User) InsertUser() (*mongo.InsertOneResult, error) {

	coll := initializers.Client.Database(os.Getenv("MONGODB_DB_NAME")).Collection(constants.UserCollectionName)
	result, err := coll.InsertOne(context.Background(), u)
	return result, err

}

func FindUserByID(id string) (*User, error) {

	coll := initializers.Client.Database(os.Getenv("MONGODB_DB_NAME")).Collection(constants.UserCollectionName)
	filter := bson.D{{"_id", id}}
	var u User
	err := coll.FindOne(context.Background(), filter).Decode(&u)
	return &u, err
}

func FindUserByEmail(email string) (*User, error) {

	coll := initializers.Client.Database(os.Getenv("MONGODB_DB_NAME")).Collection(constants.UserCollectionName)
	filter := bson.D{{"_id", email}}
	var u User
	err := coll.FindOne(context.Background(), filter).Decode(&u)
	return &u, err
}

func DeleteUserByID(id string) (*mongo.DeleteResult, error) {

	coll := initializers.Client.Database(os.Getenv("MONGODB_DB_NAME")).Collection(constants.UserCollectionName)
	filter := bson.D{{"_id", id}}
	result, err := coll.DeleteOne(context.Background(), filter)
	return result, err

}

func FindAllUsers() (*[]User, error) {

	coll := initializers.Client.Database(os.Getenv("MONGODB_DB_NAME")).Collection(constants.UserCollectionName)
	filter := bson.D{{}}
	cursor, err := coll.Find(context.Background(), filter)
	var results []User
	if err1 := cursor.All(context.Background(), &results); err1 != nil {
		log.Fatal(err1)
		err = err1
	}
	return &results, err

}
