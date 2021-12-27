package service

import (
	"juriback2/model"
	"juriback2/reposirory"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type AkunServiceInterface interface {
	SerSetSession(input model.SessiInput) (*mongo.InsertOneResult, error)
	SerGetSession(kode string, key string) (model.SessiPlain, error)
	SerDelSessionLogout(kode string, w http.ResponseWriter, r *http.Request)
	SerAkunLogin(username string, password string) (model.AkunPlain, error)
	SerCheckLogin(w http.ResponseWriter, r *http.Request)
}

type initAkun struct {
	servicesAkun reposirory.AkunInterface
}

func AkunService(akunInterface reposirory.AkunInterface) *initAkun {
	return &initAkun{akunInterface}
}

// SERVICE

func (srvc *initAkun) SerSetSession(input model.SessiInput) (*mongo.InsertOneResult, error) {

	plain := model.SessiInput{
		SessiId:    input.SessiId,
		SessiKey:   input.SessiKey,
		SessiValue: input.SessiValue,
	}

	res, err := srvc.servicesAkun.RepoSetSession(plain)
	return res, err
}

func (srvc *initAkun) SerGetSession(kode string, key string) (model.SessiPlain, error) {
	data, err := srvc.servicesAkun.RepoGetSession(kode, key)
	return data, err
}

func (srvc *initAkun) SerAkunLogin(username string, password string) (model.AkunPlain, error) {
	data, err := srvc.servicesAkun.RepoAkunLogin(username, password)
	return data, err
}

func (srvc *initAkun) SerDelSessionLogout(kode string, w http.ResponseWriter, r *http.Request) {
	res, err := srvc.servicesAkun.RepoDelSession(kode)

	if err != nil {
		log.Println(err.Error())
	}

	if res != nil {
		log.Println(res.DeletedCount)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (srvc *initAkun) SerCheckLogin(w http.ResponseWriter, r *http.Request) {
	sekarang := time.Now().Format("2006-01-02")
	data, _ := srvc.servicesAkun.RepoGetSession(sekarang, "username")

	if data.SessiValue == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

}
