package service

import (
	"juriback2/model"
	"juriback2/reposirory"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type KontenServiceInterface interface {
	SerKontenInsert(input model.KontenInput) (*mongo.InsertOneResult, error)
	SerKontenUpdate(kode string, input model.KontenUpdate) (*mongo.UpdateResult, error)

	SerKontenSemua() ([]model.KontenPlain, error)
	SerKontenCariById(IDnya string) (model.KontenPlain, error)
	SerKontenCariByKode(Kodenya string) (model.KontenPlain, error)
	SerKontenCariBySlug(Slugnya string) (model.KontenPlain, error)
	SerKontenByJenis(Jenisnya string) ([]model.KontenPlain, error)

	SerKontenTambahView(kode string, input model.KontenUpdateView) (*mongo.UpdateResult, error)

	SerKontenRandom(tipe string, limit int64) ([]model.KontenJson, error)
	SerArtikelPopular(limit int64) ([]model.KontenJson, error)
	SerArtikelTerbaru(limit int64) ([]model.KontenJson, error)
	SerSemuaPortofolio() ([]model.KontenJson, error)

	// InsertDataService(input models.PrizesInput) (models.Prize, error)
	// UpdateDataService(ID int, input models.PrizesInput) (models.Prize, error)
	// HapusDataService(ID int) (models.Prize, error)
}

type initKonten struct {
	services reposirory.KontenInterface
}

func KontenService(kontenInterface reposirory.KontenInterface) *initKonten {
	return &initKonten{kontenInterface}
}

// SERVICE

func (srvc *initKonten) SerKontenInsert(input model.KontenInput) (*mongo.InsertOneResult, error) {

	if input.Tipe == "artikel" {
		plain := model.KontenInput{
			Kode:    input.Kode,
			Tipe:    input.Tipe,
			Judul:   input.Judul,
			Isi:     input.Isi,
			Thumb:   input.Thumb,
			Tanggal: time.Now().Format("2006-01-02 15:04:05"),
			Slug:    input.Slug,
			Tahun:   "0",
			Bulan:   "0",
			Tag:     input.Tag,
			View:    0,
		}

		res, err := srvc.services.RepoKontenTambahData(plain)
		return res, err
	} else {
		plain := model.KontenInput{
			Kode:    input.Kode,
			Tipe:    input.Tipe,
			Judul:   input.Judul,
			Isi:     input.Isi,
			Thumb:   input.Thumb,
			Tanggal: time.Now().Format("2006-01-02 15:04:05"),
			Slug:    input.Slug,
			Tahun:   input.Tahun,
			Bulan:   input.Bulan,
			Tag:     input.Tag,
			View:    0,
		}

		res, err := srvc.services.RepoKontenTambahData(plain)
		return res, err
	}
}

func (srvc *initKonten) SerKontenUpdate(kode string, input model.KontenUpdate) (*mongo.UpdateResult, error) {

	plain := model.KontenUpdate{
		Judul: input.Judul,
		Isi:   input.Isi,
		Thumb: input.Thumb,
		Tahun: input.Tahun,
		Bulan: input.Bulan,
		Tag:   input.Tag,
	}

	res, err := srvc.services.RepoKontenUpdateData(kode, plain)
	return res, err
}

func (srvc *initKonten) SerKontenUpdateView(kode string, input model.KontenUpdateView) (*mongo.UpdateResult, error) {

	plain := model.KontenUpdateView{
		View: input.View,
	}

	res, err := srvc.services.RepoKontenTambahView(kode, plain)
	return res, err
}

func (srvc *initKonten) SerKontenSemua() ([]model.KontenPlain, error) {
	data, err := srvc.services.RepoKontenSemuaData()
	return data, err
}

func (srvc *initKonten) SerKontenCariById(IDnya string) (model.KontenPlain, error) {
	data, err := srvc.services.RepoKontenCariById(IDnya)
	return data, err
}

func (srvc *initKonten) SerKontenCariByKode(Kodenya string) (model.KontenPlain, error) {
	data, err := srvc.services.RepoKontenCariByKode(Kodenya)
	return data, err
}

func (srvc *initKonten) SerKontenCariBySlug(Slugnya string) (model.KontenPlain, error) {
	data, err := srvc.services.RepoKontenCariBySlug(Slugnya)
	return data, err
}

func (srvc *initKonten) SerKontenByJenis(Jenisnya string) ([]model.KontenPlain, error) {
	data, err := srvc.services.RepoKontenCariByJenis(Jenisnya)
	return data, err
}

func (srvc *initKonten) SerKontenTambahView(Slugnya string, vws model.KontenUpdateView) (*mongo.UpdateResult, error) {

	plain := model.KontenUpdateView{
		View: vws.View,
	}

	res, err := srvc.services.RepoKontenTambahView(Slugnya, plain)
	return res, err
}

func (srvc *initKonten) SerKontenRandom(tipe string, limit int64) ([]model.KontenJson, error) {
	data, err := srvc.services.RepoKontenRandom(tipe, limit)
	return data, err
}

func (srvc *initKonten) SerArtikelPopular(limit int64) ([]model.KontenJson, error) {
	data, err := srvc.services.RepoArtikelPopular(limit)
	return data, err
}

func (srvc *initKonten) SerArtikelTerbaru(limit int64) ([]model.KontenJson, error) {
	data, err := srvc.services.RepoArtikelTerbaru(limit)
	return data, err
}

func (srvc *initKonten) SerSemuaPortofolio() ([]model.KontenJson, error) {
	data, err := srvc.services.RepoSemuaPortofolio()
	return data, err
}
