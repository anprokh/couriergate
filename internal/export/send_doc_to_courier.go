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

	//fmt.Println("*****************************************************************************************************************************")
	url := Cfg.ServiceURL + "/api/document/add/" + destinationID
	//fmt.Println(url)

	value, err := json.Marshal(DocumentCard)
	if err != nil {
		return "", errors.New("Error (EU-040101): " + fmt.Sprintf("%s\n", err))
	}
	//fmt.Println(string(value))

	data, err := Web.SendPostRequest(url, string(value), "application/json", Cfg.TOKEN)
	if err != nil {
		return "", errors.New("Error (EU-040102): " + fmt.Sprintf("%s\n", err))
	}
	//fmt.Println(data)

	return data, nil
}
