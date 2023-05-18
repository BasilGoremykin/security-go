package model

type SsoUserProfile struct {
	IsOrganisation             string `json:"is_organisation"`
	Guid                       string `json:"guid"`
	MnpOperator                string `json:"mnp:operator"`
	MnpRegion                  string `json:"mnp:region"`
	Name                       string `json:"name"`
	GivenName                  string `json:"given_name"`
	MiddleName                 string `json:"middle_name"`
	LastName                   string `json:"last_name"`
	Phone                      string `json:"phone"`
	Email                      string `json:"email"`
	AccountNumber              string `json:"account:number"`
	AccountServiceProviderCode string `json:"account:service_provider_code"`
	Description                string `json:"description"`
	Picture                    string `json:"picture"`
}
