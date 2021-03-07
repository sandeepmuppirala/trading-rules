package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"trading-rules/constants"
	"trading-rules/models"
)

// FileReader : Utility function to process incoming file into bytes
func FileReader(r *http.Request) ([]byte, error) {
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile(constants.FileIdentifier)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer file.Close()

	fileInputByte, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	return fileInputByte, nil
}

// IsInputDateWithinRedisEndOfWeek : Decides if input date is within end of week of load date in Redis
func IsInputDateWithinRedisEndOfWeek(loadDateFromInput string, loadDateFromRedis string) bool{
	var endOfWeek time.Time
	var isWithinEndofWeek bool
	var dateFormat = constants.DateFormat

	parsedDateFromInput, _ := time.Parse(dateFormat, loadDateFromInput)
	parsedDateFromRedis, _ := time.Parse(dateFormat, loadDateFromRedis)

	if parsedDateFromRedis.Weekday() == time.Monday {
		endOfWeek = parsedDateFromRedis.AddDate(0, 0, 7)
	} else {
		end := parsedDateFromRedis.AddDate(0, 0, 6)
		for t := parsedDateFromRedis; t.Before(end); t = t.Add(time.Hour * 24) {
			if t.Weekday() == time.Monday {
				endOfWeek = t
			}
		}
	}
	
	if parsedDateFromInput.Before(endOfWeek) {
		isWithinEndofWeek = true
	} else {
		isWithinEndofWeek = false	
	}
	return isWithinEndofWeek
}

// GenerateJSONResponse : Build JSON response for data load attempt
func GenerateJSONResponse(lData models.DataRulesInput, dataLoadValidity bool) []byte {
	loadStatus, err := json.Marshal(models.DataRulesOutput{
		ID: lData.ID,
		CustomerID:  lData.CustomerID,
		Accepted: dataLoadValidity,
	})
	if err != nil {
		panic(err)
	}
	return loadStatus
}