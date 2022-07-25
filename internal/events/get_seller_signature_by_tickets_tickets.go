package events

import (
	Cfg "couriergate/configs"
	DB "couriergate/internal/db"
	Import "couriergate/internal/import"
	"errors"
	"fmt"
)

// ----- получить подпись продавца для УПД -----
// 01
func GetSellerSignatureByTicketsTickets(sourceID string) error {

	ticketDetailsArr, err := Import.GetTicketsTickets(sourceID)
	if err != nil {
		return errors.New("Error (EV-090101): " + fmt.Sprintf("%s\n", err))
	}
	//fmt.Printf("%s\n", ticketDetailsArr)

	for _, v := range ticketDetailsArr {
		// UpdInvoiceSellerInfo
		if (v.Type == 19) || (v.Type == 21) {
			//fmt.Printf("%d\n", v)
			for _, w := range v.Signature {

				// ----- записываем образ ЭЦП в таблицу -----
				_, err = DB.DB_COURIER.Exec("INSERT INTO [Document.Signatures] (Service, DocumentID,Type,Uuid,Content) VALUES (@p1,@p2,'UpdInvoiceSellerInfo',@p3,@p4)", Cfg.Service, sourceID, w.Uuid, w.Content)
				if err != nil {
					return errors.New("Error (EV-090102): " + fmt.Sprintf("%s\n", err))
				}

				// ----- устанавливаем флаг получения ЭЦП -----
				_, err = DB.DB_COURIER.Exec("UPDATE [Document.In] SET SignReceived = 1 WHERE (Service = @p1) AND (DocumentID = @p2)", Cfg.Service, sourceID)
				if err != nil {
					return errors.New("Error (EV-090103): " + fmt.Sprintf("%s\n", err))
				}

			}
		}
	}

	return nil
}
