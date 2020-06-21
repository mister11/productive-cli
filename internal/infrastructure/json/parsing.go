package json

import (
	"bytes"
	"io"
	"reflect"

	"github.com/google/jsonapi"
	"github.com/mister11/productive-cli/internal/utils"
)

func ToJsonEmbedded(model interface{}) []byte {
	jsonBuffer := new(bytes.Buffer)
	err := jsonapi.MarshalOnePayloadEmbedded(jsonBuffer, model)
	if err != nil {
		utils.ReportError("Marshaling of time entry payload failed.", err)
	}
	return jsonBuffer.Bytes()
}

func FromJsonMany(json io.Reader, t reflect.Type) []interface{} {
	models, err := jsonapi.UnmarshalManyPayload(json, t)
	if err != nil {
		utils.ReportError("Error parsing JSON response.", err)
	}

	return models
}
