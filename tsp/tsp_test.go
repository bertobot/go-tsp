package tsp

import "testing"

func testpoll() ([]KV, int, string) {
    var pairs []KV
    return pairs, 1234567890111, "namespace"
}

func TestTsp(t *testing.T) {
    want := "v0.2"

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

