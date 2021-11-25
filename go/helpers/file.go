package helpers

import (
	"io/ioutil"
	"log"
)

func ReadFile(path string) string {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading file: %s, %v", path, err)
	}
	return string(raw)
}
