package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Bowery/prompt"
)

// Struct for the json request
type message struct {
	AdminID  string   `json:"admin_id"`
	AdminPwd string   `json:"admin_pwd"`
	Commands []string `json:"commands"`
}

// Struct for the json response
type response struct {
	Result string `json:"result"`
}

var (
	appPassword string
	samPassword string
	err         error
)

func main() {

	// Configure flags and command line defaults
	hostname := flag.String("host", "localhost", "Appliance Hostname")

	appUsername := flag.String("applianceAdmin", "admin",
		"Appliance Admin User Name")

	samUsername := flag.String("runtimeAdmin", "sec_master",
		"Runtime Admin User Name")

	commandsFile := flag.String("commands", "commands.txt", "Commands File Path")

	credsFile := flag.String("creds", "", "Credentials File Path")

	flag.Parse()

	if *credsFile == "" {

		// Since no credentials file was passed, prompt the user for passwords
		appPassword, err = prompt.Password("Appliance Admin User Password")
		if err != nil {
			log.Fatal(err)
		}

		samPassword, err = prompt.Password("Runtime Admin User Password")
		if err != nil {
			log.Fatal(err)
		}

	} else {

		// Since a credentials file was passed, parse the file for passwords
		creds, err := readLines(*credsFile)
		if err != nil {
			log.Fatal(err)
		}

		if len(creds) > 1 {

			appPassword = creds[0]

			samPassword = creds[1]

		} else {
			log.Fatal("Credential file needs two passwords, each on their own line.")
		}

	}

	// Configure the URL based on SAM defaults
	url := "https://" + *hostname + "/isam/pdadmin"

	// Read in the commands from the passed filepath
	commands, err := readLines(*commandsFile)
	if err != nil {
		log.Fatal(err)
	}

	// Create a message
	m := &message{AdminID: *samUsername, AdminPwd: samPassword,
		Commands: commands}

	// Encode the message struct as json
	b, _ := json.Marshal(m)

	// Create an HTTP POST request with the json message
	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}

	// Add headers to the requset for json and HTTP Basic Authentication
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(*appUsername, appPassword)

	// Create a transport that ignores any untrusted SSL certificates
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Create an HTTP client with the transport
	c := &http.Client{Transport: tr}

	// Make the POST request and get the response
	res, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// Print the HTTP response code
	fmt.Println("Status Code: ", res.StatusCode)

	var result response

	// Read and unmarshal the json response
	data, err := ioutil.ReadAll(res.Body)
	if err == nil && data != nil {
		err = json.Unmarshal(data, &result)
	} else {
		log.Fatal(err)
	}

	// Print the output of each submitted command
	fmt.Printf("Result: \n%s\n", result.Result)

}

func readLines(filePath string) (lines []string, err error) {

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	rawLines := strings.Split(string(b), "\n")

	lines = make([]string, 0, len(rawLines))

	for _, l := range rawLines {

		// Trim to remove tabs (\t), Windows carriage returns (\r), etc.
		line := strings.TrimSpace(l)

		// Ignore empty lines or lines starting with a "#"
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		lines = append(lines, line)
	}

	return lines, nil

}
