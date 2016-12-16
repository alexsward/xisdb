package xisdb

import (
	"bufio"
	"fmt"
	"strings"
)

func (db *DB) load() error {
	fmt.Println("loading")

	sc := bufio.NewScanner(db.file)
	for sc.Scan() {
		_, op, kv := extract(sc.Text())
		if op == "-" {
			key, err := getKey(kv, -1, db.fileErrors)
			if err != nil {
				return err
			}
			db.remove(&Item{key, "", nil})
			continue
		}

		vi := strings.LastIndex(kv, "v~")
		if vi == -1 && op == "+" {
			return ErrIncorrectDatabaseFileFormat
		}

		key, err := getKey(kv, vi, db.fileErrors)
		if err != nil {
			return err
		}
		value := kv[vi+2 : len(kv)]
		db.insert(&Item{key, value, nil})
	}

	return nil
}

func getKey(kv string, vi int, fileErrors bool) (string, error) {
	i := strings.Index(kv, "k~")
	if i == -1 && fileErrors {
		return "", ErrIncorrectDatabaseFileFormat
	}

	if vi == -1 {
		return kv[2:len(kv)], nil
	}

	return kv[2 : vi-1], nil
}

func extract(s string) (time, op, kv string) {
	strs := strings.SplitAfterN(strings.Trim(s, " \n"), " ", 3)

	time = strs[0]
	op = strs[1][0:1]
	kv = strs[2]
	return
}

func (tx *Tx) persist() error {
	var buf []byte
	for key, item := range tx.commits {
		var s string
		if item != nil {
			s = fmt.Sprintf("%d + k~%s v~%s\n", tx.id, item.Key, item.Value)
		} else {
			s = fmt.Sprintf("%d - k~%s\n", tx.id, key)
		}

		buf = append(buf, s...)
	}

	_, err := tx.db.file.Write(buf)
	return err
}
