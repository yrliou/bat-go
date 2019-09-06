package httprequest

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

var payloadLimit10MB = int64(1024 * 1024 * 10)

// ReadBody reads a request body according to an interface and limits the size to 10MB
func ReadBody(body io.Reader, req interface{}) error {
	jsonString, err := ioutil.ReadAll(io.LimitReader(body, payloadLimit10MB))
	if err != nil {
		return errors.WithMessage(err, "Error reading body")
	}

	err = json.Unmarshal(jsonString, &req)
	if err != nil {
		return errors.WithMessage(err, "Error unmarshalling body")
	}
	return nil
}
