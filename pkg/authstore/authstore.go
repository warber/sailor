package authstore

type AuthInfo struct {
	APIEndpoint string `json:"apiEndpoint"`
	APIToken    string `json:"apiToken"`
}

type AuthStore interface {
	Save(info *AuthInfo) error
	Read() (*AuthInfo, error)
}
