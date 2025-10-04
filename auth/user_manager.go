package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type BasicUser struct {
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
}

type UserManager struct {
	users      map[string]*BasicUser
	realm      string
	mu         sync.RWMutex
	userDBFile string
}

func NewBasicAuthManager(dbFile string) (*UserManager, error) {
	auth := &UserManager{
		users:      make(map[string]*BasicUser),
		realm:      "Home Automation System",
		userDBFile: dbFile,
	}

	if err := auth.loadUsers(); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load user database: %w", err)
		}
	}
	return auth, nil
}

func (ba *UserManager) AddUser(username, password string) error {
	ba.mu.Lock()
	defer ba.mu.Unlock()

	if _, exists := ba.users[username]; exists {
		return fmt.Errorf("user '%s' already exists", username)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	ba.users[username] = &BasicUser{
		Username:     username,
		PasswordHash: string(hash),
	}

	return ba.saveUsers()
}

func (ba *UserManager) ValidateCredentials(username, password string) (*BasicUser, bool) {
	ba.mu.RLock()
	user, exists := ba.users[username]
	ba.mu.RUnlock()

	if !exists {
		return nil, false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, false // Passwords don't match
	}

	return user, true
}

func (ba *UserManager) saveUsers() error {
	data, err := json.MarshalIndent(ba.users, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal users to JSON: %w", err)
	}
	return os.WriteFile(ba.userDBFile, data, 0600)
}

func (ba *UserManager) loadUsers() error {
	data, err := os.ReadFile(ba.userDBFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &ba.users)
}
