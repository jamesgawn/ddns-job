package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	response, ipErr := http.Get("https://api.ipify.org?format=text")
	if ipErr != nil {
		fmt.Errorf("failed to retrieve IP address")
		panic(ipErr)
	}
	defer response.Body.Close()
	ip, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		fmt.Errorf("failed to parse IP address")
		panic(readErr)
	}
	fmt.Printf(string(ip))
	_, awsConfigError := config.LoadDefaultConfig(context.TODO())
	if awsConfigError != nil {
		log.Fatal(awsConfigError)
	}
}
