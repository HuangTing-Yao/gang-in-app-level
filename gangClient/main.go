package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	jobName := os.Getenv("jobName")
	serviceName := os.Getenv("serviceName")
	memberAmount, _ := strconv.Atoi(os.Getenv("memberAmount"))
	serviceName = strings.ToUpper(serviceName) + "_SERVICE_HOST"
	serviceIP := os.Getenv(serviceName)
	// flag.Parse()
	addRequest(serviceIP, jobName)
	// Check if we satisfy gang minMember or not every 5 second
	for {
		number := checkRequest(serviceIP, jobName)
		if number >= memberAmount {
			fmt.Println("satisfy gang minMember.")
			fmt.Println("start to run job.")
			time.Sleep(60 * time.Second) // means the application start running job
			break
		}
		time.Sleep(2 * time.Second)
	}
}

func addRequest(ip string, jobName string) {
	site := "http://" + ip + ":8863" + "/ws/v1/add/" + jobName
	resp, err := http.Get(site)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Get request body fail.")
	}
	fmt.Println(string(body))
}

func checkRequest(ip string, jobName string) int {
	site := "http://" + ip + ":8863" + "/ws/v1/check/" + jobName
	resp, err := http.Get(site)
	var value int
	if err != nil {
		fmt.Println("Check jobMember fail.")
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&value)
	if err != nil {
		fmt.Println("Decode fail.")
	}
	return value
}
