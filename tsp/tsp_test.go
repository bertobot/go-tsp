package tsp

import "testing"

func TestTsp(t *testing.T) {
    want := "v0.1"

    if got := Version; got != want {
        t.Errorf("Version = %q, want %q", got, want)
    }
}

