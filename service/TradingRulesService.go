package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"trading-rules/constants"
	"trading-rules/models"
	"trading-rules/redis"
	"trading-rules/utils"

	log "github.com/sirupsen/logrus"
)

// ClearPreviousLoadAttempt : Clears previous redis keys
func ClearPreviousLoadAttempt(customerInputInfo []string) bool {
	previousLoadsCleared := false
	for _, customerInput := range customerInputInfo {
		inputDataToDelete := models.DataRulesInput {}
		_ = json.Unmarshal([]byte(customerInput), &inputDataToDelete)
		_ = redis.DelOperation(inputDataToDelete.CustomerID, "")
	}
	previousLoadsCleared = true
	return previousLoadsCleared
}

// DataLoadSequence : Initiate data load sequence
func DataLoadSequence (w http.ResponseWriter, customerInputInfo []string) bool {
	dataLoadSequenceCompleted := false
	for _, customerInput := range customerInputInfo {
		inputDataToLoad := models.DataRulesInput {}
		_ = json.Unmarshal([]byte(customerInput), &inputDataToLoad)
		// If Customer ID is blank, we do not process
		if len(inputDataToLoad.CustomerID) == 0 {
			continue
		} 
		dataLoadValidity := ProcessDataLoad(inputDataToLoad)
		lStatus := utils.GenerateJSONResponse(inputDataToLoad, dataLoadValidity)
		// Printing to the console
		fmt.Println(string(lStatus))
		// Printing to the response writer
		_, _ = fmt.Fprintln(w, string(lStatus))
	}
	dataLoadSequenceCompleted = true
	return dataLoadSequenceCompleted
}

// ProcessDataLoad : Business logic to process data load
func ProcessDataLoad(inputData models.DataRulesInput) bool {
	
	// Initiate validity of data load
	isCurrentDataLoadValid := false

	// Format for storing input data: ID||LoadAmount||Time||AttemptsPerDay||DataLoadedPerWeek	
	formattedInputData := inputData.ID + constants.Delimiter + 
			inputData.LoadAmount + constants.Delimiter + 
			inputData.Time + constants.Delimiter + 
			strconv.Itoa(inputData.AttemptsPerDay) + constants.Delimiter + 
			inputData.LoadAmount

	// Check if customer exists in Redis
	customerInfoFromRedis,err := redis.GetOperation(inputData.CustomerID)
	
	// Data to persist as per the applicable rules
	dataToPersist := GetDataToPersistBasedOnRules(formattedInputData, customerInfoFromRedis)

	// Set customer info in Redis if applicable
	if len(dataToPersist) > 0 {
		err =  redis.SetOperation(inputData.CustomerID, dataToPersist)
		if err != nil {
			panic(err)
		}
		isCurrentDataLoadValid = true
	}
	return isCurrentDataLoadValid
}

// GetDataToPersistBasedOnRules : Run input data against defined rules to generate data to persist
func GetDataToPersistBasedOnRules(inputData string, customerInfoFromRedis string) string {
	dataToPersist := ""
	inputDataSplit:= strings.Split(inputData, constants.Delimiter)
	loadAmountForNewDay, _ := strconv.ParseFloat(strings.Replace(inputDataSplit[1], constants.Currency, "", -1), 64)
	
	if len(customerInfoFromRedis) > 0 {
		log.Debug("Customer exists in Redis: " , customerInfoFromRedis)
		customerInfoFromRedisSlice:= strings.Split(customerInfoFromRedis, constants.Delimiter)
		
		dayAsPerDataLoad := strings.Split(inputDataSplit[2], "T")[0]
		dayAsPerRedis := strings.Split(customerInfoFromRedisSlice[2], "T")[0]

		numberOfTransactionsPerDay, _ := strconv.Atoi(customerInfoFromRedisSlice[3])

		loadAmountForExistingDayFromRedis, _ := strconv.ParseFloat(strings.Replace(customerInfoFromRedisSlice[1], constants.Currency, "", -1), 64) 
		loadAmountForExistingDay := loadAmountForExistingDayFromRedis + loadAmountForNewDay
		
		weeklyCountFromRedis, _ := strconv.ParseFloat(strings.ReplaceAll(customerInfoFromRedisSlice[4], constants.Currency, ""), 64)
		weeklyCountFromInput, _ := strconv.ParseFloat(strings.ReplaceAll(inputDataSplit[4], constants.Currency, ""), 64)
		weeklyCount := weeklyCountFromRedis + weeklyCountFromInput
		
		if numberOfTransactionsPerDay + 1 < 3 && 
					dayAsPerRedis == dayAsPerDataLoad && 
							loadAmountForExistingDay < 5000 && weeklyCount < 20000 {
			log.Debug("New transaction for the same day")		

			numberOfTransactionsPerDay = numberOfTransactionsPerDay + 1
			var iterateCountToPersist = strconv.Itoa(numberOfTransactionsPerDay)

			dataToPersist = BuildDataToPersistAfterMatchingRules(customerInfoFromRedisSlice[0], 
				loadAmountForExistingDay, inputDataSplit[2], iterateCountToPersist, weeklyCount)
		} else if dayAsPerRedis != dayAsPerDataLoad && 
						utils.IsInputDateWithinRedisEndOfWeek(dayAsPerDataLoad, dayAsPerRedis) &&
								loadAmountForNewDay < 5000 && weeklyCount < 20000 {
			log.Debug("New transaction for another day of the same week")

			dataToPersist = BuildDataToPersistAfterMatchingRules(customerInfoFromRedisSlice[0], 
				loadAmountForNewDay, inputDataSplit[2], "0", weeklyCount)
		} else if dayAsPerRedis != dayAsPerDataLoad && 
						!utils.IsInputDateWithinRedisEndOfWeek(dayAsPerDataLoad, dayAsPerRedis) &&
								loadAmountForNewDay < 5000 && weeklyCount < 20000 {
			log.Debug("First transaction of the week")

			dataToPersist = BuildDataToPersistAfterMatchingRules(customerInfoFromRedisSlice[0], 
				loadAmountForNewDay, inputDataSplit[2], "0", loadAmountForNewDay)
		}
	} else if loadAmountForNewDay < 5000 {	// Because the record doesnt exist in Redis, we just verify the load amount is within limits
		log.Debug("No transaction history")
		dataToPersist = inputData
	}
	return dataToPersist
}

// BuildDataToPersistAfterMatchingRules : Build data to persist after matching input data against the rules 
func BuildDataToPersistAfterMatchingRules(ID string, dayLoadAmount float64, time string, iterateCount string, weekLoadAmount float64) string{
	var loadAmountToPersist = constants.Currency + fmt.Sprintf(constants.FloatingPointPrecision, dayLoadAmount)
	var weeklyCountToPersist = constants.Currency + fmt.Sprintf(constants.FloatingPointPrecision, weekLoadAmount)

	dataToPersist := 	ID + constants.Delimiter + 
						loadAmountToPersist + constants.Delimiter + 
						time + constants.Delimiter +
						iterateCount + constants.Delimiter +
						weeklyCountToPersist

	return dataToPersist
}

