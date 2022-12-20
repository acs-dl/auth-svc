package connector

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Connector struct {
	client  *http.Client
	baseUrl string
}

func NewConnector(baseUrl string) *Connector {
	return &Connector{http.DefaultClient, baseUrl}
}

func (c *Connector) post(endpoint, auth string, dst interface{}) error {
	return c.upsert(http.MethodPost, endpoint, auth, dst)
}

func (c *Connector) upsert(method, endpoint, auth string, dst interface{}) error {
	// creating request
	request, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create connector request")
	}

	// setting headers
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", auth))

	// sending request
	response, err := c.client.Do(request)
	if err != nil {
		return errors.Wrap(err, "failed to process request")
	}
	if response == nil {
		return errors.New("failed to process request: response is nil")
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		_ = response.Body.Close()
		return errors.New("Bad status")
	}

	defer func(Body io.ReadCloser) {
		if tempErr := Body.Close(); tempErr != nil {
			err = tempErr
		}
	}(response.Body)

	// if destination is nil, we don`t read the response
	if dst == nil {
		return nil
	}

	// parsing response
	raw, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}

	return json.Unmarshal(raw, &dst)
}
