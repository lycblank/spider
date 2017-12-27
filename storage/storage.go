package storage

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

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

func (ks *kvStroage) exec(cell StorageCell) error {
	// 存储真实的值
	ks.db.Put(ks.getKey(cell.Title, true), ks.getValue(cell), nil)
	return nil
}

func (ks *kvStroage) getKey(title string, isPrefix bool) []byte {
	prefix := ""
	if isPrefix {
		prefix = time.Now().Format("2006-01-02")
	}
	content := base64.StdEncoding.EncodeToString([]byte(title))
	return append([]byte{}, []byte(prefix+content)...)
}

func (ks *kvStroage) getValue(cell StorageCell) []byte {
	datas, _ := json.Marshal(cell)
	content := base64.StdEncoding.EncodeToString(datas)
	return append([]byte{}, []byte(content)...)
}

func (ks *kvStroage) getStorageContent(prefix string, filter func(cell *StorageCell) bool) []StorageCell {
	cells := []StorageCell{}
	iter := ks.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	for iter.Next() {
		v := iter.Value()
		jsonData, _ := base64.StdEncoding.DecodeString(string(v))
		cell := &StorageCell{}
		json.Unmarshal(jsonData, cell)
		if !filter(cell) {
			cells = append(cells, *cell)
		}
	}
	return cells
}

func (ks *kvStroage) check(title string) bool {
	key := ks.getKey(title, false)
	ret, err := ks.db.Has(key, nil)
	return err == nil && ret
}

type StorageCell struct {
	Href   string `json:"href"`
	Title  string `json:"title"`
	Source string `json:"source"`
}

func Storage(cells []StorageCell) {
	for _, cell := range cells {
		defaultStorage.exec(cell)
	}
}

func GetStorageContent(tm time.Time, filter func(cell *StorageCell) bool) []StorageCell {
	prefix := tm.Format("2006-01-02")
	return defaultStorage.getStorageContent(prefix, filter)
}

func Check(title string) bool {
	return defaultStorage.check(title)
}

func Put(key string, value string) error {
	return defaultStorage.db.Put(defaultStorage.getKey(key, false), []byte(value), nil)
}
