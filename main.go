package main

import (
	"net/http"
	"strings"
	"time"
	"trading-rules/constants"
	"trading-rules/redis"
	"trading-rules/service"
	"trading-rules/utils"

	log "github.com/sirupsen/logrus"
)

func init() {
	redis.InitRedisConnection()
}

func main() {
	handler := http.NewServeMux()
	handler.HandleFunc(constants.ControllerMapping, TradingRules)
	http.ListenAndServe(constants.DeployedPath, handler)
}

// TradingRules : Controller to initiate data load operation
func TradingRules(w http.ResponseWriter, r *http.Request) {
	log.Info("Starting data load")

	// Parse input files into bytes
	fileInputByte, err := utils.FileReader(r)
	if err != nil {
		log.Error(err)
		return
	}

	ProcessLoadedFile(w, fileInputByte)
}

// ProcessLoadedFile : Process input data file
func ProcessLoadedFile(w http.ResponseWriter, fileInputByte []byte) bool {
	processFileLoad := false

	customerInputInfo := strings.Split(string(fileInputByte), "\n")

	// Clear any Redis keys from previous data loads
	service.ClearPreviousLoadAttempt(customerInputInfo)

	// Sleep time for two seconds between Delete operation and Get/Set operations
	time.Sleep(2000)

	// Start loading data into redis while performing calculations
	service.DataLoadSequence(w, customerInputInfo)

	SetResponseHeaders(&w)

	log.Info("Data load completed")

	processFileLoad = true

	return processFileLoad
}

// SetResponseHeaders : Set the required standard response headers
func SetResponseHeaders(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
