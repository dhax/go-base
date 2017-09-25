package bom

import (
	"bytes"
	"io/ioutil"
	"testing"
)

var (
	test_cases = []struct {
		src      []byte
		expected []byte
	}{
		{
			[]byte{bom0, bom1, bom2}, []byte{},
		},
		{
			[]byte{0x33, bom0, bom1, bom2}, []byte{0x33, bom0, bom1, bom2},
		},
		{
			[]byte{bom0, bom1}, []byte{bom0, bom1},
		},
		{
			[]byte{bom0, bom2}, []byte{bom0, bom2},
		},
		{
			[]byte{bom0, bom1, bom2, 0x11}, []byte{0x11},
		},
	}
)

func TestNewReaderWithoutBom(t *testing.T) {
	for index, test_case := range test_cases {
		result, err := NewReaderWithoutBom(bytes.NewReader(test_case.src))
		if err != nil {
			t.Fatal(err)
		}
		bs, err := ioutil.ReadAll(result)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(bs, test_case.expected) != 0 {
			t.Fatalf("The %d th test_case failed", index)
		}
	}

}
func TestCleanBom(t *testing.T) {
	for index, test_case := range test_cases {
		result := CleanBom(test_case.src)
		if bytes.Compare(result, test_case.expected) != 0 {
			t.Fatalf("The %d th test_case failed", index)
		}
	}
}
