package tm

import(
	"fmt"
	"testing"
	"io/ioutil"
	"encoding/json"
	"bytes"
)

const Client_contr = "Client_Controller"

type TMIni struct {
	Token string `json:"token"`
	ChatID string `json:"chat_id"`
}

func TestSend(t *testing.T) {
	//read ini file
	file, err := ioutil.ReadFile("tm.json")
	if err != nil {
		t.Fatalf("%v", err)
	}
	tm_ini := TMIni{}
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
	err = json.Unmarshal([]byte(file), &tm_ini)		
	if err != nil {
		t.Fatalf("%v", err)
	}
	
	
	parameters := map[string]string{
		 "text": "Hellow, world!",
		  "chat_id": tm_ini.ChatID,
	}
	
	resp, err := ApiRequestJson(tm_ini.Token, "sendMessage", parameters)
	if err != nil {
		t.Fatalf("%v", err)
	}
	fmt.Println(string(resp))
}

