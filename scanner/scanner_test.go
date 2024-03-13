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

	tokens, errs := sc.Scan()

	if len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		t.FailNow()
	}

	for _, t := range tokens {
		fmt.Println(t.String())
	}
}

func _TestPrint(t *testing.T) {
}
