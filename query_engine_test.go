package xisdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/alexsward/xisdb/ql"
)

func TestQueryEngineExecuteGet(t *testing.T) {
	fmt.Println("-- TestQueryEngine")
	db := openTestDB()
	db.Set("key", "value")
	ch := make(chan Item, 0)
	qe := QueryEngine{}
	err := qe.Execute([]ql.Statement{createSimpleGet("key")}, &QueryEngineContext{db, ch})
	if err != nil {
		t.Errorf("Test failed. Error executing statemnt: %s", err)
		return
	}
	time.Sleep(time.Millisecond * 50)
}

func createSimpleGet(key string) ql.Statement {
	s, _ := ql.Parse(fmt.Sprintf("get %s;", key))
	return s[0]
}

func createQueryEngineContext() *QueryEngineContext {
	return &QueryEngineContext{
		DB: openTestDB(),
	}
}
