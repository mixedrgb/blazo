package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	acctID      string
	apiURL      string
	authToken   string
	creds       []string
	endpoint    string
	err         error
	fileName    string
	marshalJSON []byte
	resp        *http.Response
)

// Boilerplate[tm]
func ec(e error) {
	if e != nil {
		log.Fatalf("E: %s\n", e)
	}
}

func fetchCredentials() {
	if 1 != len(os.Args) {
		return
	} else {
		fileName = "./key.key"
		log.Printf("I: No file provided. Using %s\n", fileName)
		// Read in the file name
		tmpFile, err := ioutil.ReadFile(fileName)
		// Trim, apparently, a newline introduced by the func -.-
		tmpTrim := strings.TrimSpace(string(tmpFile))
		// Split KV into the creds slice
		creds := strings.Split(tmpTrim, ":")
		// Assign API key and auth token
		acctID, authToken = creds[0], creds[1]
		// Boilerplate[tm]
		ec(err)
		// Feeling logged, might delete later
		log.Printf("I: Using token %s:%s\n", acctID, authToken)
	}
}

func formatJson() {
	// make a map slice for the JSON (key/value)
	contentJSON := map[string]string{"accountId": acctID}
	// "convert" to JSON-readable format
	marshalJSON, err = json.Marshal(contentJSON)
	// Boilerplate[tm]
	ec(err)
}

func formatRequest() {
	//content type is json
	contentType := "application/json"
	// formulate (doesn't execute) a new POST request
	req, err := http.NewRequest("POST", apiURL+endpoint,
		bytes.NewBuffer(marshalJSON))
	// Boilerplate[tm]
	ec(err)
	// add and set appropriate headers
	req.Header.Set("Content-Type", contentType)
	req.Header.Add("Authorization", authToken)
	req.Header.Set("User-Agent", "beep boop")

	sendRequest(req)
}

func sendRequest(req *http.Request) {
	// create the client
	client := &http.Client{}
	// do the thing
	resp, err := client.Do(req)
	// Boilerplate[tm]
	ec(err)
	if resp.Status == string(401) {
		log.Printf("Token expired, getting another one")
		writeResponse(resp)
	}
	writeResponse(resp)
}

func writeResponse(resp *http.Response) {
	// read all data from byte slice (io.Reader) until EOF/err
	respBody, err := ioutil.ReadAll(resp.Body)
	// Boilerplate[tm]
	ec(err)
	// print out as a string, not weird bs
	fmt.Println(string(respBody))
}

func main() {
	apiURL = "https://api002.backblazeb2.com/b2api/v2/"
	endpoint = "b2_list_keys"
	fetchCredentials()
	formatJson()
	formatRequest()
}
