package reposirory

import (
	"context"
	"juriback2/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AkunInterface interface {
	RepoSetSession(b model.SessiInput) (*mongo.InsertOneResult, error)
	RepoGetSession(Kode string, Key string) (model.SessiPlain, error)
	RepoDelSession(sessiid string) (*mongo.DeleteResult, error)

	RepoAkunLogin(username string, password string) (model.AkunPlain, error)
}

type akunDB struct {
	databaseAkun *mongo.Database
}

func AkunRepository(db *mongo.Database) *akunDB {
	return &akunDB{db}
}

// REPO

func (db *akunDB) RepoSetSession(datas model.SessiInput) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := db.databaseAkun.Collection("sessi").InsertOne(ctx, &datas)
	return res, err
}

func (db *akunDB) RepoGetSession(kode string, key string) (model.SessiPlain, error) {
	var mdl model.SessiPlain
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{
		primitive.E{Key: "sessiid", Value: kode},
		primitive.E{Key: "sessikey", Value: key},
	}

	err := db.databaseAkun.Collection("sessi").FindOne(ctx, filter).Decode(&mdl)

	return mdl, err
}

func (db *akunDB) RepoDelSession(sessiid string) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	del, err := db.databaseAkun.Collection("sessi").DeleteOne(ctx, bson.M{"sessiid": sessiid})

	return del, err
}

func (db *akunDB) RepoAkunLogin(username string, password string) (model.AkunPlain, error) {
	var mdl model.AkunPlain
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{
		primitive.E{Key: "username", Value: username},
		primitive.E{Key: "password", Value: password},
	}

	err := db.databaseAkun.Collection("akun").FindOne(ctx, filter).Decode(&mdl)

	return mdl, err
}
