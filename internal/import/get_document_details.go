package Import

import (
	Cfg "couriergate/configs"
	Web "couriergate/internal/web"
	Models "couriergate/models"
	"encoding/json"
	"errors"
	"fmt"
)

// ----- получить данные по документу в структуре Document -----
// 01
func GetDocumentDetails(sourceID string) (Models.DocumentDetailsOptions, error) {

	var DocumentData Models.DocumentDetailsOptions

	url := Cfg.ServiceURL + "/api/document/details/" + sourceID
	data, err := Web.SendGetRequest(url, "application/json", Cfg.TOKEN)
	if err != nil {
		return DocumentData, errors.New("Error (IU-010101): " + fmt.Sprintf("%s\n", err))
	}

	if err = json.Unmarshal([]byte(data), &DocumentData); err != nil {
		return DocumentData, errors.New("Error (IU-010102): " + fmt.Sprintf("%s\n", err))
	}

	return DocumentData, nil
}
