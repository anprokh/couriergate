package export

import (
	Cfg "couriergate/configs"
	Web "couriergate/internal/web"
	Models "couriergate/models"
	"encoding/json"
	"errors"
	"fmt"
)

// ----- отправить документ в "Courier" EDM -----
// 01
func SendDocToCourier(destinationID string, DocumentCard Models.DocumentCardOptions) (string, error) {

	url := Cfg.ServiceURL + "/api/document/add/" + destinationID

	value, err := json.Marshal(DocumentCard)
	if err != nil {
		return "", errors.New("Error (EU-040101): " + fmt.Sprintf("%s\n", err))
	}

	data, err := Web.SendPostRequest(url, string(value), "application/json", Cfg.TOKEN)
	if err != nil {
		return "", errors.New("Error (EU-040102): " + fmt.Sprintf("%s\n", err))
	}

	return data, nil
}
