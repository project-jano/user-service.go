package helpers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
 * Project Jano - User microservice
 * This is the API of Project Jano
 *
 * API version: 2.0.4
 * Contact: ezequiel.aceto+project-jano@gmail.com

 */

func TestExtractKeyId(t *testing.T) {
	r := createTestRequestWithQueryParam(QueryParamKeyId, "")

	defKey, keyId := ExtractKeyId(r)
	if !defKey {
		t.Fail()
	}
	if keyId != DefaultKeyId {
		t.Fail()
	}
}

func TestExtractDefaultKey(t *testing.T) {
	testKeyId := "test-key-id"
	r := createTestRequestWithQueryParam(QueryParamKeyId, testKeyId)

	defKey, keyId := ExtractKeyId(r)
	if defKey {
		t.Fail()
	}
	if testKeyId != keyId {
		t.Fail()
	}
}

func TestSetupContentType(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		SetupDefaultContentType(w)
	}

	request := httptest.NewRequest("GET", "http://test/action", nil)
	w := httptest.NewRecorder()
	handler(w, request)

	resp := w.Result()

	if resp.Header.Get(ContentType) != DefaultContentType {
		t.Fail()
	}

}

func TestUseDefaultDevices(t *testing.T) {

	request := httptest.NewRequest("GET", "http://test/action", nil)

	useAllDevices, useDefaultDevice, devices, ok := ExtractDevicesFilters(request)

	if useAllDevices != false {
		t.Fail()
	}
	if useDefaultDevice != true {
		t.Fail()
	}
	if len(devices) != 0 {
		t.Fail()
	}
	if !ok {
		t.Fail()
	}
}

func TestUseCustomDeviceId(t *testing.T) {

	request := httptest.NewRequest("GET", "http://test/action?deviceIds=custom", nil)

	useAllDevices, useDefaultDevice, devices, ok := ExtractDevicesFilters(request)

	if useAllDevices != false {
		t.Fail()
	}
	if useDefaultDevice != false {
		t.Fail()
	}
	if len(devices) != 1 {
		t.Fail()
	}
	if !ok {
		t.Fail()
	}
}


func TestUseTwoDeviceIds(t *testing.T) {

	request := httptest.NewRequest("GET", "http://test/action?deviceIds=custom,emulator", nil)

	useAllDevices, useDefaultDevice, devices, ok := ExtractDevicesFilters(request)

	if useAllDevices != false {
		t.Fail()
	}
	if useDefaultDevice != false {
		t.Fail()
	}
	if len(devices) != 2 {
		t.Fail()
	}
	if !ok {
		t.Fail()
	}
}


func TestUseAllDevices(t *testing.T) {

	request := httptest.NewRequest("GET", "http://test/action?deviceIds=all", nil)

	useAllDevices, useDefaultDevice, devices, ok := ExtractDevicesFilters(request)

	if useAllDevices != true {
		t.Fail()
	}
	if useDefaultDevice != false {
		t.Fail()
	}
	if len(devices) != 0 {
		t.Fail()
	}
	if !ok {
		t.Fail()
	}
}

func TestDeviceIdEmptyQuery(t *testing.T) {

	request := httptest.NewRequest("GET", "http://test/action?deviceIds=", nil)

	useAllDevices, useDefaultDevice, devices, ok := ExtractDevicesFilters(request)

	if useAllDevices != false {
		t.Fail()
	}
	if useDefaultDevice != true {
		t.Fail()
	}
	if len(devices) != 0 {
		t.Fail()
	}
	if !ok {
		t.Fail()
	}
}

func createTestRequestWithQueryParam(key string, value string) *http.Request {
	request := httptest.NewRequest("GET", "http://test/action", nil)

	q := request.URL.Query()
	q.Add(key, value)
	request.URL.RawQuery = q.Encode()

	return request
}
