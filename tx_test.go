package xisdb

import (
	"fmt"
	"testing"

	"github.com/alexsward/xisdb/indexes"
)

func TestTxClose(t *testing.T) {
	fmt.Println("-- TestTxClose")
	tx := NewTransaction(false, openTestDB())
	if tx.closed {
		t.Errorf("Expected tx to not be closed, was")
	}
	tx.close()
	if len(tx.hooks) != 0 {
		t.Errorf("Expected 0 hooks after close, got %d", len(tx.hooks))
	}
	if len(tx.rollbacks) != 0 {
		t.Errorf("Expected 0 rollbacks after close, got %d", len(tx.rollbacks))
	}
	if len(tx.rollbackBuckets) != 0 {
		t.Errorf("Expected 0 rollbackBuckets after close, got %d", len(tx.rollbackBuckets))
	}
	if tx.db != nil {
		t.Errorf("Expected nil DB for transaction, wasn't")
	}
	if !tx.closed {
		t.Errorf("Expected transaction to be closed, wasn't")
	}
}

func TestTxBucketCreate(t *testing.T) {
	fmt.Println("-- TestTxBucketCreate")
	tests := []struct {
		err       error
		db, write bool
	}{
		{ErrNoDatabase, false, false},
		{ErrNotWriteTransaction, true, false},
		{nil, true, true},
	}
	for i, test := range tests {
		var db *DB
		if test.db {
			db = openTestDB()
		}
		tx := NewTransaction(test.write, db)
		_, err := tx.Bucket("name")
		if err != test.err {
			t.Errorf("Test %d failed: expected error %s, got %s", i+1, test.err, err)
		}
	}
}

func TestTxBucketDelete(t *testing.T) {
	fmt.Println("-- TextTxBucketDelete")
	tests := []struct {
		err            error
		name           string
		db, write      bool
		create, expect bool
	}{
		{nil, "b1", true, true, true, true},
		{nil, "b1", true, true, false, false},
		{ErrNoDatabase, "b1", false, true, false, false},
		{ErrNotWriteTransaction, "b1", true, false, false, false},
	}
	for i, test := range tests {
		var db *DB
		if test.db {
			db = openTestDB()
		}
		tx := NewTransaction(test.write, db)
		if test.create {
			db.buckets[test.name] = newBucket(test.name, db)
		}
		result, err := tx.DeleteBucket(test.name)
		if err != test.err {
			t.Errorf("Test %d failed: expected error %s, got %s", i+1, test.err, err)
			continue
		}
		if test.err != nil {
			continue
		}
		if result != test.expect {
			t.Errorf("Test %d failed: expected delete result:%t, got %t", i+1, test.expect, result)
		}
	}
}

func TestTxBuckets(t *testing.T) {
	fmt.Println("-- TestTxBuckets")
	tests := []struct {
		db, write bool
		err       error
		adds      []string
		expected  []string
	}{
		{true, true, nil, []string{}, []string{""}},
		{true, true, nil, []string{"b1"}, []string{"", "b1"}},
		{true, true, nil, []string{"b1", "b2"}, []string{"", "b1", "b2"}},
		{false, true, ErrNoDatabase, []string{}, []string{""}},
		{true, false, ErrNotWriteTransaction, []string{}, []string{""}},
	}
	for i, test := range tests {
		var db *DB
		if test.db {
			db = openTestDB()
		}
		tx := NewTransaction(test.write, db)
		for _, add := range test.adds {
			db.buckets[add] = newBucket(add, db)
		}
		buckets, err := tx.Buckets()
		if err != test.err {
			t.Errorf("Test %d failed: expected error '%s', got '%s'", i+1, test.err, err)
			continue
		}
		if test.err != nil {
			continue
		}
		if len(buckets) != len(test.expected) {
			t.Errorf("Test %d failed: Expected %d buckets, got %d", i+1, len(test.expected), len(buckets))
			continue
		}
		if !buckets[0].managed.isRoot() {
			t.Errorf("Test %d failed: Expected first bucket to be root, it wasn't, it was: '%s'", i+1, buckets[0].managed.name)
			continue
		}
		for _, bucket := range buckets[1:] {
			found := false
			for _, expected := range test.expected {
				if expected == bucket.managed.name {
					found = true
				}
			}
			if !found {
				t.Errorf("Test %d failed: didn't find %s bucket but it was returned", i+1, bucket.managed.name)
			}
		}
	}
}

func TestTxGet(t *testing.T) {
	fmt.Println("-- TestTxGet")
	tests := []struct {
		set, value, get string
		close           bool
		err             error
	}{
		{"key", "value", "key", false, nil},
		{"key", "value", "key", true, ErrNoDatabase},
		{"key", "value", "unknown", false, ErrKeyNotFound},
	}
	for i, test := range tests {
		db := openTestDB()
		db.root().data[test.set] = Item{Key: test.set, Value: test.value}
		tx := NewTransaction(false, db)
		if test.close {
			tx.close()
		}
		got, err := tx.Get(test.get)
		if err != test.err {
			t.Errorf("Test %d failed: expected error '%s', got '%s'", i+1, test.err, err)
			continue
		}
		if test.err != nil {
			continue
		}
		if got != test.value {
			t.Errorf("Test %d failed: Expected value '%s', got '%s'", i+1, test.value, got)
		}
	}
}

func TestTxExists(t *testing.T) {
	fmt.Println("-- TestTxExists")
	tests := []struct {
		key                string
		close, add, exists bool
		err                error
	}{
		{"key", false, true, true, nil},
		{"key", false, false, false, nil},
		{"key", true, false, false, ErrNoDatabase},
	}
	for i, test := range tests {
		db := openTestDB()
		tx := NewTransaction(true, db)
		if test.close {
			tx.close()
		}
		if test.add {
			db.Set(test.key, "test")
		}
		exists, err := tx.Exists(test.key)
		if err != test.err {
			t.Errorf("Test %d failed: expected error '%s', got '%s'", i+1, test.err, err)
			continue
		}
		if test.err != nil {
			continue
		}
		if test.exists != exists {
			t.Errorf("Test %d failed: expected exists:%t, got %t", i+1, test.exists, exists)
		}
	}
}

func TestTxSet(t *testing.T) {
	fmt.Println("-- TestTxSet")
	tests := []struct {
		set, value   string
		write, close bool
		err          error
	}{
		{"key", "value", true, false, nil},
		{"key", "value", true, true, ErrNoDatabase},
		{"key", "value", false, false, ErrNotWriteTransaction},
	}
	for i, test := range tests {
		db := openTestDB()
		tx := NewTransaction(test.write, db)
		if test.close {
			tx.close()
		}
		err := tx.Set(test.set, test.value, nil)
		if err != test.err {
			t.Errorf("Test %d failed: expected error '%s', got '%s'", i+1, test.err, err)
			continue
		}
		if test.err != nil {
			continue
		}
		val, ok := db.root().data[test.set]
		if !ok {
			t.Errorf("Test %d failed: key not found", i+1)
		}
		if val.Value != test.value {
			t.Errorf("Test %d failed: expected value '%s', got '%s'", i+1, test.value, val.Value)
		}
	}
}

func TestTxDelete(t *testing.T) {
	fmt.Println("-- TestTxDelete")
	tests := []struct {
		key, value, remove    string
		write, close, deleted bool
		err                   error
	}{
		{"key", "value", "key", true, false, true, nil},
		{"key", "value", "key", true, true, false, ErrNoDatabase},
		{"key", "value", "key", false, false, true, ErrNotWriteTransaction},
		{"key", "value", "key2", true, false, false, ErrKeyNotFound},
	}
	for i, test := range tests {
		db := openTestDB()
		db.Set(test.key, test.value)
		tx := NewTransaction(test.write, db)
		if test.close {
			tx.close()
		}
		deleted, err := tx.Delete(test.remove)
		if err != test.err {
			t.Errorf("Test %d failed: expected error '%s', got '%s'", i+1, test.err, err)
			continue
		}
		if test.err != nil {
			continue
		}
		if deleted != test.deleted {
			t.Errorf("Test %d failed: expected deleted:%t, got:%t", i+1, test.deleted, deleted)
			continue
		}
		_, found := db.root().data[test.remove]
		if found {
			t.Errorf("Expected to not still find the data")
		}
	}
}

func TestTxAddRollback(t *testing.T) {
	fmt.Println("-- TestTxAddRollback")
	db := openTestDB()
	tx := NewTransaction(false, db)
	tx.addRollback("", "key", nil)
	if len(tx.rollbacks) != 0 {
		t.Errorf("Expected no rollbacks added to read-only transaction, got 1")
	}
	tx = NewTransaction(true, db)
	tx.addRollback("bucket", "key", nil)
	tx.addRollback("bucket", "key", &Item{"key", "value1", nil})
	tx.addRollback("bucket", "key", &Item{"key", "value2", nil})
	if len(tx.rollbacks) != 1 {
		t.Errorf("Expected single rollback bucket, got %d", len(tx.rollbacks))
	}
}

func TestTxAddIndexErrors(t *testing.T) {
	fmt.Println("-- TestTxAddIndexErrors")
	db := openTestDB()
	tx := NewTransaction(true, db)
	tx.close()
	err := tx.AddIndex("", KeyIndex, nil, nil)
	if err != ErrNoDatabase {
		t.Errorf("Expected error adding index to closed transaction: '%s', got '%s'", ErrNoDatabase, err)
	}
	tx = NewTransaction(false, db)
	err = tx.AddIndex("index", KeyIndex, nil, nil)
	if err != ErrNotWriteTransaction {
		t.Errorf("Expected error adding index to closed transaction: '%s', got '%s'", ErrNotWriteTransaction, err)
	}
	tx = NewTransaction(true, db)
	err = tx.AddIndex("same", KeyIndex, nil, nil)
	err = tx.AddIndex("same", KeyIndex, nil, nil)
	if err != ErrIndexAlreadyExists {
		t.Errorf("Expected error adding index that already exists transaction: '%s', got '%s'", ErrIndexAlreadyExists, err)
	}
}

// TestTxDeleteIndex tests removal of indexes within a transaction
func TestTxDeleteIndex(t *testing.T) {
	fmt.Println("-- TestTxDeleteIndex")
	tests := []struct {
		name           string
		init, expected bool
		err            error
	}{
		{"exists", true, true, nil},
		{"doesnt-exist", false, false, nil},
	}
	for i, test := range tests {
		db := openTestDB()
		if test.init {
			db.root().indexes[test.name] = &index{}
		}
		removed, err := db.DeleteIndex(test.name)
		if err != test.err {
			t.Errorf("Test %d failed: expected error '%s', got '%s'", i+1, test.err, err)
			continue
		}
		if test.err != nil {
			continue
		}

		if removed != test.expected {
			t.Errorf("Test %d failed: expected removal %t, got %t", i+1, test.expected, removed)
		}
	}
}

func TestTxAddIndexKeyWildCard(t *testing.T) {
	fmt.Println("-- TestTxAddIndexKeyWildCard")
	tests := []struct {
		add []string
	}{
		{[]string{}},
		{[]string{"a", "abc", "1234", "some data", "json?", "{}"}},
	}
	for i, test := range tests {
		db := openTestDB()
		tx := NewTransaction(true, db)
		for _, k := range test.add {
			db.Set(k, "value")
		}
		err := tx.AddIndex("test", KeyIndex, indexes.WildcardMatcher, NaturalOrderKeyComparison)
		if err != nil {
			t.Errorf("Test %d failed: didn't expect error, got: '%s'", i+i, err)
			continue
		}
		idx, exists := db.root().indexes["test"]
		if !exists {
			t.Errorf("Index didn't exist after Add")
			continue
		}
		if idx.tree.Size() != uint(len(test.add)) {
			t.Errorf("Test %d failed: index expected to have %d items, had: %d", i+1, idx.tree.Size(), len(test.add))
		}
	}
}
