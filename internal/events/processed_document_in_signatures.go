package events

import (
	Cfg "couriergate/configs"
	DB "couriergate/internal/db"
	"errors"
	"fmt"
	"github.com/fatih/color"
)

// ----- получить и сохранить подписи покупателей для УПД -----
// 01
func Processed_Document_In_Signatures() error {

	// ----- формируем текст запроса -----
	requestText := fmt.Sprintf("SELECT ROW_ID, DocumentID FROM [Document.In] (NOLOCK) "+
		"WHERE (Action = 'Accept') AND (SignReceived = 0) AND (Service = '%s') AND (Account = '%s') ORDER BY ROW_ID", Cfg.Service, Cfg.Account)

	// выполнение запроса
	rows, err := DB.DB_COURIER.Query(requestText)
	if err != nil {
		return errors.New("Error (EV-080101): SQL request failure, " + fmt.Sprintf("%s\n", err))
	}
	defer rows.Close()

	// ----- выборка и обработка данных запроса -----
	var DocumentID string
	var ROW_ID int

	for rows.Next() {

		err := rows.Scan(&ROW_ID, &DocumentID)
		if err != nil {
			return errors.New("Error (EV-080102): " + fmt.Sprintf("%s\n", err))
		}

		err = GetSellerSignatureByTicketsTickets(DocumentID)
		if err != nil {
			fmt.Fprintf(color.Output, "%s SellerSignature %s not received: %s", color.RedString("[warning]"), color.CyanString(DocumentID), err)
		} else {
			fmt.Fprintf(color.Output, "%s SellerSignature %s successfully received\n", color.GreenString("[info]"), color.CyanString(DocumentID))
		}

	}

	return nil
}
