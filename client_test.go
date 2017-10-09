package statsd

import (
	"reflect"
	"testing"
	"net"
	"strings"
)

var udpConnectionStubIO chan []byte = make(chan []byte)
type udpConnectionStub struct {
	net.Conn
}
func (stub *udpConnectionStub) Write(p []byte) (n int, err error) {
	udpConnectionStubIO <- p
	n = 0
	err = nil
	return
}

func TestClient_NewClient(t *testing.T) {
	client := NewClient("127.0.0.1", 9876)
	typeName := reflect.TypeOf(*client).Name()
	if typeName != "Client" {
		t.Errorf("Wrong type of factory method return value: \"%s\"", typeName)
	}
}

func TestClient_SetAutoflush(t *testing.T) {
	client := NewClient("127.0.0.1", 9876)

	client.SetAutoflush(true)
	if client.autoflush != true {
		t.Errorf("Autoflush must be true")
	}

	client.SetAutoflush(false)
	if client.autoflush != false {
		t.Errorf("Autoflush must be false")
	}
}

func TestClient_Timing_BigRate(t *testing.T) {
	client := NewClient("127.0.0.1", 9876)
	client.Timing("a.b.c", 320, 0.9999)

	if len(client.keyBuffer) > 0 && client.keyBuffer["a.b.c"] != "320|ms|@0.9999" {
		t.Errorf("Wrong timing metric with big rate: \"%s\"", client.keyBuffer["a.b.c"])
	}
}

func TestClient_Timing_SmallRate(t *testing.T) {
	client := NewClient("127.0.0.1", 9876)
	client.Timing("a.b.c", 320, 0.0001)

	if len(client.keyBuffer) > 0 && client.keyBuffer["a.b.c"] != "320|ms|@0.0001" {
		t.Errorf("Wrong timing metric with small rate: \"%s\"", client.keyBuffer["a.b.c"])
	}
}

func TestClient_Timing_NoRate(t *testing.T) {
	client := NewClient("127.0.0.1", 9876)
	client.Timing("a.b.c", 320, 1)

	if client.keyBuffer["a.b.c"] != "320|ms" {
		t.Errorf("Wrong timing metric without rate: \"%s\"", client.keyBuffer["a.b.c"])
	}
}

func TestClient_Count_BigRate(t *testing.T) {
	client := NewClient("127.0.0.1", 9876)
	client.Count("a.b.c", 320, 0.9999)

	if len(client.keyBuffer) > 0 && client.keyBuffer["a.b.c"] != "320|c|@0.9999" {
		t.Errorf("Wrong count metric with big rate: \"%s\"", client.keyBuffer["a.b.c"])
	}
}

func TestClient_Count_SmallRate(t *testing.T) {
	client := NewClient("127.0.0.1", 9876)
	client.Count("a.b.c", 320, 0.0001)

	if len(client.keyBuffer) > 0 && client.keyBuffer["a.b.c"] != "320|c|@0.0001" {
		t.Errorf("Wrong count metric with small rate: \"%s\"", client.keyBuffer["a.b.c"])
	}
}

func TestClient_Count_NoRate(t *testing.T) {
	client := NewClient("127.0.0.1", 9876)
	client.Count("a.b.c", 320, 1)

	if client.keyBuffer["a.b.c"] != "320|c" {
		t.Errorf("Wrong count metric without rate: \"%s\"", client.keyBuffer["a.b.c"])
	}
}

func TestClient_Gauge(t *testing.T) {
	client := NewClient("127.0.0.1", 9876)
	client.Gauge("a.b.c", 320)

	if client.keyBuffer["a.b.c"] != "320|g" {
		t.Errorf("Wrong gauge metric: \"%s\"", client.keyBuffer["a.b.c"])
	}
}

func TestClient_GaugeShift(t *testing.T) {
	client := NewClient("127.0.0.1", 9876)

	client.GaugeShift("a.b.c", 320)
	if client.keyBuffer["a.b.c"] != "+320|g" {
		t.Errorf("Wrong gauge metric: \"%s\"", client.keyBuffer["a.b.c"])
	}

	client.GaugeShift("a.b.c", -320)
	if client.keyBuffer["a.b.c"] != "-320|g" {
		t.Errorf("Wrong gauge metric: \"%s\"", client.keyBuffer["a.b.c"])
	}
}

func TestClient_Set(t *testing.T) {
	client := NewClient("127.0.0.1", 9876)
	client.Set("a.b.c", 320)

	if client.keyBuffer["a.b.c"] != "320|s" {
		t.Errorf("Wrong set metric: \"%s\"", client.keyBuffer["a.b.c"])
	}
}

func TestClient_addToBuffer(t *testing.T) {
	client := NewClient("127.0.0.1", 9876)
	client.addToBuffer("a.b.c", "320|s")
	client.addToBuffer("a.b.d", "321|ms|@0.0001")

	if len(client.keyBuffer) != 2 {
		t.Errorf("Must be 2 keys in buffer")
	}

	if client.keyBuffer["a.b.c"] != "320|s" {
		t.Errorf("Wrong metric added to buffer: \"%s\"", client.keyBuffer["a.b.c"])
	}

	if client.keyBuffer["a.b.d"] != "321|ms|@0.0001" {
		t.Errorf("Wrong metric added to buffer: \"%s\"", client.keyBuffer["a.b.c"])
	}
}

func TestClient_isSendAcceptedBySampleRate(t *testing.T) {
	client := NewClient("127.0.0.1", 9876)

	if client.isSendAcceptedBySampleRate(2) == false {
		t.Errorf("2 must be accepred by sample rate")
	}

	if client.isSendAcceptedBySampleRate(1) == false {
		t.Errorf("1 must be accepred by sample rate")
	}

	if client.isSendAcceptedBySampleRate(0.00000001) == true {
		t.Errorf("0.00000001 must not be accepred by sample rate")
	}

	if client.isSendAcceptedBySampleRate(0.99999999) == false {
		t.Errorf("0.99999999 must be accepred by sample rate")
	}
}

func TestClient_Flush(t *testing.T) {
	client := NewClient("127.0.0.1", 9876)

	client.conn = new(udpConnectionStub)

	client.Count("a.a", 42, 1)
	client.Timing("a.b", 43, 1)
	client.Gauge("a.c", 44)
	client.GaugeShift("a.d", 45)
	client.GaugeShift("a.e", 46)
	client.Set("a.f", 47)
	client.Flush()

	metricPacketBytes := <- udpConnectionStubIO
	actualMetricPacket := strings.Replace(
		string(metricPacketBytes),
		"\n",
		"@",
		-1,
	)

	expectedMetricPacket := "a.f:47|s@a.a:42|c@a.b:43|ms@a.c:44|g@a.d:+45|g@a.e:+46|g"

	if expectedMetricPacket != actualMetricPacket {
		t.Errorf(
			"Wrong metric packet send: %s, expected: %s",
			actualMetricPacket,
			expectedMetricPacket,
		)
	}
}

