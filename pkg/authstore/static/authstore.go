package static

import (
	"github.com/warber/sailor/pkg/authstore"
)

type StaticAuthStore struct {
	currentEndpoint string
	currentAPIToken string
	data            map[string]*authstore.AuthInfo
}

func New(currentEndpoint, currentAPIToken string) *StaticAuthStore {
	return &StaticAuthStore{currentEndpoint: currentEndpoint, currentAPIToken: currentAPIToken, data: make(map[string]*authstore.AuthInfo)}
}

func (i *StaticAuthStore) Save(info *authstore.AuthInfo) error {
	i.currentEndpoint = info.APIEndpoint
	i.currentAPIToken = info.APIToken
	return nil
}

func (i *StaticAuthStore) Read() (*authstore.AuthInfo, error) {

	return &authstore.AuthInfo{
		APIEndpoint: i.currentEndpoint,
		APIToken:    i.currentAPIToken,
	}, nil
}
