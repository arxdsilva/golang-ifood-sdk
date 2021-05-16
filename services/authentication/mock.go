package authentication

import (
	"github.com/stretchr/testify/mock"
)

// AuthMock mock of auth service
type AuthMock struct {
	mock.Mock
}

// V2GetUserCode mock of auth service
func (a *AuthMock) V2GetUserCode() (uc *UserCode, err error) {
	args := a.Called()
	if res, ok := args.Get(0).(*UserCode); ok {
		return res, nil
	}
	return nil, args.Error(1)
}

func (a *AuthMock) V2Authenticate(authType, authCode, authCodeVerifier, refreshToken string) (c *V2Credentials, err error) {
	args := a.Called(authType, authCode, authCodeVerifier, refreshToken)
	if res, ok := args.Get(0).(*V2Credentials); ok {
		return res, nil
	}
	return nil, args.Error(1)
}

// Authenticate mock of auth service
func (a *AuthMock) Authenticate(user, pass string) (c *Credentials, err error) {
	args := a.Called(user, pass)
	if res, ok := args.Get(0).(*Credentials); ok {
		return res, nil
	}
	return nil, args.Error(1)
}

// Validate mock of auth service
func (a *AuthMock) Validate() (err error) {
	args := a.Called()
	return args.Error(0)
}

// GetToken mock of auth service
func (a *AuthMock) GetToken() (token string) {
	args := a.Called()
	return args.Get(0).(string)
}
