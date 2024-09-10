package db

import (
	"context"
	"fmt"
	"sync"

	"github.com/nahuelojea/handballscore/dto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient  *mongo.Client
	DatabaseName string
	once         sync.Once
)

func Connect(ctx context.Context) error {
	var err error
	once.Do(func() {
		user := ctx.Value(dto.Key("user")).(string)
		password := ctx.Value(dto.Key("password")).(string)
		host := ctx.Value(dto.Key("host")).(string)
		connStr := fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority", user, password, host)

		clientOptions := options.Client().ApplyURI(connStr)
		MongoClient, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			fmt.Println("Error al conectar a MongoDB:", err)
			return
		}

		err = MongoClient.Ping(ctx, nil)
		if err != nil {
			fmt.Println("Error al hacer ping a MongoDB:", err)
			return
		}

		fmt.Println("Conexión exitosa a la base de datos")
		DatabaseName = ctx.Value(dto.Key("database")).(string)
	})

	return err
}
