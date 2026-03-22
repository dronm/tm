// Package tm provides a method to send messages to telegram bot.
package tm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/proxy"
)

const (
	defMaxIdleConns        = 100
	defIdleConnTimeout     = 90 * time.Second
	defTLSHandshakeTimeout = 10 * time.Second
)

const apiURLTmpl = "https://api.telegram.org/bot%s/%s"

type HTTPTransportConfig struct {
	MaxIdleConns        int
	IdleConnTimeout     time.Duration
	TLSHandshakeTimeout time.Duration
}

type ProxyConfig struct {
	Address  string
	Username string
	Password string
}

func RequestJSON(
	botToken string,
	method string,
	parameters map[string]string,
	proxyCfg *ProxyConfig,
	httpTransportCfg *HTTPTransportConfig,
) ([]byte, error) {
	jsonParams, err := json.Marshal(parameters)
	if err != nil {
		return nil, err
	}

	u := fmt.Sprintf(apiURLTmpl, botToken, method)

	transpCfg := HTTPTransportConfig{
		MaxIdleConns:        defMaxIdleConns,
		IdleConnTimeout:     defIdleConnTimeout,
		TLSHandshakeTimeout: defTLSHandshakeTimeout,
	}
	if httpTransportCfg != nil {
		if httpTransportCfg.MaxIdleConns > 0 {
			transpCfg.MaxIdleConns = httpTransportCfg.MaxIdleConns
		}
		if httpTransportCfg.IdleConnTimeout > 0 {
			transpCfg.IdleConnTimeout = httpTransportCfg.IdleConnTimeout
		}
		if httpTransportCfg.TLSHandshakeTimeout > 0 {
			transpCfg.TLSHandshakeTimeout = httpTransportCfg.TLSHandshakeTimeout
		}
	}
	client, err := newHTTPClient(transpCfg, proxyCfg)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, u, bytes.NewBuffer(jsonParams))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Close = true

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	msg := struct {
		Ok          bool   `json:"ok"`
		ErrorCode   int    `json:"error_code"`
		Description string `json:"description"`
	}{}

	if err := json.Unmarshal(body, &msg); err != nil {
		return nil, err
	}

	if resp.StatusCode >= 500 {
		return nil, fmt.Errorf("http error: %d", resp.StatusCode)
	}

	if resp.StatusCode == http.StatusBadRequest {
		return nil, errors.New("HTTP error: 400, Invalid access token provided")
	}

	if resp.StatusCode != http.StatusOK {
		return body, fmt.Errorf("http error: %d, description: %s", resp.StatusCode, msg.Description)
	}

	return body, nil
}

func newHTTPClient(transportConfig HTTPTransportConfig, proxyConfig *ProxyConfig) (*http.Client, error) {
	transport := &http.Transport{
		ForceAttemptHTTP2:   true,
		MaxIdleConns:        transportConfig.MaxIdleConns,
		IdleConnTimeout:     transportConfig.IdleConnTimeout,
		TLSHandshakeTimeout: transportConfig.TLSHandshakeTimeout,
	}
	fmt.Println("transpCfg:", transportConfig)

	if proxyConfig != nil && proxyConfig.Address != "" {
		var auth *proxy.Auth
		if proxyConfig.Username != "" {
			auth = &proxy.Auth{
				User:     proxyConfig.Username,
				Password: proxyConfig.Password,
			}
		}

		dialer, err := proxy.SOCKS5("tcp", proxyConfig.Address, auth, proxy.Direct)
		if err != nil {
			return nil, err
		}

		contextDialer, ok := dialer.(proxy.ContextDialer)
		if !ok {
			return nil, errors.New("proxy dialer does not support context dialing")
		}

		transport.DialContext = func(
			ctx context.Context,
			network string,
			addr string,
		) (net.Conn, error) {
			return contextDialer.DialContext(ctx, network, addr)
		}
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}, nil
}
