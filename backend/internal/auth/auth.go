package auth

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Role string

const (
	RoleAdmin      Role = "admin"
	RoleTeacher    Role = "teacher"
	RoleStudent    Role = "student"
	RoleEnterprise Role = "enterprise"
	RoleProbe      Role = "probe"
)

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Role      Role      `json:"role"`
	ClassroomID string  `json:"classroom_id"`
	PasswordHash string `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	LastLogin  time.Time `json:"last_login"`
	IsActive   bool      `json:"is_active"`
}

type Session struct {
	ID        string
	UserID    string
	Role      Role
	ClassroomID string
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
	IPAddress string
}

type Permission struct {
	Resource string
	Action   string
}

var RolePermissions = map[Role][]Permission{
	RoleAdmin: {
		{Resource: "*", Action: "*"},
	},
	RoleTeacher: {
		{Resource: "classroom", Action: "read"},
		{Resource: "classroom", Action: "write"},
		{Resource: "student", Action: "read"},
		{Resource: "error_cluster", Action: "read"},
		{Resource: "teaching_alert", Action: "read"},
		{Resource: "vm", Action: "read"},
	},
	RoleStudent: {
		{Resource: "workspace", Action: "read"},
		{Resource: "workspace", Action: "write"},
		{Resource: "ai_companion", Action: "read"},
		{Resource: "ai_companion", Action: "write"},
		{Resource: "vm", Action: "read"},
		{Resource: "vm", Action: "write"},
	},
	RoleEnterprise: {
		{Resource: "sandbox", Action: "read"},
		{Resource: "sandbox", Action: "write"},
		{Resource: "audit_log", Action: "read"},
		{Resource: "vm", Action: "read"},
	},
	RoleProbe: {
		{Resource: "probe_event", Action: "write"},
		{Resource: "heartbeat", Action: "write"},
	},
}

type AuthService struct {
	users     map[string]*User
	sessions  map[string]*Session
	apiKeys   map[string]*User
	mu        sync.RWMutex
}

func NewAuthService() *AuthService {
	return &AuthService{
		users:    make(map[string]*User),
		sessions: make(map[string]*Session),
		apiKeys:  make(map[string]*User),
	}
}

func (a *AuthService) CreateUser(id, username, password string, role Role, classroomID string) *User {
	user := &User{
		ID:           id,
		Username:     username,
		Role:         role,
		ClassroomID:  classroomID,
		PasswordHash: hashPassword(password),
		CreatedAt:    time.Now(),
		IsActive:     true,
	}

	a.mu.Lock()
	a.users[id] = user
	a.mu.Unlock()

	return user
}

func (a *AuthService) Authenticate(username, password string) (*User, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	for _, user := range a.users {
		if user.Username == username && user.IsActive {
			if subtle.ConstantTimeCompare([]byte(user.PasswordHash), []byte(hashPassword(password))) == 1 {
				user.LastLogin = time.Now()
				return user, nil
			}
		}
	}

	return nil, errors.New("invalid username or password")
}

func (a *AuthService) CreateSession(userID, ipAddress string) (*Session, error) {
	a.mu.RLock()
	user, exists := a.users[userID]
	a.mu.RUnlock()

	if !exists {
		return nil, errors.New("user not found")
	}

	sessionID := generateSessionID()
	token := generateToken(userID, sessionID, user.Role)

	session := &Session{
		ID:          sessionID,
		UserID:      userID,
		Role:        user.Role,
		ClassroomID: user.ClassroomID,
		Token:       token,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		IPAddress:   ipAddress,
	}

	a.mu.Lock()
	a.sessions[sessionID] = session
	a.mu.Unlock()

	return session, nil
}

func (a *AuthService) ValidateToken(token string) (*Session, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	for _, session := range a.sessions {
		if session.Token == token {
			if time.Now().After(session.ExpiresAt) {
				return nil, errors.New("session expired")
			}
			return session, nil
		}
	}

	return nil, errors.New("invalid token")
}

func (a *AuthService) CheckPermission(role Role, resource, action string) bool {
	permissions, exists := RolePermissions[role]
	if !exists {
		return false
	}

	for _, perm := range permissions {
		if (perm.Resource == "*" || perm.Resource == resource) &&
			(perm.Action == "*" || perm.Action == action) {
			return true
		}
	}

	return false
}

func (a *AuthService) RevokeSession(sessionID string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.sessions, sessionID)
}

func (a *AuthService) GetUser(userID string) (*User, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	user, exists := a.users[userID]
	return user, exists
}

func (a *AuthService) GetActiveSessions() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.sessions)
}

func (a *AuthService) CleanupExpiredSessions() {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := time.Now()
	for id, session := range a.sessions {
		if now.After(session.ExpiresAt) {
			delete(a.sessions, id)
		}
	}
}

func (a *AuthService) RegisterAPIKey(key string, userID string) error {
	a.mu.RLock()
	user, exists := a.users[userID]
	a.mu.RUnlock()

	if !exists {
		return errors.New("user not found")
	}

	a.mu.Lock()
	a.apiKeys[key] = user
	a.mu.Unlock()

	return nil
}

func (a *AuthService) ValidateAPIKey(key string) (*User, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	user, exists := a.apiKeys[key]
	if !exists {
		return nil, errors.New("invalid API key")
	}

	return user, nil
}

func hashPassword(password string) string {
	return base64.StdEncoding.EncodeToString([]byte(password + "cloud-smart-tutor-salt"))
}

func generateSessionID() string {
	return fmt.Sprintf("sess-%d", time.Now().UnixNano())
}

func generateToken(userID, sessionID string, role Role) string {
	data := fmt.Sprintf("%s:%s:%s:%d", userID, sessionID, role, time.Now().UnixNano())
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func ExtractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	return ""
}

func ExtractTokenFromQuery(r *http.Request) string {
	return r.URL.Query().Get("token")
}

func ExtractAPIKeyFromHeader(r *http.Request) string {
	return r.Header.Get("X-API-Key")
}

func ExtractAPIKeyFromQuery(r *http.Request) string {
	return r.URL.Query().Get("api_key")
}

func Middleware(authService *AuthService, requiredRole Role, requiredResource string, requiredAction string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			token := ExtractTokenFromHeader(r)
			if token == "" {
				token = ExtractTokenFromQuery(r)
			}

			if token != "" {
				session, err := authService.ValidateToken(token)
				if err != nil {
					http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
					return
				}

				if !authService.CheckPermission(session.Role, requiredResource, requiredAction) {
					http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
					return
				}

				next(w, r)
				return
			}

			apiKey := ExtractAPIKeyFromHeader(r)
			if apiKey == "" {
				apiKey = ExtractAPIKeyFromQuery(r)
			}

			if apiKey != "" {
				user, err := authService.ValidateAPIKey(apiKey)
				if err != nil {
					http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
					return
				}

				if !authService.CheckPermission(user.Role, requiredResource, requiredAction) {
					http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
					return
				}

				next(w, r)
				return
			}

			http.Error(w, "Unauthorized: no authentication credentials provided", http.StatusUnauthorized)
		}
	}
}
