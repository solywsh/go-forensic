package db

import (
	"fmt"
	"github.com/solywsh/go-forensic/db/mmkv"
	"path/filepath"
	"strings"
)

type MMKV struct {
	dirPath  string
	filePath string
	keywords []string
	cryptKey string
}

type MMKVResult struct {
	key   string
	value []byte
	id    string
	path  string
}

func NewMMKV() *MMKV {
	return &MMKV{}
}

func (m *MMKV) SetDir(dir string) *MMKV {
	m.dirPath = dir
	return m
}

func (m *MMKV) SetFilePath(filePath string) *MMKV {
	m.filePath = filePath
	return m
}

func (m *MMKV) SetKeywords(keywords []string) *MMKV {
	m.keywords = keywords
	return m
}

func (m *MMKV) SetCryptKey(cryptKey string) *MMKV {
	m.cryptKey = cryptKey
	return m
}

func (m *MMKV) SearchByKeywords(keywords ...string) error {
	if len(keywords) == 0 {
		return fmt.Errorf("please input keywords")
	}
	if m.dirPath == "" && m.filePath == "" {
		return fmt.Errorf("please input dir path or filename")
	}
	var resNum int
	var err error
	defer log.Info("search result num", "num", resNum)
	if m.filePath != "" {
		resNum, err = m.searchByKeywordWithFilepath(keywords...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MMKV) searchByKeywordWithFilepath(keywords ...string) (int, error) {
	var dirPath string
	if m.dirPath != "" {
		dirPath = m.dirPath
	} else {
		dirPath = filepath.Dir(m.filePath)
	}
	_m := mmkv.FactoryNewMMKV(dirPath)
	err := _m.WithId(filepath.Base(m.filePath), m.cryptKey)
	if err != nil {
		return 0, err
	}
	keys := _m.Keys()
	if len(keys) == 0 {
		return 0, nil
	}
	var resultNum int
	for _, key := range keys {
		var hit bool
		_values := _m.Get(key)
		for _, keyword := range keywords {
			if strings.Contains(key, keyword) {
				resultNum++
				hit = true
			} else if strings.Contains(string(_values), keyword) {
				resultNum++
				hit = true
			}
			if hit {
				log.Info("hit keywords", "key", key, "value", _values, "keywords", keyword, "fileName", filepath.Base(m.filePath))
				break
			}
		}
	}
	return resultNum, nil
}

func (m *MMKV) searchByKeywordsWithDir(keywords ...string) error {
	return nil
}
