package events

import (
	Cfg "couriergate/configs"
	DB "couriergate/internal/db"
	"errors"
	"fmt"
	"github.com/fatih/color"
)

// ----- получить и обработать pdf файлы документов -----
// 01
func Processed_Document_Out_Pdf() error {

	// ----- формируем текст запроса -----
	requestText := fmt.Sprintf("SELECT ROW_ID, DocumentID FROM [Document.Out] (NOLOCK) "+
		"WHERE (DocumentStatus = 'Signed') AND (PdfReceived = 0) AND (ROW_ID >= 3119) AND (Service = '%s') AND (Account = '%s') ORDER BY ROW_ID", Cfg.Service, Cfg.Account)
	//fmt.Println(requestText)

	// выполнение запроса
	rows, err := DB.DB_COURIER.Query(requestText)
	if err != nil {
		return errors.New("Error (EV-040101): SQL request failure, " + fmt.Sprintf("%s\n", err))
	}
	defer rows.Close()

	// ----- выборка и обработка данных запроса -----
	var DocumentID string
	var ROW_ID int

	for rows.Next() {

		err := rows.Scan(&ROW_ID, &DocumentID)
		if err != nil {
			return errors.New("Error (EV-040102): " + fmt.Sprintf("%s\n", err))
		}
		//fmt.Printf("%d  %s\n", ROW_ID, DocumentID)

		_, err = GetDocumentPdf(DocumentID, "Document.Out", 1)
		if err != nil {
			fmt.Fprintf(color.Output, "%s Pdf %s not received: %s", color.RedString("[warning]"), color.CyanString(DocumentID), err)
		} else {
			fmt.Fprintf(color.Output, "%s Pdf %s successfully received\n", color.GreenString("[info]"), color.CyanString(DocumentID))
		}

	}

	return nil
}
