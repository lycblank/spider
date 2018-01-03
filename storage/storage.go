package storage

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type StorageCell struct {
	Key   string
	Value interface{}
}

var defaultStorage *kvStroage

func init() {
	db, err := leveldb.OpenFile("datas/db", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defaultStorage = &kvStroage{
		db: db,
	}
}

type kvStroage struct {
	db *leveldb.DB
}

func (ks *kvStroage) exec(prefix string, cell StorageCell) error {
	// 存储真实的值
	key := ks.getKey(prefix, cell.Key, true)
	ks.db.Put(key, ks.getValue(cell.Value), nil)
	return nil
}

func (ks *kvStroage) getKey(prefix string, key string, isPrefix bool) []byte {
	if isPrefix {
		prefix = time.Now().Format("2006-01-02") + "-" + prefix
	}
	w := md5.New()
	w.Write([]byte(key))
	content := base64.StdEncoding.EncodeToString(w.Sum(nil))
	return append([]byte{}, []byte(prefix+"-"+content)...)
}

func (ks *kvStroage) getValue(value interface{}) []byte {
	datas, _ := json.Marshal(value)
	content := base64.StdEncoding.EncodeToString(datas)
	return append([]byte{}, []byte(content)...)
}

func (ks *kvStroage) getStorageContent(prefix string, filter func(raw []byte) bool) [][]byte {
	cells := [][]byte{}
	iter := ks.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	for iter.Next() {
		v := iter.Value()
		jsonData, _ := base64.StdEncoding.DecodeString(string(v))
		if !filter(jsonData) {
			cells = append(cells, jsonData)
		}
	}
	return cells
}

func (ks *kvStroage) check(prefix string, k string) bool {
	key := ks.getKey(prefix, k, false)
	ret, err := ks.db.Has(key, nil)
	return err == nil && ret
}

func Storage(prefix string, cells []StorageCell) {
	for _, cell := range cells {
		defaultStorage.exec(prefix, cell)
	}
}

func GetStorageContent(tm time.Time, prefix string, filter func(raw []byte) bool) [][]byte {
	prefix = tm.Format("2006-01-02") + "-" + prefix + "-"
	return defaultStorage.getStorageContent(prefix, filter)
}

func Check(prefix string, key string) bool {
	return defaultStorage.check(prefix, key)
}

func Put(prefix string, key string, value string) error {
	return defaultStorage.db.Put(defaultStorage.getKey(prefix, key, false), []byte(value), nil)
}
