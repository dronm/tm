package tm

import(
	"fmt"
	"testing"
	"os"
)

type TMIni struct {
	Token string
	ChatID string
}

func TestSend(t *testing.T) {
	tm_ini := TMIni{Token:os.Getenv("TM_TEST_BOT_TOKEN"), ChatID: os.Getenv("TM_TEST_CHAT_ID")}
	if tm_ini.Token == "" {
		t.Fatalf("environment variable 'TM_TEST_BOT_TOKEN' is not initialized")
	}
	if tm_ini.ChatID == "" {
		t.Fatalf("environment variable 'TM_TEST_CHAT_ID' is not initialized")
	}
	parameters := map[string]string{
		 "text": "Hello, world!",
		  "chat_id": tm_ini.ChatID,
	}
	
	resp, err := ApiRequestJson(tm_ini.Token, "sendMessage", parameters)
	if err != nil {
		t.Fatalf("%v", err)
	}
	fmt.Println(string(resp))
}

