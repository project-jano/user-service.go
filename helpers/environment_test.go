package helpers

import (
	"os"
	"testing"
)

const (
	testKey      = "test-key"
	defaultValue = "default-value"
	anotherValue = "another-value"
)

func TestMissingKey(t *testing.T) {
	os.Clearenv()

	value := GetEnvVar(testKey, defaultValue)
	if value != defaultValue {
		t.Fail()
	}
}

func TestHasEnvVar(t *testing.T) {
	os.Clearenv()
	os.Setenv(testKey, anotherValue)

	value := GetEnvVar(testKey, defaultValue)
	if value != anotherValue {
		t.Fail()
	}
}
