package signatures

import (
	Cfg "couriergate/configs"
	DB "couriergate/internal/db"
	"errors"
	"fmt"
)

// ----- получить имя сертификата из Document.In по id документа -----
// 01
func GetCertificateNameByDocumentID(documentID string) (string, error) {

	requestText := fmt.Sprintf("SELECT ISNULL (Certificate, '') FROM [Document.In] (NOLOCK) WHERE (Service = '%s') AND (DocumentID = %s) ORDER BY ROW_ID", Cfg.Service, documentID)
	//fmt.Println(requestText)

	// выполнение запроса
	rows, err := DB.DB_COURIER.Query(requestText)
	if err != nil {
		return "", errors.New("Error (SU-030101): SQL request failure, " + fmt.Sprintf("%s\n", err))
	}
	defer rows.Close()

	// ----- выборка и обработка данных запроса -----
	var certificateName string

	for rows.Next() {

		err := rows.Scan(&certificateName)
		if err != nil {
			return "", errors.New("Error (SU-030102): " + fmt.Sprintf("%s\n", err))
		}
		//fmt.Printf("%s\n", certificateName)
	}

	return certificateName, nil
}
