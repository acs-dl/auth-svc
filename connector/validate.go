package connector

import (
	"fmt"
	"net/http"
	"strings"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

const validateEndpoint = "validate"

func (c *Connector) Validate(req *http.Request, modulePrefix string) (err error) {
	defer wrapErr(err)

	endpoint := fmt.Sprintf("%s/%s/%s", c.baseUrl, modulePrefix, validateEndpoint)

	splitedAuthHeader := strings.Split(req.Header.Get("Authorization"), " ")
	if len(splitedAuthHeader) != 2 {
		return
	}

	token := splitedAuthHeader[1]
	if token == "" {
		return
	}

	err = c.post(endpoint, token, nil)
	return
}

func wrapErr(err error) error {
	if err != nil {
		err = errors.Wrap(err, "failed to send message")
	}

	return err
}
