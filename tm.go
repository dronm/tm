package tm

import(
	"fmt"
	"errors"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"bytes"	
)

const API_URL_TMPL = "https://api.telegram.org/bot%s/"

func ApiRequestJson(botToken string, method string, parameters map[string]string) ([]byte,error) {

	parameters["method"] = method
	json_params, err := json.Marshal(parameters)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf(API_URL_TMPL, botToken)
	
	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, u, bytes.NewBuffer(json_params))
	r.Header.Add("Content-Type", "application/json")
	r.Close = true
	resp, err := client.Do(r)	
	if err != nil {				
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 500 {
		return nil, errors.New(fmt.Sprintf("HTTP error: %d", resp.StatusCode));
	
	} else if resp.StatusCode == 400 {
		return nil, errors.New("HTTP error: 400, Invalid access token provided");
		
	}

	//200 or error need body
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	msg := struct {
		Ok bool `json:"ok"`
		Error_code int `json:"error_code"`
		Description string `json:"description"`
	}{}

	if err := json.Unmarshal(body, &msg); err != nil {
		return nil, err
	}			
	
	if resp.StatusCode != 200 {
		//error
		return body, errors.New(fmt.Sprintf("HTTP error: %d, description: %s", resp.StatusCode, msg.Description))
		
	}
	
	//OK 200
	return body, nil
}
