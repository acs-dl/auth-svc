package connector

import (
	"fmt"
	"net/http"
	"strings"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

const validateEndpoint = "validate"

func (c *Connector) ValidateToken(req *http.Request) (err error) {
	defer func(e error) {
		err = errors.Wrap(err, "failed to send request via connector")
	}(err)

	endpoint := fmt.Sprintf("%s/%s", c.baseUrl, validateEndpoint)

	splitAuthHeader := strings.Split(req.Header.Get("Authorization"), " ")
	if len(splitAuthHeader) != 2 {
		return errors.New("No auth token provided")
	}

	token := splitAuthHeader[1]
	if token == "" {
		return errors.New("No auth token provided")
	}

	err = c.post(endpoint, token, nil)
	return
}
