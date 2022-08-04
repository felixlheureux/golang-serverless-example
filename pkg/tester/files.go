package tester

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func ReadFile(t *testing.T, paths ...string) []byte {
	t.Helper()
	path := filepath.Join(paths...)

	b, err := ioutil.ReadFile(path)

	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return b
}
