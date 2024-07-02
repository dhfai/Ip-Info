package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type IPInfo struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
}

type UserAgentInfo struct {
	Parse struct {
		SimpleSoftwareString          string   `json:"simple_software_string"`
		SimpleSubDescriptionString    string   `json:"simple_sub_description_string"`
		SimpleOperatingPlatformString string   `json:"simple_operating_platform_string"`
		Software                      string   `json:"software"`
		SoftwareName                  string   `json:"software_name"`
		SoftwareNameCode              string   `json:"software_name_code"`
		SoftwareVersion               string   `json:"software_version"`
		SoftwareVersionFull           []string `json:"software_version_full"`
		OperatingSystem               string   `json:"operating_system"`
		OperatingSystemName           string   `json:"operating_system_name"`
		OperatingSystemNameCode       string   `json:"operating_system_name_code"`
		OperatingSystemFlavour        string   `json:"operating_system_flavour"`
		OperatingSystemFlavourCode    string   `json:"operating_system_flavour_code"`
		OperatingSystemVersion        string   `json:"operating_system_version"`
		OperatingSystemVersionFull    []string `json:"operating_system_version_full"`
		IsAbusive                     bool     `json:"is_abusive"`
		UserAgent                     string   `json:"user_agent"`
	} `json:"parse"`
	Result struct {
		Code        string `json:"code"`
		MessageCode string `json:"message_code"`
		Message     string `json:"message"`
	} `json:"result"`
}

func getIPInfo(ip, apiKey string) (*IPInfo, error) {
	url := fmt.Sprintf("https://ipinfo.io/%s/json?token=%s", ip, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ipInfo IPInfo
	if err := json.NewDecoder(resp.Body).Decode(&ipInfo); err != nil {
		return nil, err
	}

	return &ipInfo, nil
}

func getUserAgentInfo(userAgent, apiKey string) (*UserAgentInfo, error) {
	url := "https://api.whatismybrowser.com/api/v2/user_agent_parse"
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	payload := map[string]string{"user_agent": userAgent}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	payloadBuffer := bytes.NewBuffer(payloadBytes)

	req, err := http.NewRequest("POST", url, payloadBuffer)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("User Agent Response JSON:", string(body))

	var userAgentInfo UserAgentInfo
	if err := json.Unmarshal(body, &userAgentInfo); err != nil {
		return nil, err
	}

	return &userAgentInfo, nil
}

func main() {
	ipApiKey := "IP_API_KEY"
	userAgentApiKey := "USER_AGENT_API_KEY"
	ip := "8.8.8.8"

	if ipApiKey == "" || userAgentApiKey == "" {
		fmt.Println("API key is missing. Please set your API keys.")
		os.Exit(1)
	}

	ipInfo, err := getIPInfo(ip, ipApiKey)
	if err != nil {
		fmt.Println("Error getting IP info:", err)
		return
	}

	userAgentString := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9"
	userAgentInfo, err := getUserAgentInfo(userAgentString, userAgentApiKey)
	if err != nil {
		fmt.Println("Error getting User Agent info:", err)
		return
	}

	fmt.Printf("IP Address: %s\n", ipInfo.IP)
	fmt.Printf("Hostname: %s\n", ipInfo.Hostname)
	fmt.Printf("City: %s\n", ipInfo.City)
	fmt.Printf("Region: %s\n", ipInfo.Region)
	fmt.Printf("Country: %s\n", ipInfo.Country)
	fmt.Printf("Location: %s\n", ipInfo.Loc)
	fmt.Printf("Organization: %s\n", ipInfo.Org)
	fmt.Printf("Postal Code: %s\n", ipInfo.Postal)
	fmt.Printf("Timezone: %s\n", ipInfo.Timezone)
	fmt.Printf("User Agent String: %s\n", userAgentInfo.Parse.UserAgent)
	fmt.Printf("Software Name: %s\n", userAgentInfo.Parse.SoftwareName)
	fmt.Printf("Software Version: %s\n", userAgentInfo.Parse.SoftwareVersion)
	fmt.Printf("Operating System Name: %s\n", userAgentInfo.Parse.OperatingSystemName)
	fmt.Printf("Operating System Version: %s\n", userAgentInfo.Parse.OperatingSystemVersion)
}
