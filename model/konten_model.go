package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type KontenFull struct {
	// type KontenPlain struct {
	Id      primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Kode    string             `json:"kode"`
	Tipe    string             `json:"tipe"`
	Judul   string             `json:"judul"`
	Short   string             `json:"short"`
	Isi     string             `json:"isi"`
	Thumb   string             `json:"thumb"`
	Tanggal string             `json:"tanggal"`
	Slug    string             `json:"slug"`
	Tag     string             `json:"tag"`
	View    int                `json:"view"`
}

type KontenPlain struct {
	// type KontenFull struct {
	Id      primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Kode    string
	Tipe    string
	Judul   string
	Short   string
	Isi     string
	Thumb   string
	Tanggal string
	Slug    string
	Tahun   string
	Bulan   string
	Tag     string
	View    int
}

type KontenInput struct {
	Kode    string `json:"kode"`
	Tipe    string `json:"tipe"`
	Judul   string `json:"judul"`
	Short   string `json:"short"`
	Isi     string `json:"isi"`
	Thumb   string `json:"thumb"`
	Tanggal string `json:"tanggal"`
	Slug    string `json:"slug"`
	Tahun   string `json:"tahun"`
	Bulan   string `json:"bulan"`
	Tag     string `json:"tag"`
	View    int    `json:"view"`
}

type KontenUpdate struct {
	Judul string `json:"judul"`
	Short string `json:"short"`
	Isi   string `json:"isi"`
	Thumb string `json:"thumb"`
	Tahun string `json:"tahun"`
	Bulan string `json:"bulan"`
	Tag   string `json:"tag"`
}

type KontenUpdateView struct {
	View int `json:"view"`
}

type AkunFull struct {
	Id       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Username string             `json:"username"`
	Password string             `json:"password"`
}

type AkunPlain struct {
	Id       string
	Username string
	Password string
}

type Pesan struct {
	IsError  bool
	Pesannya string
}

type Push struct {
	ActivePage string
}

type ContentWriteError struct {
	ActivePage string
	Pesan      []Pesan
	Judul      string
	Isi        string
	Tahun      string
	Bulan      string
	Tag        string
}

type ContentUpdateError struct {
	ActivePage string
	Pesan      []Pesan
	Kode       string
	Judul      string
	Isi        string
	Tahun      string
	Bulan      string
	Tag        string
	Thumb      string
	Tipe       string
}

type PushKontenPlain struct {
	ActivePage  string
	KontenPlain []KontenPlain
}

type PushKontenPlainObj struct {
	ActivePage  string
	KontenPlain KontenPlain
}

type PushKontenFull struct {
	ActivePage string
	KontenFull []KontenFull
}

type SessiInput struct {
	SessiId    string `json:"sessiid"`
	SessiKey   string `json:"sessikey"`
	SessiValue string `json:"sessivalue"`
}

type SessiPlain struct {
	SessiId    string
	SessiKey   string
	SessiValue string
}

// JSON

type KontenJson struct {
	Kode    string `json:"kode"`
	Tipe    string `json:"tipe"`
	Judul   string `json:"judul"`
	Short   string `json:"short"`
	Isi     string `json:"isi"`
	Thumb   string `json:"thumb"`
	Tanggal string `json:"tanggal"`
	Slug    string `json:"slug"`
	Tag     string `json:"tag"`
	View    int    `json:"view"`
}

type JsonApi struct {
	Stat       int
	KontenJson []KontenJson
}
