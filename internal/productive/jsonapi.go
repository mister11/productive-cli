package productive

import (
	"bytes"
	"io"
	"reflect"

	"github.com/google/jsonapi"
)

func ToJsonEmbedded(model interface{}) ([]byte, error) {
	jsonBuffer := new(bytes.Buffer)
	err := jsonapi.MarshalOnePayloadEmbedded(jsonBuffer, model)
	if err != nil {
		return nil, err
	}
	return jsonBuffer.Bytes(), nil
}

func FromJsonMany(json io.Reader, t reflect.Type) ([]interface{}, error) {
	models, err := jsonapi.UnmarshalManyPayload(json, t)
	if err != nil {
		return nil, err
	}

	return models, nil
}
