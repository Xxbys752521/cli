// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/larksuite/cli/internal/core"
)

// AuthCodeResult is the result of the Authorization Code flow.
type AuthCodeResult struct {
	OK      bool
	Token   *DeviceFlowTokenData
	Error   string
	Message string
}

// AuthCodeServer holds the local HTTP server state for the OAuth callback.
type AuthCodeServer struct {
	listener    net.Listener
	state       string
	codeCh      chan string
	errCh       chan error
	redirectURI string
}

// NewAuthCodeServer creates and starts a local HTTP server for OAuth callback.
// It listens on localhost:8080 (matching the official Go example).
func NewAuthCodeServer() (*AuthCodeServer, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		return nil, fmt.Errorf("failed to start local server on port 8080 (is it in use?): %v", err)
	}

	stateBytes := make([]byte, 16)
	if _, err := rand.Read(stateBytes); err != nil {
		listener.Close()
		return nil, fmt.Errorf("failed to generate state: %v", err)
	}

	s := &AuthCodeServer{
		listener:    listener,
		state:       hex.EncodeToString(stateBytes),
		codeCh:      make(chan string, 1),
		errCh:       make(chan error, 1),
		redirectURI: "http://localhost:8080/callback",
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", s.handleCallback)
	server := &http.Server{Handler: mux}
	go server.Serve(listener)

	return s, nil
}

// RedirectURI returns the redirect URI for the OAuth authorize request.
func (s *AuthCodeServer) RedirectURI() string {
	return s.redirectURI
}

// State returns the CSRF state parameter.
func (s *AuthCodeServer) State() string {
	return s.state
}

// Close shuts down the local server.
func (s *AuthCodeServer) Close() {
	s.listener.Close()
}

// WaitForCode blocks until the authorization code is received or context is cancelled.
func (s *AuthCodeServer) WaitForCode(ctx context.Context) (string, error) {
	select {
	case code := <-s.codeCh:
		return code, nil
	case err := <-s.errCh:
		return "", err
	case <-ctx.Done():
		return "", fmt.Errorf("authorization timed out")
	}
}

func (s *AuthCodeServer) handleCallback(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	state := query.Get("state")
	code := query.Get("code")

	if state != s.state {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, htmlPage("Authorization Failed", "State mismatch. Please try again."))
		s.errCh <- fmt.Errorf("state mismatch: possible CSRF attack")
		return
	}

	if code == "" {
		errMsg := query.Get("error")
		if errMsg == "" {
			errMsg = "no authorization code received"
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, htmlPage("Authorization Failed", errMsg))
		s.errCh <- fmt.Errorf("authorization failed: %s", errMsg)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, htmlPage("Authorization Successful", "You can close this page and return to the terminal."))
	s.codeCh <- code
}

func htmlPage(title, message string) string {
	return fmt.Sprintf(`<!DOCTYPE html><html><head><meta charset="utf-8"><title>%s</title>
<style>body{font-family:system-ui,sans-serif;display:flex;justify-content:center;align-items:center;height:100vh;margin:0;background:#f5f5f5}
.card{background:#fff;padding:40px;border-radius:12px;box-shadow:0 2px 8px rgba(0,0,0,.1);text-align:center;max-width:400px}
h1{margin:0 0 16px;font-size:24px;color:#333}p{color:#666;font-size:16px}</style>
</head><body><div class="card"><h1>%s</h1><p>%s</p></div></body></html>`, title, title, message)
}

// BuildAuthorizeURL constructs the Feishu/Lark OAuth authorize URL.
// Uses the accounts domain endpoint: /open-apis/authen/v1/authorize
// with client_id per official docs.
// Note: URL is constructed manually to avoid url.Values encoding colons in
// scope values (e.g. "im:message:readonly" must NOT become "im%3Amessage%3Areadonly").
func BuildAuthorizeURL(brand core.LarkBrand, appId, redirectURI, state, scope string) string {
	ep := core.ResolveEndpoints(brand)
	q := "client_id=" + appId +
		"&redirect_uri=" + url.QueryEscape(redirectURI)
	if scope != "" {
		q += "&scope=" + strings.ReplaceAll(scope, " ", "%20")
	}
	if state != "" {
		q += "&state=" + state
	}
	return ep.Accounts + "/open-apis/authen/v1/authorize?" + q
}

// ExchangeAuthCode exchanges an authorization code for an access token.
func ExchangeAuthCode(httpClient *http.Client, appId, appSecret string, brand core.LarkBrand, code, redirectURI string) (*AuthCodeResult, error) {
	ep := core.ResolveEndpoints(brand)
	tokenURL := ep.Open + "/open-apis/authen/v2/oauth/token"

	reqBody := map[string]string{
		"grant_type":    "authorization_code",
		"code":          code,
		"client_id":     appId,
		"client_secret": appSecret,
		"redirect_uri":  redirectURI,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal token request: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(string(bodyBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token exchange failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("token exchange failed: read body: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("token exchange failed: invalid JSON response")
	}

	if errStr := getStr(data, "error"); errStr != "" {
		desc := getStr(data, "error_description")
		if desc == "" {
			desc = errStr
		}
		return &AuthCodeResult{OK: false, Error: errStr, Message: desc}, nil
	}

	accessToken := getStr(data, "access_token")
	if accessToken == "" {
		return nil, fmt.Errorf("token exchange failed: no access_token in response")
	}

	refreshToken := getStr(data, "refresh_token")
	tokenExpiresIn := getInt(data, "expires_in", 7200)
	refreshExpiresIn := getInt(data, "refresh_token_expires_in", 604800)
	if refreshToken == "" {
		refreshExpiresIn = tokenExpiresIn
	}

	return &AuthCodeResult{
		OK: true,
		Token: &DeviceFlowTokenData{
			AccessToken:      accessToken,
			RefreshToken:     refreshToken,
			ExpiresIn:        tokenExpiresIn,
			RefreshExpiresIn: refreshExpiresIn,
			Scope:            getStr(data, "scope"),
		},
	}, nil
}
