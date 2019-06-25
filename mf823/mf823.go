package mf823

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

//Status retrieves the status from a ZTE MF823 4G Modem. The individual
//parameters of the status have to be supplied.
func Status(address string, params ...string) (map[string]interface{}, error) {

	paramsList := ""

	if len(params) > 0 {
		paramsList = params[0]
	}

	// concatenate the parameters to a comma separated value string
	for i := 1; i < len(params); i++ {
		paramsList = fmt.Sprintf("%s,%s", paramsList, params[i])
	}

	// url parameters. They have been reverse engineered by inspecting the
	// ajax calls from the webui
	urlParams := url.Values{}
	urlParams.Add("isTest", "false")
	urlParams.Add("multi_data", "1")
	urlParams.Add("cmd", paramsList)
	urlParams.Add("_", fmt.Sprintf("%v", time.Now().UnixNano()/1000000))

	// assemble the url to be queried
	_url := "http://" + address + "/goform/goform_get_cmd_process?" + urlParams.Encode()

	client := http.Client{}
	req, err := http.NewRequest(
		http.MethodGet,
		_url,
		nil,
	)

	if err != nil {
		return nil, err
	}

	// Referer header is mandatory. Otherwise no data will be returned.
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Referer", "http://"+address+"/status.html")

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(resp.Status)
	}

	if resp.Body == nil {
		return nil, fmt.Errorf("received empty response")
	}

	m := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return nil, err
	}

	return m, nil
}
