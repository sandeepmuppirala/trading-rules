package main

import (
	"bufio"
	//"bytes"
	"fmt"
	"io/ioutil"
	//"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"trading-rules/redis"
	"trading-rules/service"
	"trading-rules/utils"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// TestDataLoadSequence : Trigger business rules and perform data load
func TestDataLoadSequence(t *testing.T) {
	dataloadInput, err := ioutil.ReadFile("input_test.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	//r := httptest.NewRequest("GET", "/foo", nil)
	w := httptest.NewRecorder()

	customerInputInfo := strings.Split(string(dataloadInput), "\n")
	dataLoadSequence := service.DataLoadSequence(w, customerInputInfo)
	assert.True(t, dataLoadSequence, "Del operation did not work as expected")
}

// TestProcessLoadedFile : Process input data file that is loaded
func TestProcessLoadedFile(t *testing.T) {
	dataloadInput, err := ioutil.ReadFile("input_test.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	w := httptest.NewRecorder()
	clearPreviousLoads := ProcessLoadedFile(w, dataloadInput)
	assert.True(t, clearPreviousLoads, "Processing loaded file fails")
}

func Router() *mux.Router {
	router:= mux.NewRouter()
	router.HandleFunc("/v1/dataload", TradingRules).Methods("POST")
	return router
}

// TestFileReader : Test dataload request
func TestFileReader(t *testing.T) {
	file, err := os.Open("input_test.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	request, _ := http.NewRequest("POST", "/v1/dataload", bufio.NewReader(file))
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
}

// TestGetAndSetOperation : Redis Set and Get key tests
func TestGetAndSetOperation(t *testing.T) {
	redis.SetOperation("sample_custID", "sample data")
	custIdValue, _ := redis.GetOperation("sample_custID")
	assert.Equal(t, "sample data", custIdValue, "Get/Set operation did not work as expected")
}

// TestDelOperation : Redis Delete key test
func TestDelOperation(t *testing.T) {
	redis.DelOperation("sample_custID", "")
	custIdValueAfterDeletion, _ := redis.GetOperation("sample_custID")
	assert.Equal(t, "", custIdValueAfterDeletion, "Del operation did not work as expected")
}

// TestIsInputDateWithinRedisEndOfWeek : Tests if the input date provided is within required limits
func TestIsInputDateWithinRedisEndOfWeek(t *testing.T) {
	isWithinEndOfWeek := utils.IsInputDateWithinRedisEndOfWeek("2021-01-26", "2021-01-30")
	assert.True(t, isWithinEndOfWeek, "Del operation did not work as expected")
}

// TestClearPreviousLoadAttempt : Clears previous load attempts if any
func TestClearPreviousLoadAttempt(t *testing.T) {
	dataloadInput, err := ioutil.ReadFile("input_test.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	customerInputInfo := strings.Split(string(dataloadInput), "\n")
	clearPreviousLoads := service.ClearPreviousLoadAttempt(customerInputInfo)
	assert.True(t, clearPreviousLoads, "Del operation did not work as expected")
}