package server_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/christiangda/ingester/server"
)

func TestConnectionStore_CountSet(t *testing.T) {
	c1 := &server.Connection{}
	cs := server.NewConnectionStore()
	cs.Set(uuid.New(), c1)
	expected := 1
	actual := cs.Count()
	if actual != expected {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestConnectionStore_CountPut(t *testing.T) {
	c1 := &server.Connection{}
	cs := server.NewConnectionStore()
	cs.Put(c1)
	expected := 1
	actual := cs.Count()
	if actual != expected {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestConnectionStore_CountPutAndDelete(t *testing.T) {
	c1 := &server.Connection{}
	cs := server.NewConnectionStore()
	key := cs.Put(c1)
	cs.Delete(key)
	expected := 0
	actual := cs.Count()
	if actual != expected {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}
