package reposirory

import (
	"context"
	"juriback2/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type KontenInterface interface {
	RepoKontenTambahData(b model.KontenInput) (*mongo.InsertOneResult, error)
	RepoKontenUpdateData(selector string, b model.KontenUpdate) (*mongo.UpdateResult, error)

	RepoKontenSemuaData() ([]model.KontenPlain, error)
	RepoKontenCariById(ID string) (model.KontenPlain, error)
	RepoKontenCariByKode(Kode string) (model.KontenPlain, error)
	RepoKontenCariBySlug(Slug string) (model.KontenPlain, error)
	RepoKontenCariByJenis(jenisnya string) ([]model.KontenPlain, error)

	RepoKontenTambahView(Slug string, b model.KontenUpdateView) (*mongo.UpdateResult, error)

	RepoKontenRandom(tipe string, limit int64) ([]model.KontenJson, error)
	RepoArtikelPopular(limit int64) ([]model.KontenJson, error)
	RepoArtikelTerbaru(limit int64) ([]model.KontenJson, error)
	RepoSemuaPortofolio() ([]model.KontenJson, error)
	RepoKontenCariLike(like string) ([]model.KontenJson, error)

	// UbahData(b model.KontenPlain) (model.KontenPlain, error)
	// HapusData(b model.KontenPlain) (model.KontenPlain, error)
}

type kontenDB struct {
	database *mongo.Database
}

func KontenRepository(db *mongo.Database) *kontenDB {
	return &kontenDB{db}
}

// REPO

func (db *kontenDB) RepoKontenTambahData(datas model.KontenInput) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := db.database.Collection("konten").InsertOne(ctx, &datas)
	return res, err
}

func (db *kontenDB) RepoKontenUpdateData(selector string, datas model.KontenUpdate) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{
		primitive.E{Key: "kode", Value: selector},
	}

	update := bson.M{
		"$set": datas,
	}

	res, err := db.database.Collection("konten").UpdateOne(ctx, filter, update)
	return res, err
}

func (db *kontenDB) RepoKontenSemuaData() ([]model.KontenPlain, error) {
	var mdl []model.KontenPlain
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cur, err := db.database.Collection("konten").Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	err = cur.All(ctx, &mdl)
	if err != nil {
		return nil, err
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return mdl, err
}

func (db *kontenDB) RepoKontenCariById(ID string) (model.KontenPlain, error) {
	var mdl model.KontenPlain
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(ID)

	filter := bson.D{
		primitive.E{Key: "_id", Value: objID},
	}

	err := db.database.Collection("konten").FindOne(ctx, filter).Decode(&mdl)

	return mdl, err
}

func (db *kontenDB) RepoKontenCariByKode(kode string) (model.KontenPlain, error) {
	var mdl model.KontenPlain
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{
		primitive.E{Key: "kode", Value: kode},
	}

	err := db.database.Collection("konten").FindOne(ctx, filter).Decode(&mdl)

	return mdl, err
}

func (db *kontenDB) RepoKontenCariBySlug(slug string) (model.KontenPlain, error) {
	var mdl model.KontenPlain
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{
		primitive.E{Key: "slug", Value: slug},
	}

	err := db.database.Collection("konten").FindOne(ctx, filter).Decode(&mdl)

	return mdl, err
}

func (db *kontenDB) RepoKontenCariByJenis(jenis string) ([]model.KontenPlain, error) {
	var mdl []model.KontenPlain
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cur, err := db.database.Collection("konten").Find(ctx, bson.M{"tipe": jenis})

	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	err = cur.All(ctx, &mdl)
	if err != nil {
		return nil, err
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return mdl, err
}

func (db *kontenDB) RepoKontenTambahView(selector string, datas model.KontenUpdateView) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{
		primitive.E{Key: "slug", Value: selector},
	}

	update := bson.M{
		"$set": datas,
	}

	res, err := db.database.Collection("konten").UpdateOne(ctx, filter, update)
	return res, err
}

func (db *kontenDB) RepoKontenRandom(tipe string, limit int64) ([]model.KontenJson, error) {
	var mdl []model.KontenJson
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if tipe != "" {
		pipeline := []bson.M{
			{"$match": bson.M{"tipe": tipe}},
			{"$sample": bson.M{"size": limit}},
		}

		cur, err := db.database.Collection("konten").Aggregate(ctx, pipeline)
		if err != nil {
			return nil, err
		}

		defer cur.Close(ctx)

		err = cur.All(ctx, &mdl)
		if err != nil {
			return nil, err
		}

		if err := cur.Err(); err != nil {
			return nil, err
		}

		return mdl, err
	} else {
		pipeline := []bson.M{
			{"$sample": bson.M{"size": limit}},
		}

		cur, err := db.database.Collection("konten").Aggregate(ctx, pipeline)
		if err != nil {
			return nil, err
		}

		defer cur.Close(ctx)

		err = cur.All(ctx, &mdl)
		if err != nil {
			return nil, err
		}

		if err := cur.Err(); err != nil {
			return nil, err
		}

		return mdl, err
	}

}

func (db *kontenDB) RepoArtikelPopular(limit int64) ([]model.KontenJson, error) {
	var mdl []model.KontenJson
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := db.database.Collection("konten").Find(ctx, bson.M{"tipe": "artikel"}, options.Find().SetSort(map[string]int{"view": -1}).SetLimit(limit))

	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	err = cur.All(ctx, &mdl)
	if err != nil {
		return nil, err
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return mdl, err
}

func (db *kontenDB) RepoArtikelTerbaru(limit int64) ([]model.KontenJson, error) {
	var mdl []model.KontenJson
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := db.database.Collection("konten").Find(ctx, bson.M{"tipe": "artikel"}, options.Find().SetSort(map[string]int{"tanggal": 1}).SetLimit(limit))

	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	err = cur.All(ctx, &mdl)
	if err != nil {
		return nil, err
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return mdl, err
}

func (db *kontenDB) RepoSemuaPortofolio() ([]model.KontenJson, error) {
	var mdl []model.KontenJson
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := db.database.Collection("konten").Find(ctx, bson.M{"tipe": "portofolio"}, options.Find().SetSort(bson.D{
		{Key: "tahun", Value: -1},
		{Key: "bulan", Value: -1},
	}))

	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	err = cur.All(ctx, &mdl)
	if err != nil {
		return nil, err
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return mdl, err
}

func (db *kontenDB) RepoKontenCariLike(like string) ([]model.KontenJson, error) {
	var mdl []model.KontenJson
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := bson.M{
		"judul": bson.M{
			"$regex": primitive.Regex{
				Pattern: like,
				Options: "i",
			},
		},
	}

	cur, err := db.database.Collection("konten").Find(ctx, query)

	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	err = cur.All(ctx, &mdl)
	if err != nil {
		return nil, err
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return mdl, err
}
