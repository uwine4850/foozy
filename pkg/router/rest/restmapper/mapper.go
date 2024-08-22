package restmapper

import (
	"github.com/uwine4850/foozy/pkg/interfaces/itypeopr"
	"github.com/uwine4850/foozy/pkg/router/rest"
)

// JsonToMessage converts JSON data into the selected message.
// It is important that the message is safe.
func JsonToMessage(jsonData *map[string]interface{}, dto *rest.DTO, output itypeopr.IPtr) error {
	if err := rest.DeepCheckSafeMessage(dto, output); err != nil {
		return err
	}
	if err := FillMessageFromMap(jsonData, output); err != nil {
		return err
	}
	return nil
}
