package repos

import (
	"github.com/rpsoftech/bullion-server/src/env"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const bullionSiteInfoCollectionName = "BullionSiteInfo"

var BullionSiteInfoRepo *bullionSiteInfoRepo

func init() {
	if env.Env.APP_ENV == env.APP_ENV_DEVELOPE {
		return
	}
	coll := mongodb.MongoDatabase.Collection(bullionSiteInfoCollectionName)
	BullionSiteInfoRepo = &bullionSiteInfoRepo{
		collection: coll,
	}
	addUniqueIndexesToCollection([]string{"id", "domain", "shortName"}, BullionSiteInfoRepo.collection)
}

type bullionSiteInfoRepo struct {
	collection *mongo.Collection
}

func (repo *bullionSiteInfoRepo) Save(entity *interfaces.BullionSiteInfo) (result interfaces.BullionSiteInfo, err error) {
	err = repo.collection.FindOneAndUpdate(mongodb.MongoCtx, bson.D{{
		Key: "_id", Value: entity.ID,
	}}, bson.D{{Key: "$set", Value: entity}}, findOneAndUpdateOptions).Decode(&result)
	return
}

func (repo *bullionSiteInfoRepo) FindOne(id string) (result interfaces.BullionSiteInfo) {
	repo.collection.FindOne(mongodb.MongoCtx, bson.D{{
		Key: "id", Value: id,
	}}).Decode(&result)
	return
}

func (repo *bullionSiteInfoRepo) FindOneByDomain(domain string) (result interfaces.BullionSiteInfo) {
	repo.collection.FindOne(mongodb.MongoCtx, bson.D{{
		Key: "domain", Value: domain,
	}}).Decode(&result)
	return
}

func (repo *bullionSiteInfoRepo) FindByShortName(name string) (result interfaces.BullionSiteInfo) {
	repo.collection.FindOne(mongodb.MongoCtx, bson.D{{
		Key: "shortName", Value: name,
	}}).Decode(&result)
	return
}
