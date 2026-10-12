
package googledrive

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const driveScope = "https://www.googleapis.com/auth/drive"

type storedToken struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	Expiry      time.Time `json:"expiry"`
}

type oauthResult struct {
	token *oauth2.Token
	err   error
}

func tokenPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "rioctl", "drive-token.json"), nil
}

func loadToken() (*oauth2.Token, error) {
	path, err := tokenPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var saved storedToken
	if err := json.Unmarshal(data, &saved); err != nil {
		return nil, fmt.Errorf("decode cached access token: %w", err)
	}

	token := &oauth2.Token{
		AccessToken: saved.AccessToken,
		TokenType:   saved.TokenType,
		Expiry:      saved.Expiry,
	}

	// Require an expiry; never treat an incomplete cached token as valid.
	if token.AccessToken == "" || token.Expiry.IsZero() || !token.Valid() {
		return nil, errors.New("cached access token is missing or expired")
	}

	return token, nil
}

func saveToken(token *oauth2.Token) error {
	path, err := tokenPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	// Deliberately persist only the access token. Do not save RefreshToken.
	saved := storedToken{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		Expiry:      token.Expiry,
	}

	data, err := json.MarshalIndent(saved, "", "  ")
	if err != nil {
		return err
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}

	return nil
}

func deleteToken() error {
	path, err := tokenPath()
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Start()
}

// Authorize returns an HTTP client authenticated with Google Drive.
//
// A valid cached access token is reused. If it is absent or expired,
// the user is prompted to authorize in a browser.
//
// reauth forces a new browser authorization.
func Authorize(
	ctx context.Context,
	reauth bool,
) (*http.Client, error) {
	if !reauth {
		token, err := loadToken()
		if err == nil {
			return oauth2.NewClient(
				ctx,
				oauth2.StaticTokenSource(token),
			), nil
		}
	}

	if err := deleteToken(); err != nil {
		return nil, fmt.Errorf("remove cached access token: %w", err)
	}

	credentials, err := Credentials()
	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(credentials, driveScope)
	if err != nil {
		return nil, fmt.Errorf("parse embedded OAuth credentials: %w", err)
	}

	token, err := authorizeInBrowser(ctx, config)
	if err != nil {
		return nil, err
	}

	// Do not retain an offline refresh token.
	token.RefreshToken = ""

	if token.AccessToken == "" || token.Expiry.IsZero() {
		return nil, errors.New("Google returned an access token without an expiry")
	}

	if err := saveToken(token); err != nil {
		return nil, fmt.Errorf("save access token: %w", err)
	}

	return oauth2.NewClient(
		ctx,
		oauth2.StaticTokenSource(token),
	), nil
}

func authorizeInBrowser(
	ctx context.Context,
	config *oauth2.Config,
) (*oauth2.Token, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		return nil, fmt.Errorf("start OAuth callback listener: %w", err)
	}
	defer listener.Close()

	state, err := randomState()
	if err != nil {
		return nil, err
	}

	config.RedirectURL = fmt.Sprintf(
		"http://%s/oauth2/callback",
		listener.Addr().String(),
	)

	results := make(chan oauthResult, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/oauth2/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.URL.Query().Get("state") != state {
			http.Error(w, "Invalid OAuth state", http.StatusBadRequest)
			return
		}

		if oauthErr := r.URL.Query().Get("error"); oauthErr != "" {
			http.Error(w, "Authorization was not completed. You can close this window.", http.StatusBadRequest)
			select {
			case results <- oauthResult{err: fmt.Errorf("Google authorization failed: %s", oauthErr)}:
			default:
			}
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Missing authorization code", http.StatusBadRequest)
			return
		}

		token, exchangeErr := config.Exchange(r.Context(), code)
		if exchangeErr != nil {
			http.Error(w, "Failed to exchange authorization code. You can close this window.", http.StatusBadGateway)
			select {
			case results <- oauthResult{err: fmt.Errorf("exchange OAuth code: %w", exchangeErr)}:
			default:
			}
			return
		}

		fmt.Fprintln(w, "Authorization successful. You can close this window and return to rioctl.")

		select {
		case results <- oauthResult{token: token}:
		default:
		}
	})

	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	serverDone := make(chan struct{})
	go func() {
		defer close(serverDone)
		_ = server.Serve(listener)
	}()

	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
		<-serverDone
	}()

	authURL := config.AuthCodeURL(state)

	fmt.Println("Opening your browser for Google Drive authorization...")
	if err := openBrowser(authURL); err != nil {
		fmt.Printf("Could not open browser automatically. Open this URL:\n%s\n", authURL)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-results:
		if result.err != nil {
			return nil, result.err
		}
		if result.token == nil {
			return nil, errors.New("Google authorization returned no token")
		}
		return result.token, nil
	}
}

func randomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate OAuth state: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
