package main

import (
	"bytes"
	"code.google.com/p/leveldb-go/leveldb/db"
	"code.google.com/p/leveldb-go/leveldb/table"
	"encoding/gob"
	//"fmt"
)

var (
	buf    bytes.Buffer
	dbname = "my.db"
	dbfs   = db.DefaultFileSystem
)

type Zuoye struct {
	Ip       string
	UserName string
	FileName string
	Score    string
}

func (zy Zuoye) save() error {
	enc := gob.NewEncoder(&buf)
	enc.Encode(zy)
	f0, err := dbfs.Open(dbname)
	if err != nil {
		f0, _ = dbfs.Create(dbname)
	}
	w := table.NewWriter(f0, nil)
	w.Set([]byte(zy.Ip), buf.Bytes(), nil)
	w.Close()
	return nil

}

func zyFromGob(gobdata []byte) Zuoye {
	gobreader := bytes.NewReader(gobdata)
	dec := gob.NewDecoder(gobreader)
	var zy Zuoye
	dec.Decode(&zy)
	return zy
}

func getZuoye(ip string) Zuoye {
	f1, _ := dbfs.Open(dbname)
	r := table.NewReader(f1, nil)
	gobdata, _ := r.Get([]byte(ip), nil)
	zy := zyFromGob(gobdata)
	return zy
}

func rmZuoye(key string) error {
	f1, _ := dbfs.Open(dbname)
	r := table.NewWriter(f1, nil)
	err := r.Delete([]byte(key), nil)
	return err
}

func scoreZuoye(ip, score string) error {
	zy := getZuoye(ip)
	zy.Score = score
	enc := gob.NewEncoder(&buf)
	enc.Encode(zy)
	f0, _ := dbfs.Open(dbname)
	w := table.NewWriter(f0, nil)
	w.Set([]byte(zy.Ip), buf.Bytes(), nil)
	w.Close()
	return nil
}

func listZuoye() []Zuoye {
	f1, _ := dbfs.Open(dbname)
	r := table.NewReader(f1, nil)
	records := r.Find(nil, nil)
	var zys []Zuoye
	for records.Next() {
		zys = append(zys, zyFromGob(records.Value()))
	}
	return zys
}
