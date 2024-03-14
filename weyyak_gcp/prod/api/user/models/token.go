package models

import (
	"time"

	"github.com/go-oauth2/oauth2/v4"
)

// NewToken create to token model instance
func NewToken() *Token {
	return &Token{}
}

// Token token model
type Token struct {
	ClientID            string        `bson:"ClientID"`
	UserID              string        `bson:"UserID"`
	RedirectURI         string        `bson:"RedirectURI"`
	Scope               string        `bson:"Scope"`
	Code                string        `bson:"Code"`
	CodeChallenge       string        `bson:"CodeChallenge"`
	CodeChallengeMethod string        `bson:"CodeChallengeMethod"`
	CodeCreateAt        time.Time     `bson:"CodeCreateAt"`
	CodeExpiresIn       time.Duration `bson:"CodeExpiresIn"`
	Access              string        `bson:"Access"`
	AccessCreateAt      time.Time     `bson:"AccessCreateAt"`
	AccessExpiresIn     time.Duration `bson:"AccessExpiresIn"`
	Refresh             string        `bson:"Refresh"`
	RefreshCreateAt     time.Time     `bson:"RefreshCreateAt"`
	RefreshExpiresIn    time.Duration `bson:"RefreshExpiresIn"`
	// Below fields added as per requirement
	GrantType        string    `bson:"GrantType"`
	Username         string    `bson:"Username"`
	Password         string    `bson:"Password"`
	DeviceID         string    `bson:"DeviceID"`
	DeviceName       string    `bson:"DeviceName"`
	DevicePlatform   string    `bson:"DevicePlatform"`
	LanguageId       string    `bson:"LanguageId"`
	Role             string    `bson:"Role"`
	IsBackOfficeUser bool      `bson:"IsBackOfficeUser"`
	ExpiresAt        time.Time `bson:"ExpiresAt"`
}

// New create to token model instance
func (t *Token) New() oauth2.TokenInfo {
	return NewToken()
}

// GetClientID the client id
func (t *Token) GetClientID() string {
	return t.ClientID
}

// SetClientID the client id
func (t *Token) SetClientID(clientID string) {
	t.ClientID = clientID
}

// GetUserID the user id
func (t *Token) GetUserID() string {
	return t.UserID
}

// SetUserID the user id
func (t *Token) SetUserID(userID string) {
	t.UserID = userID
}

// GetRedirectURI redirect URI
func (t *Token) GetRedirectURI() string {
	return t.RedirectURI
}

// SetRedirectURI redirect URI
func (t *Token) SetRedirectURI(redirectURI string) {
	t.RedirectURI = redirectURI
}

// GetScope get scope of authorization
func (t *Token) GetScope() string {
	return t.Scope
}

// SetScope get scope of authorization
func (t *Token) SetScope(scope string) {
	t.Scope = scope
}

// GetCode authorization code
func (t *Token) GetCode() string {
	return t.Code
}

// SetCode authorization code
func (t *Token) SetCode(code string) {
	t.Code = code
}

// GetCodeCreateAt create Time
func (t *Token) GetCodeCreateAt() time.Time {
	return t.CodeCreateAt
}

// SetCodeCreateAt create Time
func (t *Token) SetCodeCreateAt(createAt time.Time) {
	t.CodeCreateAt = createAt
}

// GetCodeExpiresIn the lifetime in seconds of the authorization code
func (t *Token) GetCodeExpiresIn() time.Duration {
	return t.CodeExpiresIn
}

// SetCodeExpiresIn the lifetime in seconds of the authorization code
func (t *Token) SetCodeExpiresIn(exp time.Duration) {
	t.CodeExpiresIn = exp
}

// GetCodeChallenge challenge code
func (t *Token) GetCodeChallenge() string {
	return t.CodeChallenge
}

// SetCodeChallenge challenge code
func (t *Token) SetCodeChallenge(code string) {
	t.CodeChallenge = code
}

// GetCodeChallengeMethod challenge method
func (t *Token) GetCodeChallengeMethod() oauth2.CodeChallengeMethod {
	return oauth2.CodeChallengeMethod(t.CodeChallengeMethod)
}

// SetCodeChallengeMethod challenge method
func (t *Token) SetCodeChallengeMethod(method oauth2.CodeChallengeMethod) {
	t.CodeChallengeMethod = string(method)
}

// GetAccess access Token
func (t *Token) GetAccess() string {
	return t.Access
}

// SetAccess access Token
func (t *Token) SetAccess(access string) {
	t.Access = access
}

// GetAccessCreateAt create Time
func (t *Token) GetAccessCreateAt() time.Time {
	return t.AccessCreateAt
}

// SetAccessCreateAt create Time
func (t *Token) SetAccessCreateAt(createAt time.Time) {
	t.AccessCreateAt = createAt
}

// GetAccessExpiresIn the lifetime in seconds of the access token
func (t *Token) GetAccessExpiresIn() time.Duration {
	return t.AccessExpiresIn
}

// SetAccessExpiresIn the lifetime in seconds of the access token
func (t *Token) SetAccessExpiresIn(exp time.Duration) {
	t.AccessExpiresIn = exp
}

// GetRefresh refresh Token
func (t *Token) GetRefresh() string {
	return t.Refresh
}

// SetRefresh refresh Token
func (t *Token) SetRefresh(refresh string) {
	t.Refresh = refresh
}

// GetRefreshCreateAt create Time
func (t *Token) GetRefreshCreateAt() time.Time {
	return t.RefreshCreateAt
}

// SetRefreshCreateAt create Time
func (t *Token) SetRefreshCreateAt(createAt time.Time) {
	t.RefreshCreateAt = createAt
}

// GetRefreshExpiresIn the lifetime in seconds of the refresh token
func (t *Token) GetRefreshExpiresIn() time.Duration {
	return t.RefreshExpiresIn
}

// SetRefreshExpiresIn the lifetime in seconds of the refresh token
func (t *Token) SetRefreshExpiresIn(exp time.Duration) {
	t.RefreshExpiresIn = exp
}

// GetDeviceID the lifetime in seconds of the refresh token
func (t *Token) GetDeviceID() string {
	return t.DeviceID
}

// SetDeviceID the lifetime in seconds of the refresh token
func (t *Token) SetDeviceID(deviceId string) {
	t.DeviceID = deviceId
}

// GetGrantType the lifetime in seconds of the refresh token
func (t *Token) GetGrantType() string {
	return t.GrantType
}

// SetGrantType the lifetime in seconds of the refresh token
func (t *Token) SetGrantType(grantType string) {
	t.GrantType = grantType
}

// GetUsername name of the user in token
func (t *Token) GetUsername() string {
	return t.Username
}

// SetUsername name of the user in token
func (t *Token) SetUsername(username string) {
	t.Username = username
}

// GetPassword password of the user in token
func (t *Token) GetPassword() string {
	return t.Password
}

// SetPassword password of the user in token
func (t *Token) SetPassword(password string) {
	t.Password = password
}

// GetDeviceName device of the user in token
func (t *Token) GetDeviceName() string {
	return t.DeviceName
}

// SetDeviceName device of the user in token
func (t *Token) SetDeviceName(deviceName string) {
	t.DeviceName = deviceName
}

// GetDevicePlatform device of the user in token
func (t *Token) GetDevicePlatform() string {
	return t.DevicePlatform
}

// SetDevicePlatform device of the user in token
func (t *Token) SetDevicePlatform(devicePlatform string) {
	t.DevicePlatform = devicePlatform
}

// GetLanguageId language of the user
func (t *Token) GetLanguageId() string {
	return t.LanguageId
}

// SetLanguageId language of the user setting
func (t *Token) SetLanguageId(language string) {
	t.LanguageId = language
}

// GetRole role of the user
func (t *Token) GetRole() string {
	return t.Role
}

// SetRole role of the user setting
func (t *Token) SetRole(role string) {
	t.Role = role
}

// GetIsbackofficeuser user role of the user
func (t *Token) GetIsBackOfficeUser() bool {
	return t.IsBackOfficeUser
}

// SetIsBackOfficeUser userrole of the user setting
func (t *Token) SetIsBackOfficeUser(isBackOfficeUser bool) {
	t.IsBackOfficeUser = isBackOfficeUser
}

// GetExpiresAt create Time
func (t *Token) GetExpiresAt() time.Time {
	return t.ExpiresAt
}

// SetExpiresAt create Time
func (t *Token) SetExpiresAt(expiresAt time.Time) {
	t.ExpiresAt = expiresAt
}
