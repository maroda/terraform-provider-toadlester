package toadlester

import (
	"io"
	"net/http"
	"strings"
	"time"
)

// Config is used to contain the Client
// use to check against interface{} parameters
// e.g.: config := interface{}(.*Config)
type Config struct {
	Client *APIClient
}

type APIClient struct {
	BaseURL    string
	HttpClient *http.Client
}

// Setting defines the change being made or value being read
// Name ::: The full envvar name (e.g. INT_SIZE)
// Value ::: The value to set (e.g. 100)
// Algo ::: Only used for reading, can be 'up' or 'down'
type Setting struct {
	Name  string
	Value string
	Algo  string
}

// NewAPIClient creates a client with the base url of the API
func NewAPIClient(url string) *APIClient {
	return &APIClient{
		BaseURL: url,
		HttpClient: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

// CreateType | UpdateType returns the body of the response from the API
// This API does not need initializing, so "Create" here means
// that a new looping buffer of numbers is being created.
func (c *APIClient) CreateType(setting *Setting) (string, error) {
	url := c.BaseURL + "/" + setting.Name + "/" + setting.Value
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *APIClient) ReadType(setting *Setting) (string, error) {
	parts := strings.Split(setting.Name, "_")
	baretype := strings.ToLower(parts[0])
	url := c.BaseURL + "/series/" + baretype + "/" + setting.Algo

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// UpdateType is identical to CreateType
// It updates the output by creating a new looping buffer
func (c *APIClient) UpdateType(setting *Setting) (string, error) {
	return c.CreateType(setting)
}

// DeleteType is identical to CreateType
// It deletes the current output by creating a new looping buffer
func (c *APIClient) DeleteType(setting *Setting) (string, error) {
	return c.CreateType(setting)
}
