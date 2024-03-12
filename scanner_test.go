package scanner

import (
	"fmt"
	"testing"
)

func TestScanFile(t *testing.T) {
	sc, err := NewScanner("./fixture/file.lox")

	if err != nil {
		t.Fatal(err)
	}

	for _, t := range sc.ScanToken() {
		fmt.Println(t.String())
	}
}
