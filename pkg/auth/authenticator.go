package auth

import (
	"fmt"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/warber/sailor/pkg/authstore"
)

type Authenticator struct {
	authStore authstore.AuthStore
}

func New(authstore authstore.AuthStore) *Authenticator {
	return &Authenticator{authStore: authstore}
}
func (a *Authenticator) Auth(apiEndpoint, apiToken string) (*api.APISet, error) {
	apiSet, err := api.New(apiEndpoint, api.WithAuthToken(apiToken))
	if err != nil {
		return nil, fmt.Errorf("could not create API client: %w", err)
	}
	_, authErr := apiSet.AuthV1().Authenticate()
	if err != nil {
		return nil, fmt.Errorf("could not authenticate: %s", *authErr.Message)
	}

	authInfo := authstore.AuthInfo{
		APIEndpoint: apiEndpoint,
		APIToken:    apiToken,
	}
	err = a.authStore.Save(&authInfo)
	if err != nil {
		return nil, fmt.Errorf("could not save auth info: %w", err)
	}

	return apiSet, nil
}

func (a *Authenticator) ReAuth() (*api.APISet, error) {
	authInfo, err := a.authStore.Read()
	if err != nil {
		return nil, fmt.Errorf("could not authenticate: %w", err)
	}

	apiSet, err := api.New(authInfo.APIEndpoint, api.WithAuthToken(authInfo.APIToken))
	if err != nil {
		return nil, fmt.Errorf("could not create API client: %w", err)
	}
	return apiSet, nil
}
