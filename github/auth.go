package github

import (
	"os"
	"path/filepath"

	"github.com/99designs/keyring"
)

const (
	serviceName = "tgr"
	userKey     = "github_username"
	tokenKey    = "github_token"
)

// AuthService handles credential storage
type AuthService struct {
	ring keyring.Keyring
}

// NewAuthService creates a new authentication service
func NewAuthService() (*AuthService, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		home, _ := os.UserHomeDir()
		configDir = filepath.Join(home, ".config")
	}

	keyringDir := filepath.Join(configDir, "tgr", "keyring")
	if err := os.MkdirAll(keyringDir, 0700); err != nil {
		return nil, err
	}

	ring, err := keyring.Open(keyring.Config{
		ServiceName:      serviceName,
		FileDir:          keyringDir,
		FilePasswordFunc: keyring.TerminalPrompt,
	})
	if err != nil {
		return nil, err
	}
	return &AuthService{ring: ring}, nil
}

// GetCredentials retrieves the stored username and token
func (s *AuthService) GetCredentials() (string, string, error) {
	userItem, err := s.ring.Get(userKey)
	if err != nil && err != keyring.ErrKeyNotFound {
		return "", "", err
	}

	tokenItem, err := s.ring.Get(tokenKey)
	if err != nil && err != keyring.ErrKeyNotFound {
		return "", "", err
	}

	return string(userItem.Data), string(tokenItem.Data), nil
}

// SaveCredentials stores the username and token
func (s *AuthService) SaveCredentials(username, token string) error {
	err := s.ring.Set(keyring.Item{
		Key:  userKey,
		Data: []byte(username),
	})
	if err != nil {
		return err
	}

	return s.ring.Set(keyring.Item{
		Key:  tokenKey,
		Data: []byte(token),
	})
}
