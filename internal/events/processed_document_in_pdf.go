package events

import (
	Cfg "couriergate/configs"
	DB "couriergate/internal/db"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"time"
)

// ----- получить и обработать pdf файлы входящих документов -----
// 01
func Processed_Document_In_Pdf() error {

	// ----- формируем текст запроса -----
	requestText := fmt.Sprintf("SELECT ROW_ID, DocumentID, ISNULL(ActionProcessed, getdate()) AS ActionProcessed FROM [Document.In] (NOLOCK) "+
		"WHERE (Action = 'Accept') AND (PdfReceived = 0) AND (Service = '%s') AND (Account = '%s') ORDER BY ROW_ID", Cfg.Service, Cfg.Account)
	//fmt.Println(requestText)

	// выполнение запроса
	rows, err := DB.DB_COURIER.Query(requestText)
	if err != nil {
		return errors.New("Error (EV-070101): SQL request failure, " + fmt.Sprintf("%s\n", err))
	}
	defer rows.Close()

	// ----- выборка и обработка данных запроса -----
	var DocumentID string
	var ROW_ID int

	var ActionProcessed string

	for rows.Next() {

		err := rows.Scan(&ROW_ID, &DocumentID, &ActionProcessed)
		if err != nil {
			return errors.New("Error (EV-070102): " + fmt.Sprintf("%s\n", err))
		}
		//fmt.Printf("(1)  %d  %s  %s\n", ROW_ID, DocumentID, ActionProcessed)
		ActionProcessed = ActionProcessed[:19]

		// ----- обеспечиваем задержку 120 сек на прикрепление подписи в Courier -----
		layout := "2006-01-02T15:04:05"
		actionTime, err := time.Parse(layout, ActionProcessed)
		if err != nil {
			return errors.New("Error (EV-070103): " + fmt.Sprintf("%s\n", err))
		}
		//fmt.Println(actionTime)

		// компенсируем разницу между UTC и MSK
		actionTime = actionTime.Add(time.Hour * -3)

		t1 := time.Now()
		//fmt.Println(t1)

		diff := t1.Sub(actionTime).Seconds()
		fmt.Printf("diff:  %f\n", diff)

		if diff < 120 {
			continue
		}
		// ---------------------------------------------------------------------------

		_, err = GetDocumentPdf(DocumentID, "Document.In", 1)
		if err != nil {
			fmt.Fprintf(color.Output, "%s Pdf %s not received: %s", color.RedString("[warning]"), color.CyanString(DocumentID), err)
		} else {
			fmt.Fprintf(color.Output, "%s Pdf %s successfully received\n", color.GreenString("[info]"), color.CyanString(DocumentID))
		}

	}

	return nil
}
