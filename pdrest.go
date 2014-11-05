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

func main() {

	// Configure flags and command line defaults
	hostnamePtr := flag.String("hostname", "localhost", "Appliance Hostname")

	appUsernamePtr := flag.String("applianceAdmin", "admin", "Appliance Admin User Name")

	samUsernamePtr := flag.String("runtimeAdmin", "sec_master", "Runtime Admin User Name")

	filePtr := flag.String("commands", "commands.txt", "Commands File Path")

	flag.Parse()

	// Get the passwords from standard input; alternatively these can be feed
	// in using standard input redirection like "< credentials.txt"

	// Use Scanln instead of Scanf to deal with Windows line endings:
	// https://code.google.com/p/go/issues/detail?id=5391

	var appPassword string
	fmt.Println("Appliance Admin User Password:")
	fmt.Scanln(&appPassword)

	var samPassword string
	fmt.Println("Runtime Admin User Password:")
	fmt.Scanln(&samPassword)

	// Configure the URL based on SAM defaults
	url := "https://" + *hostnamePtr + "/isam/pdadmin"

	// Read in the commands from the passed filepath
	commands, err := readCommands(*filePtr)
	if err != nil {
		log.Fatal(err)
	}

	// Create a message
	m := &message{AdminID: *samUsernamePtr, AdminPwd: samPassword, Commands: commands}

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
	req.SetBasicAuth(*appUsernamePtr, appPassword)

	// Create a transport that ignores any untrusted SSL certificates on the appliance
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

func readCommands(fname string) (commands []string, err error) {

	b, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(b), "\n")

	commands = make([]string, 0, len(lines))

	for _, l := range lines {

		// Trim to remove tabs (\t), Windows carriage returns (\r), etc.
		command := strings.TrimSpace(l)

		// Ignore empty lines or lines starting with a "#"
		if len(command) == 0 || strings.HasPrefix(command, "#") {
			continue
		}

		commands = append(commands, command)
	}

	return commands, nil

}
