package accept

import (
	Cfg "couriergate/configs"
	DB "couriergate/internal/db"
	Web "couriergate/internal/web"
	Models "couriergate/models"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// ----- создать файлы Титул покупателя -----
//
func CreateAcceptTitleFiles() error {

	// ----- формируем текст запроса -----
	requestText := fmt.Sprintf("SELECT ROW_ID, DocumentID, ISNULL (Certificate, '') FROM [Document.In] (NOLOCK) WHERE (Action = 'Accept') AND (ResponseFileCreated = 0) AND (Service = '%s') AND (Account = '%s')", Cfg.Service, Cfg.Account)

	// выполнение запроса
	rows, err := DB.DB_COURIER.Query(requestText)
	if err != nil {
		return errors.New("Error (AU-010101): SQL request failure, " + fmt.Sprintf("%s\n", err))
	}
	defer rows.Close()

	// ----- выборка и обработка данных запроса -----
	var DocumentID, CertificateName string
	var ROW_ID int

	for rows.Next() {

		err := rows.Scan(&ROW_ID, &DocumentID, &CertificateName)
		if err != nil {
			return errors.New("Error (AU-010102): " + fmt.Sprintf("%s\n", err))
		}
		//fmt.Printf("%d  %s  %s\n", ROW_ID, DocumentID, CertificateName)

		universalDocumentAcceptInfo, err := GetUniversalDocumentAcceptInfo(CertificateName)
		if err != nil {
			return errors.New("Error (AU-010103): " + fmt.Sprintf("%s\n", err))
		}

		url := Cfg.ServiceURL + "/v2.0/document/" + fmt.Sprintf("%s", DocumentID) + "/acceptTitle"

		switch Cfg.Service {
		case "CislinkTest":
			url = Cfg.ServiceURL + "/api/v2.0/document/" + fmt.Sprintf("%s", DocumentID) + "/acceptTitle"
		case "Cislink":
			url = Cfg.ServiceURL + "/api/v2.0/document/" + fmt.Sprintf("%s", DocumentID) + "/acceptTitle"
		default:
		}

		data, err := Web.SendPostRequest(url, universalDocumentAcceptInfo, "application/json", Cfg.TOKEN)
		if err != nil {
			return errors.New("Error (AU-010104): " + fmt.Sprintf("%s\n", err))
		}

		var FileContentData Models.SignedContentOptions

		if err = json.Unmarshal([]byte(data), &FileContentData); err != nil {
			fmt.Printf("Error (AU-010105): convert post data to FileContentOptions type failure")
			continue
		}

		Filename := FileContentData.FileName
		data64 := FileContentData.Content

		// ----- декодируем содержимое документа из Base64, ожидаем что там win-1251 -----
		sDec, err := base64.StdEncoding.DecodeString(data64) // string -> []byte
		if err != nil {
			return errors.New("Error (AU-010106): " + fmt.Sprintf("%s\n", err))
		}

		// ----- записываем xml-файл документа -----
		fullFileName := fmt.Sprintf("%s\\%s.%s.AcceptTitle", Cfg.ExPath, Filename, DocumentID)
		fmt.Println(fullFileName)

		f, err := os.Create(fullFileName)
		if err != nil {
			return errors.New("Error (AU-010107): " + fmt.Sprintf("%s\n", err))
		}
		defer f.Close()

		_, err = f.Write([]byte(sDec))
		if err != nil {
			return errors.New("Error (AU-010108): " + fmt.Sprintf("%s\n", err))
		}

		// ----- устанавливаем флаг создания файла -----
		_, err = DB.DB_COURIER.Exec("UPDATE [Document.In] SET ResponseFileCreated = 1 WHERE ROW_ID = @p1", ROW_ID)
		if err != nil {
			return errors.New("Error (AU-010109): " + fmt.Sprintf("%s\n", err))
		}
	}

	return nil
}
