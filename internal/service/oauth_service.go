package service

import "go-security/internal/model"

type OauthService interface {
	Introspect(token string) (string, error)
	GetSsoProfile(token string) (*model.SsoUserProfile, error)
}
