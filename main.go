package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"io/ioutil"
	"os"
)

type EventPayloadStruct struct {
	Source string `json:"source"`
}

func main() {
	checkEnvironment()
	loadEnvVariables()
	if isLambda {
		lambda.Start(HandleRequest)
	}

	HandleRequest(context.Background(), loadDummyPayload())

}

func HandleRequest(ctx context.Context, eventPayload EventPayloadStruct) (string, error) {
	fmt.Printf("===================================================\n")
	fmt.Printf("Event received from %s\n", eventPayload.Source)
	fmt.Printf("isAWS %t -- isLambda %t -- isDocker %t\n", isAWS, isLambda, isDocker)
	fmt.Printf("%s\n", os.Getenv("TESTVAR"))
	fmt.Printf("===================================================\n")
	fmt.Printf("Made with   ❤️   and  🍝\n")
	fmt.Printf("===================================================\n")

	return "PROCESS COMPLETED", nil
}

var isLambda bool
var isDocker bool
var isAWS bool

// checkEnvironment
func checkEnvironment() {

	//default value
	isLambda = false
	isDocker = false
	isAWS = false

	if len(os.Getenv("AWS_REGION")) != 0 {
		isLambda = true
	}
	//even if we are in a lambda environment, it could be the local docker dev container...
	//so let's use another env var to understand if docker
	//after comparing docker lambda and AWS lambda I noticed that AWS_SESSION_TOKEN env var is (for the moment) only available in AWS
	if isLambda && len(os.Getenv("AWS_SESSION_TOKEN")) == 0 {
		isDocker = true
	}
	//Try to understand if we are running in AWS
	isAWS = isLambda && !isDocker

}

// loadDummyPayload
func loadDummyPayload() EventPayloadStruct {
	var content string
	var eventPayload EventPayloadStruct
	content = readFileContent("./dummyPayload.json")
	_ = json.Unmarshal([]byte(content), &eventPayload)
	return eventPayload
}

// readFileContent .
func readFileContent(filename string) string {
	var content []byte
	var err error
	if fileExists(filename) {
		content, err = ioutil.ReadFile(filename)
		if err != nil {
			fmt.Printf("Unable to find file to load %s\n", filename)
			return "{}"
		}
	}
	return string(content)
}

// fileExists
func fileExists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	}

	return false
}

// loadEnvVariables
func loadEnvVariables() bool {
	//read .env variables
	if err := godotenv.Load(".env"); err != nil {
		return false
	}
	return true
}
