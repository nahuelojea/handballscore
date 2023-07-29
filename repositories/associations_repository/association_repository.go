package associations_repository

import (
	"context"

	"github.com/nahuelojea/handballscore/config/db"
	"github.com/nahuelojea/handballscore/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	association_collection = "associations"
)

func GetAssociation(ID string) (models.Association, bool, error) {
	ctx := context.TODO()
	db := db.MongoClient.Database(db.DatabaseName)
	collection := db.Collection(association_collection)

	var association models.Association
	objId, _ := primitive.ObjectIDFromHex(ID)

	condicion := bson.M{
		"_id": objId,
	}

	err := collection.FindOne(ctx, condicion).Decode(&association)
	if err != nil {
		return association, false, err
	}

	return association, true, nil
}
