package impl

import (
	"encoding/json"
	"errors"
	"go-security/internal/model"
	"go-security/internal/properties"
	"go-security/internal/service"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type TokenIntrospection struct {
	UserId string `json:"user_id"`
}

type ssoService struct {
	properties properties.SsoProperties
}

func NewSsoService(props properties.SsoProperties) service.OauthService {
	return &ssoService{properties: props}
}

func (s *ssoService) Introspect(token string) (string, error) {
	formData := url.Values{
		"token": {token},
	}

	req, err := http.NewRequest(http.MethodPost, s.properties.IntrospectionUrl, ioutil.NopCloser(strings.NewReader(formData.Encode())))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(s.properties.ClientId, s.properties.ClientSecret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to introspect token")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var userIdIntrospected TokenIntrospection
	err = json.Unmarshal(body, &userIdIntrospected)
	if err != nil {
		return "", err
	}

	return userIdIntrospected.UserId, nil
}

func (s *ssoService) GetSsoProfile(token string) (*model.SsoUserProfile, error) {
	var profile model.SsoUserProfile

	client := &http.Client{}
	req, err := http.NewRequest("GET", s.properties.UserProfileUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("no ok status code")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}
