package statsd

import (
	"reflect"
	"testing"
)

func TestClient_NewClient(t *testing.T) {
	client := NewClient("127.0.0.1", 9876)
	typeName := reflect.TypeOf(*client).Name()
	if typeName != "Client" {
		t.Errorf("Wrong type of factory method return value: \"%s\"", typeName)
	}
}

func TestClient_Timing(t *testing.T) {
	client := NewClient("127.0.0.1", 9876)
	client.Timing("a.b.c", 320, 0.954)

	if client.keyBuffer["a.b.c"] != "320|ms|@0.954" {
		t.Errorf("Wrong timing metric: \"%s\"", client.keyBuffer["a.b.c"])
	}
}

func TestClient_Count(t *testing.T) {
	client := NewClient("127.0.0.1", 9876)
	client.Count("a.b.c", 320, 0.954)

	if client.keyBuffer["a.b.c"] != "320|c|@0.954" {
		t.Errorf("Wrong timing metric: \"%s\"", client.keyBuffer["a.b.c"])
	}
}
