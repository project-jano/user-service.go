package helpers

import (
	"testing"
)

func TestSplitStringExactly(t *testing.T) {
	str := "AAABBBCCC"
	chunkSize := 3
	splitted := SplitString(str, chunkSize)

	for _, s := range splitted {
		if len(s) != chunkSize {
			t.Fail()
		}
	}
}

func TestSplitString(t *testing.T) {
	str := "AAABBBC"
	chunkSize := 3

	splitted := SplitString(str, chunkSize)

	for _, s := range splitted {
		if len(s) > chunkSize {
			t.Fail()
		}
	}
}

func TestArrayContainsString(t *testing.T) {
	str := "ABC"
	arr := []string{"A", "B", "C", "ABC"}
	if !ContainsStringInStringArray(arr, str) {
		t.Fail()
	}
}

func TestArrayContainsEmptyString(t *testing.T) {
	str := ""
	arr := []string{"A", "B", "C", "ABC"}
	if ContainsStringInStringArray(arr, str) {
		t.Fail()
	}
}

func TestArrayDoesNotContainsString(t *testing.T) {
	str := "D"
	arr := []string{"A", "B", "C"}
	if ContainsStringInStringArray(arr, str) {
		t.Fail()
	}
}
