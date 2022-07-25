package clarification

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
	"strings"
)

// ----- создать файлы УОУ (уведомление об уточнении) -----
//
func CreateClarificationFiles() error {

	// ----- формируем текст запроса -----
	requestText := fmt.Sprintf("SELECT ROW_ID, DocumentID, ISNULL (RejectionReason, 'ошибка обработки') FROM [Document.In] (NOLOCK) WHERE (Action = 'Clarification') AND (ResponseFileCreated = 0) AND (Service = '%s') AND (Account = '%s')", Cfg.Service, Cfg.Account)
	//fmt.Println(requestText)

	// выполнение запроса
	rows, err := DB.DB_COURIER.Query(requestText)
	if err != nil {
		return errors.New("Error (CU-010101): SQL request failure, " + fmt.Sprintf("%s\n", err))
	}
	defer rows.Close()

	// ----- выборка и обработка данных запроса -----
	var DocumentID, RejectionReason string
	var ROW_ID int

	for rows.Next() {

		err := rows.Scan(&ROW_ID, &DocumentID, &RejectionReason)
		if err != nil {
			return errors.New("Error (CU-010102): " + fmt.Sprintf("%s\n", err))
		}
		//fmt.Printf("%d  %s : %s\n", ROW_ID, DocumentID, RejectionReason)

		url := Cfg.ServiceURL + "/api/tickets/createclarification/" + fmt.Sprintf("%s", DocumentID)
		//fmt.Println(url)

		RejectionReason = strings.Replace(RejectionReason, "\"", "\\\"", -1)
		requestBody := fmt.Sprintf("{ \"comment\":\"%s\" }", RejectionReason)

		data, err := Web.SendPostRequest(url, requestBody, "application/json", Cfg.TOKEN)
		if err != nil {
			return errors.New("Error (CU-010103): " + fmt.Sprintf("%s\n", err))
		}
		//fmt.Println(data)

		var FileContentData Models.SignedContentOptions

		if err = json.Unmarshal([]byte(data), &FileContentData); err != nil {
			fmt.Printf("Error (CU-010104): convert post data to SignedContentOptions type failure")
			continue
		}
		//fmt.Printf("FileContentData: %s\n", FileContentData)

		Filename := FileContentData.FileName
		data64 := FileContentData.Content

		// ----- декодируем содержимое документа из Base64 и приводим к win-1251 -----
		sDec, err := base64.StdEncoding.DecodeString(data64) // string -> []byte
		if err != nil {
			return errors.New("Error (CU-010105): " + fmt.Sprintf("%s\n", err))
		}
		//fmt.Printf("sDec: %s\n", sDec)

		// ----- записываем xml-файл документа -----
		fullFileName := fmt.Sprintf("%s\\%s.%s.Clarification", Cfg.ExPath, Filename, DocumentID)
		fmt.Println(fullFileName)

		f, err := os.Create(fullFileName)
		if err != nil {
			return errors.New("Error (CU-010106): " + fmt.Sprintf("%s\n", err))
		}
		defer f.Close()

		_, err = f.Write(sDec)
		if err != nil {
			return errors.New("Error (CU-010107): " + fmt.Sprintf("%s\n", err))
		}

		// ----- устанавливаем флаг создания файла -----
		_, err = DB.DB_COURIER.Exec("UPDATE [Document.In] SET ResponseFileCreated = 1 WHERE ROW_ID = @p1", ROW_ID)
		if err != nil {
			return errors.New("Error (CU-010108): " + fmt.Sprintf("%s\n", err))
		}
	}

	return nil
}
