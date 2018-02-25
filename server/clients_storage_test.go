package server_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/christiangda/ingester/server"
)

func TestClientsStore_CountSet(t *testing.T) {
	c1 := &server.Client{}
	cs := server.NewClientStore()
	cs.Set(uuid.New(), c1)
	expected := 1
	actual := cs.Count()
	if actual != expected {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestClientsStore_CountPut(t *testing.T) {
	c1 := &server.Client{}
	cs := server.NewClientStore()
	cs.Put(c1)
	expected := 1
	actual := cs.Count()
	if actual != expected {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}

func TestClientsStore_CountPutAndDelete(t *testing.T) {
	c1 := &server.Client{}
	cs := server.NewClientStore()
	key := cs.Put(c1)
	cs.Delete(key)
	expected := 0
	actual := cs.Count()
	if actual != expected {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}
