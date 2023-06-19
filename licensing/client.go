package licensing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type Response struct {
	Data       interface{} `json:"data"`
	Message    string      `json:"message"`
	StatusCode int         `json:"status_code"`
}

// ENDPOINTS CONSTANT
const (
	CLIENT = "v1/client"
)

func invokeGetAPI(opts map[string]string, URL string) {
	params := url.Values{}

	params.Add("licenseId", opts["licenseId"])
	params.Add("entitlementId", opts["entitlementId"])

	// fmt.Println("params are :", params)
	key, check := os.LookupEnv("CHEF_LICENSE_SERVER")
	if check {
		URL = key
	}
	res, err := http.Get(URL + "/" + CLIENT + "?" + params.Encode())
	if err != nil {
		fmt.Print(err.Error())
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
		fmt.Println("err is", err)
	}
	// fmt.Println("response is---", PrettyPrint(response))
	if response.Data == false {
		fmt.Println("Error:", response.Message)
		os.Exit(1)
	}
}

// func PrettyPrint(i interface{}) string {
// 	s, _ := json.MarshalIndent(i, "", "\t")
// 	return string(s)
// }
