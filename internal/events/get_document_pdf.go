package events

import (
	Cfg "couriergate/configs"
	DB "couriergate/internal/db"
	Web "couriergate/internal/web"
	Models "couriergate/models"
	"encoding/json"
	"errors"
	"fmt"
)

// ----- получить PDF с печатной формой документа -----
// 01
func GetDocumentPdf(sourceID string, tableName string, setPdfReceived int) (string, error) {

	var FileContent Models.SignedContentOptions

	url := Cfg.ServiceURL + "/api/document/pdf/" + sourceID
	data, err := Web.SendGetRequest(url, "application/json", Cfg.TOKEN)
	if err != nil {
		return "", errors.New("Error (EV-030101): " + fmt.Sprintf("%s\n", err))
	}

	if err = json.Unmarshal([]byte(data), &FileContent); err != nil {
		return "", errors.New("Error (EV-030102): " + fmt.Sprintf("%s\n", err))
	}

	data64 := FileContent.Content

	// ----- записываем pdf base64 в таблицу, устанавливаем флаг получения pdf файла -----
	_, err = DB.DB_COURIER.Exec("UPDATE ["+tableName+"] SET PdfBase64 = @p1, PdfReceived = @p2 WHERE Service = @p3 AND DocumentID = @p4", data64, setPdfReceived, Cfg.Service, sourceID)
	if err != nil {
		return "", errors.New("Error (EV-030103): " + fmt.Sprintf("%s\n", err))
	}

	return "", nil
}
