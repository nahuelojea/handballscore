package db

import (
	"context"
	"fmt"

	"github.com/nahuelojea/handballscore/dto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var DatabaseName string

func Connect(ctx context.Context) error {
	user := ctx.Value(dto.Key("user")).(string)
	password := ctx.Value(dto.Key("password")).(string)
	host := ctx.Value(dto.Key("host")).(string)
	connStr := fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority", user, password, host)

	var clientOptions = options.Client().ApplyURI(connStr)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Successful connection to database")
	MongoClient = client
	DatabaseName = ctx.Value(dto.Key("database")).(string)
	return nil
}

func isActive() bool {
	err := MongoClient.Ping(context.TODO(), nil)
	return err == nil
}
