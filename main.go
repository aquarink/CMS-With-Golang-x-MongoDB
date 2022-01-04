package main

import (
	"context"
	"juriback2/controller"
	"juriback2/reposirory"
	"juriback2/service"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	r := http.NewServeMux()

	// assets
	assets := http.FileServer((http.Dir("/var/user/back/assets/")))
	r.Handle("/static/", http.StripPrefix("/static/", assets))

	// REPO SERVICE CONTROLLER
	kontenRepo := reposirory.KontenRepository(conn())
	kontenService := service.KontenService(kontenRepo)

	akunRepo := reposirory.AkunRepository(conn())
	akunService := service.AkunService(akunRepo)

	kontenController := controller.KontenController(kontenService, akunService)

	r.HandleFunc("/in", kontenController.Landing)
	r.HandleFunc("/masuk", kontenController.Masuk)
	r.HandleFunc("/keluar", kontenController.Keluar)
	r.HandleFunc("/dashboard", kontenController.Dashboard)

	r.HandleFunc("/artikellist", kontenController.ArtikelList)
	r.HandleFunc("/artikelform", kontenController.ArtikelWrite)

	r.HandleFunc("/portofoliolist", kontenController.PortofoliolList)
	r.HandleFunc("/portofolioform", kontenController.PortofolioWrite)

	r.HandleFunc("/edit", kontenController.ContentEdit)
	r.HandleFunc("/save", kontenController.ContentSave)
	r.HandleFunc("/update", kontenController.ContentUpdate)

	r.HandleFunc("/gambar", kontenController.Gambar)

	// API ----------------------
	r.HandleFunc("/random-artikel", kontenController.RandomArtikel)
	r.HandleFunc("/random-portofolio", kontenController.RandomPortofolio)
	r.HandleFunc("/random-konten", kontenController.RandomKonten)
	r.HandleFunc("/artikel-terbaru", kontenController.ArtikelTerbaru)
	r.HandleFunc("/artikel-popular", kontenController.ArtikelPopular)
	r.HandleFunc("/artikel-semua", kontenController.SemuaArtikelRandom)
	r.HandleFunc("/portofolio-semua", kontenController.SemuaPortofolioDesc)

	r.HandleFunc("/detail-artikel", kontenController.DetailArtikel)
	r.HandleFunc("/random-artikel6", kontenController.RandomArtikel6)
	r.HandleFunc("/konten-cari", kontenController.CariKonten)
	r.HandleFunc("/view", kontenController.Views)

	// SERVER
	err := http.ListenAndServe(":8899", r)
	if err != nil {
		log.Println("Error 8899 " + err.Error())
	}
}

func conn() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	monitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			// log.Print(evt.Command)
		},
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017").SetMonitor(monitor))
	if err != nil {
		return nil
	}

	return client.Database("pebridb")
}
