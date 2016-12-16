package xisdb

import (
	"fmt"
	"testing"
	"time"
)

// TestSubscribe validates that subscriptions work
func TestPubSubSubscribe(t *testing.T) {
	fmt.Println("--- TestPubSubSubscribe")
	db := openTestDB()
	_, err := db.Subscribe("*", 1)
	if err != nil {
		t.Error(err)
	}
}

// TestUnsubscribe validates that unsubscribing works
func TestPubSubUnsubscribe(t *testing.T) {
	fmt.Println("--- TestPubSubUnsubscribe")
	db := openTestDB()
	ch, err := db.Subscribe("*", 1)
	if err != nil {
		t.Error(err)
	}

	err = db.Unsubscribe("*", ch)
	if err != nil {
		t.Error(err)
	}
}

// TestPubSubPublish validates that proper items are received over a channel
func TestPubSubPublish(t *testing.T) {
	fmt.Println("--- TestPubSubPublish")
	db := openTestDB()
	ch, err := db.Subscribe("pubsub:", 1)
	if err != nil {
		t.Error(err)
	}

	db.Set("pubsub:test", "some_value")

	select {
	case item := <-ch:
		if item.Key != "pubsub:test" {
			t.Errorf("Expected [pubsub:test], got [%s]", item.Key)
		}
		return
	case <-time.After(1 * time.Second):
		t.Error("Test timed out")
	}
}
