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
	client.Timing("a.b.c", 320, 0.9999)

	if client.keyBuffer["a.b.c"] != "320|ms|@0.9999" {
		t.Errorf("Wrong timing metric: \"%s\"", client.keyBuffer["a.b.c"])
	}
}

func TestClient_Timing_NoRate(t *testing.T) {
	client := NewClient("127.0.0.1", 9876)
	client.Timing("a.b.c", 320, 1)

	if client.keyBuffer["a.b.c"] != "320|ms" {
		t.Errorf("Wrong timing metric: \"%s\"", client.keyBuffer["a.b.c"])
	}
}

func TestClient_Count(t *testing.T) {
	client := NewClient("127.0.0.1", 9876)
	client.Count("a.b.c", 320, 0.9999)

	if client.keyBuffer["a.b.c"] != "320|c|@0.9999" {
		t.Errorf("Wrong timing metric: \"%s\"", client.keyBuffer["a.b.c"])
	}
}

func TestClient_Count_NoRate(t *testing.T) {
	client := NewClient("127.0.0.1", 9876)
	client.Count("a.b.c", 320, 1)

	if client.keyBuffer["a.b.c"] != "320|c" {
		t.Errorf("Wrong timing metric: \"%s\"", client.keyBuffer["a.b.c"])
	}
}

