package licensing

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	CLIENT_VERSION = "v1"
)

var responseMap = map[string]interface{}{
	"client":   &ClientResponse{},
	"validate": &ValidateResponse{},
}

func invokeGetAPI(path string, opts map[string]string, suppressExp ...bool) interface{} {
	URL := constructURL(opts["licenseId"], path)
	// log.Println("URL: " + URL)

	res, err := http.Get(URL)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	resData := responseMap[path]
	// fmt.Println("response is---", PrettyPrint(resData))
	// fmt.Println("Response body is ", string(body))
	if err := json.Unmarshal(body, &resData); err != nil { // Parse []byte to go struct pointer
		log.Println("Can not unmarshal JSON")
		log.Fatal(err)
	}

	var statusCode int
	var message string
	switch v := resData.(type) {
	case *ValidateResponse:
		statusCode = v.StatusCode
		message = v.Message
	case *ClientResponse:
		statusCode = v.StatusCode
		message = v.Message
	default:
		statusCode = 400
		message = "Unknown error occured."
	}

	// fmt.Println("After this status code: ", statusCode)
	// fmt.Println("Message: ", message)
	var suppress bool
	if len(suppressExp) > 0 {
		suppress = suppressExp[0]
	}

	if !suppress && statusCode != 200 {
		log.Fatal(message)
	}

	return resData
}

func constructURL(licenseID string, path string) string {
	config := GetConfig()
	params := url.Values{}

	URL := config.LicenseServerURL

	params.Add("licenseId", licenseID)
	params.Add("entitlementId", config.EntitlementID)

	key, check := os.LookupEnv("CHEF_LICENSE_SERVER")
	if check {
		URL = key
	}

	URL += "/" + CLIENT_VERSION + "/" + path + "?" + params.Encode()

	return URL
}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
