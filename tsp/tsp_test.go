package tsp

import "testing"

func testpoll() ([]KV, int, string) {
    var pairs []KV
    return pairs, 1234567890111, "namespace"
}

func TestTsp(t *testing.T) {
    want := "v0.5"

    if got := Version; got != want {
        t.Errorf("Version = %q, want %q", got, want)
    }
}

func TestPoller(t *testing.T) {
    client := NewClient("localhost", 33030)
    poller := NewPoller(client, 60.0, nil, nil)
    poller.Add(testpoll)
    poller.Add(func()([]KV, int, string) {
        var pairs []KV
        return pairs, 1234567890222, "key"
    })
}

func TestKV(t *testing.T) {
    kv := NewKV("key", float64(-1.0))
    if kv.key != "key" || kv.value != -1.0 {
        t.Errorf("kv mismatch - key: %s, value: %0.1f", kv.key, kv.value)
    }
}

func TestTV(t *testing.T) {
    tv := NewTV(1234567890111, -1.0)
    if tv.timestamp != 1234567890111 || tv.value != -1.0 {
        t.Errorf("tv mismatch - tv.timestamp: %d, tv.value: %0.1f", tv.timestamp, tv.value)
    }
}

