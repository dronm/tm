package tm

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	_ = godotenv.Load(".env")
	os.Exit(m.Run())
}

func TestRequestJSON(t *testing.T) {
	botToken := os.Getenv("BOT_TOKEN")
	chatID := os.Getenv("CHAT_ID")
	proxyAddr := os.Getenv("PROXY_ADDR")
	proxyUser := os.Getenv("PROXY_USER")
	proxyPassword := os.Getenv("PROXY_PASSWORD")

	if botToken == "" {
		t.Fatal("BOT_TOKEN is not set")
	}
	if chatID == "" {
		t.Fatal("CHAT_ID is not set")
	}

	var proxyConfig *ProxyConfig
	if proxyAddr != "" {
		proxyConfig = &ProxyConfig{
			Address:  proxyAddr,
			Username: proxyUser,
			Password: proxyPassword,
		}
	}

	body, err := RequestJSON(
		botToken,
		"sendMessage",
		map[string]string{
			"chat_id": chatID,
			"text":    "test message from go test",
		},
		proxyConfig,
	)
	if err != nil {
		t.Fatalf("RequestJSON() error = %v", err)
	}

	if len(body) == 0 {
		t.Fatal("expected non-empty response body")
	}
}
