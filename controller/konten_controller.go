package controller

import (
	"encoding/json"
	"html/template"
	"io"
	"juriback2/model"
	"juriback2/service"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type initKontenController struct {
	kontenServiceInterface service.KontenServiceInterface
	akunServiceInterface   service.AkunServiceInterface
}

func KontenController(kontenService service.KontenServiceInterface, akunService service.AkunServiceInterface) *initKontenController {
	return &initKontenController{kontenService, akunService}
}

// Controller

func (handler *initKontenController) Landing(w http.ResponseWriter, r *http.Request) {
	// sekarang := time.Now().Format("2006-01-02")
	// getSes, _ := handler.akunServiceInterface.SerGetSession(sekarang, "username")
	// if getSes.SessiValue != "" {
	// 	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	// 	return
	// }

	tmpl, _ := template.ParseFiles(path.Join("pages", "Landing.htm"))
	tmpl.Execute(w, nil)
}

func (handler *initKontenController) Masuk(w http.ResponseWriter, r *http.Request) {
	sekarang := time.Now().Format("2006-01-02")

	username := r.FormValue("username")
	password := r.FormValue("pass")

	if username == "" && password == "" {

		tmpl, _ := template.ParseFiles(path.Join("pages", "Landing.htm"))

		psn := model.Pesan{
			IsError:  true,
			Pesannya: "Username dan password tidak boleh ada yang kosong",
		}

		tmpl.Execute(w, psn)
	}

	data, err := handler.akunServiceInterface.SerAkunLogin(username, password)
	if err != nil {
		http.Error(w, "ERROR SerAkunLogin "+err.Error(), http.StatusInternalServerError)
		return
	}

	if data.Username != "" {
		sessIns := model.SessiInput{
			SessiId:    sekarang,
			SessiKey:   "username",
			SessiValue: data.Username,
		}

		ses, _ := handler.akunServiceInterface.SerSetSession(sessIns)
		if err != nil {
			http.Error(w, "ERROR SerSetSession "+err.Error(), http.StatusInternalServerError)
		}

		log.Println(ses)

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles(path.Join("pages", "landing.htm"))
	if err != nil {
		http.Error(w, "Error tmlp "+err.Error(), http.StatusInternalServerError)
		return
	}

	psn := model.Pesan{
		IsError:  true,
		Pesannya: "Username dan password tidak cocok",
	}

	tmpl.Execute(w, psn)
}

func (handler *initKontenController) Keluar(w http.ResponseWriter, r *http.Request) {
	sekarang := time.Now().Format("2006-01-02")
	handler.akunServiceInterface.SerDelSessionLogout(sekarang, w, r)
}

func (handler *initKontenController) Dashboard(w http.ResponseWriter, r *http.Request) {
	// addrs, err := net.InterfaceAddrs()
	// if err != nil {
	// 	os.Stderr.WriteString("Oops: " + err.Error() + "\n")
	// 	os.Exit(1)
	// }

	// log.Println(addrs)

	// SerArtikelRandom
	// d, e := handler.kontenServiceInterface.SerArtikelRandom(1)
	// if e != nil {
	// 	os.Stderr.WriteString("SerArtikelRandom : " + e.Error() + "\n")
	// 	os.Exit(1)
	// }

	// SerArtikelPopular
	// d, e := handler.kontenServiceInterface.SerArtikelPopular(1)
	// if e != nil {
	// 	os.Stderr.WriteString("SerArtikelPopular : " + e.Error() + "\n")
	// 	os.Exit(1)
	// }

	// SerArtikelTerbaru
	d, e := handler.kontenServiceInterface.SerArtikelTerbaru(1)
	if e != nil {
		os.Stderr.WriteString("SerArtikelTerbaru : " + e.Error() + "\n")
		os.Exit(1)
	}

	log.Println(d)

	tmpl, _ := template.ParseFiles(path.Join("pages", "Dashboard.htm"), path.Join("pages", "Layout.htm"))

	data := model.Push{
		ActivePage: "dashboard",
	}

	tmpl.Execute(w, data)
}

func (handler *initKontenController) ArtikelList(w http.ResponseWriter, r *http.Request) {
	handler.akunServiceInterface.SerCheckLogin(w, r)

	tmpl, err := template.ParseFiles(path.Join("pages", "ArtikelList.htm"), path.Join("pages", "Layout.htm"))
	if err != nil {
		http.Error(w, "Error tmlp", http.StatusInternalServerError)
		return
	}

	dt, err := handler.kontenServiceInterface.SerKontenByJenis("artikel")
	if err != nil {
		http.Error(w, "ERROR Landing "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := model.PushKontenPlain{
		ActivePage:  "artikel",
		KontenPlain: dt,
	}

	tmpl.Execute(w, data)
}

func (handler *initKontenController) ArtikelWrite(w http.ResponseWriter, r *http.Request) {
	handler.akunServiceInterface.SerCheckLogin(w, r)

	tmpl, _ := template.ParseFiles(path.Join("pages", "ArtikelTulis.htm"), path.Join("pages", "Layout.htm"))

	data := model.Push{
		ActivePage: "artikel",
	}

	tmpl.Execute(w, data)
}

func (handler *initKontenController) PortofoliolList(w http.ResponseWriter, r *http.Request) {
	handler.akunServiceInterface.SerCheckLogin(w, r)

	tmpl, err := template.ParseFiles(path.Join("pages", "PortofolioList.htm"), path.Join("pages", "Layout.htm"))
	if err != nil {
		http.Error(w, "ERROR tmpl "+err.Error(), http.StatusInternalServerError)
		return
	}

	dt, err := handler.kontenServiceInterface.SerKontenByJenis("portofolio")
	if err != nil {
		http.Error(w, "ERROR SerKontenByJenis "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := model.PushKontenPlain{
		ActivePage:  "portofolio",
		KontenPlain: dt,
	}

	tmpl.Execute(w, data)
}

func (handler *initKontenController) PortofolioWrite(w http.ResponseWriter, r *http.Request) {
	handler.akunServiceInterface.SerCheckLogin(w, r)

	tmpl, _ := template.ParseFiles(path.Join("pages", "PortofolioTulis.htm"), path.Join("pages", "Layout.htm"))

	data := model.Push{
		ActivePage: "portofolio",
	}

	tmpl.Execute(w, data)
}

func (handler *initKontenController) ContentSave(w http.ResponseWriter, r *http.Request) {
	handler.akunServiceInterface.SerCheckLogin(w, r)

	judul := r.FormValue("judul")
	isi := r.FormValue("isi")
	tipe := r.FormValue("tipe")
	tahun := r.FormValue("tahun")
	bulan := r.FormValue("bulan")
	tag := r.FormValue("tag")

	if judul != "" && isi != "" && tipe != "" {

		file, hand, err := r.FormFile("thumb")
		if err != nil {
			if tipe == "artikel" {
				tmpl, _ := template.ParseFiles(path.Join("pages", "ArtikelTulisError.htm"), path.Join("pages", "Layout.htm"))
				data := model.ContentWriteError{
					ActivePage: "artikel",
					Pesan: []model.Pesan{
						{
							IsError:  true,
							Pesannya: "Error gambar",
						},
					},
					Judul: judul,
					Isi:   isi,
					Tahun: tahun,
					Bulan: bulan,
					Tag:   tag,
				}
				tmpl.Execute(w, data)
				return
			} else {
				tmpl, _ := template.ParseFiles(path.Join("pages", "PortofolioTulisError.htm"), path.Join("pages", "Layout.htm"))
				data := model.ContentWriteError{
					ActivePage: "portofolio",
					Pesan: []model.Pesan{
						{
							IsError:  true,
							Pesannya: "Error gambar",
						},
					},
					Judul: judul,
					Isi:   isi,
					Tahun: tahun,
					Bulan: bulan,
					Tag:   tag,
				}
				tmpl.Execute(w, data)
				return
			}
		}

		defer file.Close()

		ran := primitive.NewObjectID().Hex()

		var re = regexp.MustCompile("[^a-z0-9]+")
		slug := strings.Trim(re.ReplaceAllString(strings.ToLower(judul), "-"), "-") + "-" + ran

		dtSlug, _ := handler.kontenServiceInterface.SerKontenCariBySlug(slug)

		if len(dtSlug.Judul) > 0 {
			if tipe == "artikel" {
				tmpl, _ := template.ParseFiles(path.Join("pages", "ArtikelTulisError.htm"), path.Join("pages", "Layout.htm"))
				data := model.ContentWriteError{
					ActivePage: "artikel",
					Pesan: []model.Pesan{
						{
							IsError:  true,
							Pesannya: "Judul sudah ada",
						},
					},
					Judul: judul,
					Isi:   isi,
					Tahun: tahun,
					Bulan: bulan,
					Tag:   tag,
				}
				tmpl.Execute(w, data)
				return
			} else {
				tmpl, _ := template.ParseFiles(path.Join("pages", "PortofolioTulisError.htm"), path.Join("pages", "Layout.htm"))
				data := model.ContentWriteError{
					ActivePage: "portofolio",
					Pesan: []model.Pesan{
						{
							IsError:  true,
							Pesannya: "Judul sudah ada",
						},
					},
					Judul: judul,
					Isi:   isi,
					Tahun: tahun,
					Bulan: bulan,
					Tag:   tag,
				}
				tmpl.Execute(w, data)
				return
			}
		}

		// image
		ext := ""
		ext_tipe := hand.Header["Content-Type"][0]
		if ext_tipe == "image/png" {
			ext = "png"
		} else if ext_tipe == "image/jpeg" {
			ext = "jpeg"
		}

		if ext == "" {
			if tipe == "artikel" {
				tmpl, _ := template.ParseFiles(path.Join("pages", "ArtikelTulisError.htm"), path.Join("pages", "Layout.htm"))
				data := model.ContentWriteError{
					ActivePage: "artikel",
					Pesan: []model.Pesan{
						{
							IsError:  true,
							Pesannya: "File yang anda pilih bukan jenis gambar",
						},
					},
					Judul: judul,
					Isi:   isi,
					Tahun: tahun,
					Bulan: bulan,
					Tag:   tag,
				}
				tmpl.Execute(w, data)
				return
			} else {
				tmpl, _ := template.ParseFiles(path.Join("pages", "PortofolioTulisError.htm"), path.Join("pages", "Layout.htm"))
				data := model.ContentWriteError{
					ActivePage: "portofolio",
					Pesan: []model.Pesan{
						{
							IsError:  true,
							Pesannya: "File yang anda pilih bukan jenis gambar",
						},
					},
					Judul: judul,
					Isi:   isi,
					Tahun: tahun,
					Bulan: bulan,
					Tag:   tag,
				}
				tmpl.Execute(w, data)
				return
			}
		}

		rename := "./assets/gambar/" + ran + "." + ext
		fileName := "assets/gambar/" + ran + "." + ext

		tmpfile, err := os.Create(rename)
		if err != nil {
			log.Println("os.Create : " + err.Error())
			http.Error(w, "Error os.Create "+err.Error(), http.StatusInternalServerError)
			return
		}

		defer tmpfile.Close()

		_, err = io.Copy(tmpfile, file)
		if err != nil {
			http.Error(w, "Error io.Copy "+err.Error(), http.StatusInternalServerError)
			return
		}

		var konten model.KontenInput

		konten.Kode = ran
		konten.Tipe = tipe
		konten.Judul = judul
		konten.Isi = isi
		konten.Thumb = fileName
		konten.Slug = slug
		konten.Tahun = tahun
		konten.Bulan = bulan
		konten.Tag = tag

		ins, err := handler.kontenServiceInterface.SerKontenInsert(konten)
		if err != nil {
			http.Error(w, "ERROR SerKontenInsert "+err.Error(), http.StatusInternalServerError)
			return
		}

		redir := "/" + tipe + "list"

		log.Println(ins)
		http.Redirect(w, r, redir, http.StatusSeeOther)
		return
	} else {
		if tipe == "artikel" {
			tmpl, _ := template.ParseFiles(path.Join("pages", "ArtikelTulisError.htm"), path.Join("pages", "Layout.htm"))
			data := model.ContentWriteError{
				ActivePage: "artikel",
				Pesan: []model.Pesan{
					{
						IsError:  true,
						Pesannya: "Judul dan konten tidak boleh kosong",
					},
				},
				Judul: judul,
				Isi:   isi,
				Tahun: tahun,
				Bulan: bulan,
				Tag:   tag,
			}
			tmpl.Execute(w, data)
			return
		} else {
			tmpl, _ := template.ParseFiles(path.Join("pages", "PortofolioTulisError.htm"), path.Join("pages", "Layout.htm"))
			data := model.ContentWriteError{
				ActivePage: "portofolio",
				Pesan: []model.Pesan{
					{
						IsError:  true,
						Pesannya: "Judul dan konten tidak boleh kosong",
					},
				},
				Judul: judul,
				Isi:   isi,
				Tahun: tahun,
				Bulan: bulan,
				Tag:   tag,
			}
			tmpl.Execute(w, data)
			return
		}
	}
}

func (handler *initKontenController) ContentEdit(w http.ResponseWriter, r *http.Request) {
	handler.akunServiceInterface.SerCheckLogin(w, r)

	tmpl, _ := template.ParseFiles(path.Join("pages", "Edit.htm"), path.Join("pages", "Layout.htm"))

	id := r.URL.Query().Get("id")

	datas, err := handler.kontenServiceInterface.SerKontenCariById(id)
	if err != nil {
		http.Error(w, "ERROR SerKontenCariBySlug "+err.Error(), http.StatusInternalServerError)
		return
	}

	tipenya := datas.Tipe

	push := model.PushKontenPlainObj{
		ActivePage:  tipenya,
		KontenPlain: datas,
	}

	err = tmpl.Execute(w, push)
	if err != nil {
		http.Error(w, "ERROR ContentEdit "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *initKontenController) ContentUpdate(w http.ResponseWriter, r *http.Request) {
	handler.akunServiceInterface.SerCheckLogin(w, r)

	id := r.FormValue("id")
	kode := r.FormValue("kode")
	judul := r.FormValue("judul")
	isi := r.FormValue("isi")
	tipe := r.FormValue("tipe")
	tahun := r.FormValue("tahun")
	bulan := r.FormValue("bulan")
	tag := r.FormValue("tag")
	fileName := r.FormValue("gambar_existing")

	if kode != "" && judul != "" && isi != "" && tipe != "" {

		file, hand, err := r.FormFile("thumb")
		if err != nil {
			if err.Error() == "http: no such file" {

				var konten model.KontenUpdate

				konten.Judul = judul
				konten.Isi = isi
				konten.Thumb = fileName
				konten.Tahun = tahun
				konten.Bulan = bulan
				konten.Tag = tag

				dt, err := handler.kontenServiceInterface.SerKontenUpdate(kode, konten)
				if err != nil {
					tmpl, _ := template.ParseFiles(path.Join("pages", "EditError.htm"), path.Join("pages", "Layout.htm"))
					data := model.ContentUpdateError{
						ActivePage: tipe,
						Pesan: []model.Pesan{
							{
								IsError:  true,
								Pesannya: "Gagal update konten tanpa thumbnail error : " + err.Error(),
							},
						},

						Kode:  kode,
						Judul: judul,
						Isi:   isi,
						Tahun: tahun,
						Bulan: bulan,
						Tag:   tag,
						Thumb: fileName,
						Tipe:  tipe,
					}

					tmpl.Execute(w, data)
					return
				}

				redir := "/" + tipe + "list"
				log.Println(dt)
				http.Redirect(w, r, redir, http.StatusSeeOther)
				return
			}

			tmpl, _ := template.ParseFiles(path.Join("pages", "EditError.htm"), path.Join("pages", "Layout.htm"))
			data := model.ContentUpdateError{
				ActivePage: tipe,
				Pesan: []model.Pesan{
					{
						IsError:  true,
						Pesannya: "Error pada file thumbnail saat update konten",
					},
				},

				Kode:  kode,
				Judul: judul,
				Isi:   isi,
				Tahun: tahun,
				Bulan: bulan,
				Tag:   tag,
				Thumb: fileName,
				Tipe:  tipe,
			}

			tmpl.Execute(w, data)
			return
		}

		defer file.Close()

		ran := primitive.NewObjectID().Hex()
		// image
		ext_tipe := hand.Header["Content-Type"][0]
		ext := "png"
		if ext_tipe == "image/jpeg" {
			ext = "jpeg"
		}

		rename := "./assets/gambar/" + ran + "." + ext
		fileNameChange := "assets/gambar/" + ran + "." + ext

		tmpfile, err := os.Create(rename)
		if err != nil {
			http.Error(w, "ERROR os.Create "+err.Error(), http.StatusInternalServerError)
		}

		defer tmpfile.Close()

		_, err = io.Copy(tmpfile, file)
		if err != nil {
			http.Error(w, "ERROR io.Copy "+err.Error(), http.StatusInternalServerError)
		}

		var konten model.KontenUpdate

		konten.Judul = judul
		konten.Isi = isi
		konten.Thumb = fileNameChange
		konten.Tahun = tahun
		konten.Bulan = bulan
		konten.Tag = tag

		dt, err := handler.kontenServiceInterface.SerKontenUpdate(kode, konten)
		if err != nil {
			tmpl, _ := template.ParseFiles(path.Join("pages", "EditError.htm"), path.Join("pages", "Layout.htm"))
			data := model.ContentUpdateError{
				ActivePage: tipe,
				Pesan: []model.Pesan{
					{
						IsError:  true,
						Pesannya: "Gagal update konten dan thumbnail error : " + err.Error(),
					},
				},

				Kode:  kode,
				Judul: judul,
				Isi:   isi,
				Tahun: tahun,
				Bulan: bulan,
				Tag:   tag,
				Tipe:  tipe,
				Thumb: fileName,
			}

			tmpl.Execute(w, data)
			return
		}

		redir := "/" + tipe + "list"
		log.Println(id)
		log.Println(dt)
		http.Redirect(w, r, redir, http.StatusSeeOther)
		return
	}

	tmpl, _ := template.ParseFiles(path.Join("pages", "EditError.htm"), path.Join("pages", "Layout.htm"))
	data := model.ContentUpdateError{
		ActivePage: tipe,
		Pesan: []model.Pesan{
			{
				IsError:  true,
				Pesannya: "Judul dan isi konten harus diisi",
			},
		},

		Kode:  kode,
		Judul: judul,
		Isi:   isi,
		Tahun: tahun,
		Bulan: bulan,
		Tag:   tag,
		Tipe:  tipe,
		Thumb: fileName,
	}

	tmpl.Execute(w, data)
}

// func (handler *initKontenController) ContentUpdate(w http.ResponseWriter, r *http.Request) {
func (handler *initKontenController) Gambar(w http.ResponseWriter, r *http.Request) {

	file, hand, err := r.FormFile("file")
	if err != nil {
		log.Println("Error Retrieving the File")
		log.Println(err)
		return
	}

	defer file.Close()

	ext := ""
	ext_tipe := hand.Header["Content-Type"][0]
	if ext_tipe == "image/png" {
		ext = "png"
	} else if ext_tipe == "image/jpeg" {
		ext = "jpeg"
	}

	type TinyMCE struct {
		FileName string `json:"location"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if ext != "" {

		ran := primitive.NewObjectID().Hex()
		rename := "./assets/gambar/" + ran + "." + ext
		fileName := "https://api.juripebrianto.my.id/assets/gambar/" + ran + "." + ext

		tmpfile, err := os.Create(rename)
		if err != nil {
			log.Println("ERROR os.Create : " + err.Error())

			p := TinyMCE{
				FileName: "https://api.juripebrianto.my.id/assets/no-image-1.png",
			}

			json.NewEncoder(w).Encode(p)
			return
		}

		defer tmpfile.Close()

		_, err = io.Copy(tmpfile, file)
		if err != nil {
			log.Println("ERROR osio.Copy : " + err.Error())

			p := TinyMCE{
				FileName: "https://api.juripebrianto.my.id/assets/no-image-2.png",
			}

			json.NewEncoder(w).Encode(p)
			return
		}

		p := TinyMCE{
			FileName: fileName,
		}

		json.NewEncoder(w).Encode(p)
		return
	}

	p := TinyMCE{
		FileName: "https://api.juripebrianto.my.id/assets/no-image-3.png",
	}

	json.NewEncoder(w).Encode(p)
}

// API ------------------------------------------------------------------------------------
func (handler *initKontenController) RandomArtikel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(http.StatusOK)

	data, err := handler.kontenServiceInterface.SerKontenRandom("artikel", 1)
	if err != nil {
		push := model.JsonApi{
			Stat:       400,
			KontenJson: nil,
		}

		json.NewEncoder(w).Encode(push)
		return
	}

	// Kode    string `json:"kode"`
	// Tipe    string `json:"tipe"`
	// Judul   string `json:"judul"`
	// Short   string `json:"short"`
	// Isi     string `json:"isi"`
	// Thumb   string `json:"thumb"`
	// Tanggal string `json:"tanggal"`
	// Slug    string `json:"slug"`
	// Tag     string `json:"tag"`
	// View    int    `json:"view"`

	push := model.JsonApi{
		Stat:       400,
		KontenJson: data,
	}

	json.NewEncoder(w).Encode(push)
}

func (handler *initKontenController) RandomPortofolio(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(http.StatusOK)

	data, err := handler.kontenServiceInterface.SerKontenRandom("portofolio", 1)
	if err != nil {
		push := model.JsonApi{
			Stat:       400,
			KontenJson: nil,
		}

		json.NewEncoder(w).Encode(push)
		return
	}

	push := model.JsonApi{
		Stat:       400,
		KontenJson: data,
	}

	json.NewEncoder(w).Encode(push)
}

func (handler *initKontenController) RandomKonten(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(http.StatusOK)

	data, err := handler.kontenServiceInterface.SerKontenRandom("", 3)
	if err != nil {
		push := model.JsonApi{
			Stat:       400,
			KontenJson: nil,
		}

		json.NewEncoder(w).Encode(push)
		return
	}

	push := model.JsonApi{
		Stat:       400,
		KontenJson: data,
	}

	json.NewEncoder(w).Encode(push)
}

func (handler *initKontenController) ArtikelTerbaru(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(http.StatusOK)

	data, err := handler.kontenServiceInterface.SerArtikelTerbaru(3)
	if err != nil {
		push := model.JsonApi{
			Stat:       400,
			KontenJson: nil,
		}

		json.NewEncoder(w).Encode(push)
		return
	}

	push := model.JsonApi{
		Stat:       400,
		KontenJson: data,
	}

	json.NewEncoder(w).Encode(push)
}

func (handler *initKontenController) ArtikelPopular(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(http.StatusOK)

	data, err := handler.kontenServiceInterface.SerArtikelPopular(6)
	if err != nil {
		push := model.JsonApi{
			Stat:       400,
			KontenJson: nil,
		}

		json.NewEncoder(w).Encode(push)
		return
	}

	push := model.JsonApi{
		Stat:       400,
		KontenJson: data,
	}

	json.NewEncoder(w).Encode(push)
}

func (handler *initKontenController) SemuaArtikelRandom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(http.StatusOK)

	data, err := handler.kontenServiceInterface.SerKontenRandom("artikel", 10000)
	if err != nil {
		push := model.JsonApi{
			Stat:       400,
			KontenJson: nil,
		}

		json.NewEncoder(w).Encode(push)
		return
	}

	push := model.JsonApi{
		Stat:       400,
		KontenJson: data,
	}

	json.NewEncoder(w).Encode(push)
}

func (handler *initKontenController) SemuaPortofolioDesc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(http.StatusOK)

	data, err := handler.kontenServiceInterface.SerSemuaPortofolio()
	if err != nil {
		push := model.JsonApi{
			Stat:       400,
			KontenJson: nil,
		}

		json.NewEncoder(w).Encode(push)
		return
	}

	push := model.JsonApi{
		Stat:       400,
		KontenJson: data,
	}

	json.NewEncoder(w).Encode(push)
}

func (handler *initKontenController) DetailArtikel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(http.StatusOK)

	slug := r.URL.Query().Get("id")

	data, err := handler.kontenServiceInterface.SerKontenCariBySlug(slug)
	if err != nil {
		push := model.JsonApi{
			Stat:       400,
			KontenJson: nil,
		}

		json.NewEncoder(w).Encode(push)
		return
	}

	srt := ""
	if len(data.Isi) > 45 {
		srt = data.Isi[3:45]
	}

	if len(data.Isi) < 45 {
		max := len(data.Isi) - 4
		srt = data.Isi[3:max]
	}

	dataJsonnya := []model.KontenJson{
		{
			Kode:    data.Kode,
			Tipe:    data.Tipe,
			Judul:   data.Judul,
			Short:   srt,
			Isi:     data.Isi,
			Thumb:   data.Thumb,
			Tanggal: data.Tanggal,
			Slug:    data.Slug,
			Tag:     data.Tag,
			View:    0,
		},
	}

	push := model.JsonApi{
		Stat:       200,
		KontenJson: dataJsonnya,
	}

	json.NewEncoder(w).Encode(push)
}

func (handler *initKontenController) RandomArtikel6(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(http.StatusOK)

	data, err := handler.kontenServiceInterface.SerKontenRandom("artikel", 6)
	if err != nil {
		push := model.JsonApi{
			Stat:       400,
			KontenJson: nil,
		}

		json.NewEncoder(w).Encode(push)
		return
	}

	push := model.JsonApi{
		Stat:       200,
		KontenJson: data,
	}

	json.NewEncoder(w).Encode(push)
}

func (handler *initKontenController) CariKonten(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(http.StatusOK)

	cari := r.URL.Query().Get("find")

	data, err := handler.kontenServiceInterface.SerKontenCariLike(cari)
	if err != nil {
		push := model.JsonApi{
			Stat:       400,
			KontenJson: nil,
		}

		json.NewEncoder(w).Encode(push)
		return
	}

	push := model.JsonApi{
		Stat:       200,
		KontenJson: data,
	}

	json.NewEncoder(w).Encode(push)
}

func (handler *initKontenController) Views(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(http.StatusOK)

	type PushData struct {
		Stat int
	}

	slug := r.URL.Query().Get("id")

	dt, err := handler.kontenServiceInterface.SerKontenCariBySlug(slug)
	if err != nil {
		push := PushData{
			Stat: 400,
		}

		json.NewEncoder(w).Encode(push)
		return
	}

	var vw model.KontenUpdateView
	vw.View = dt.View + 1

	data, err := handler.kontenServiceInterface.SerKontenTambahView(slug, vw)
	if err != nil {
		push := PushData{
			Stat: 400,
		}

		json.NewEncoder(w).Encode(push)
		return
	}

	push := PushData{
		Stat: 200,
	}

	log.Println(data)
	json.NewEncoder(w).Encode(push)
}

// ----------------------------------------------------------------------------------------
