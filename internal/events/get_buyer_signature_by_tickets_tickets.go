package events

import (
	Cfg "couriergate/configs"
	DB "couriergate/internal/db"
	Import "couriergate/internal/import"
	"errors"
	"fmt"
)

// ----- получить подпись покупателя для УПД -----
// 01
func GetBuyerSignatureByTicketsTickets(sourceID string) error {

	ticketDetailsArr, err := Import.GetTicketsTickets(sourceID)
	if err != nil {
		return errors.New("Error (EV-050101): " + fmt.Sprintf("%s\n", err))
	}

	for _, v := range ticketDetailsArr {
		// UpdInvoiceBuyerInfo
		if (v.Type == 20) || (v.Type == 22) {
			//fmt.Printf("%d\n", v)
			for _, w := range v.Signature {

				// ----- записываем образ ЭЦП в таблицу -----
				_, err = DB.DB_COURIER.Exec("INSERT INTO [Document.Signatures] (Service, DocumentID,Type,Uuid,Content) VALUES (@p1, @p2,'UpdInvoiceBuyerInfo',@p3,@p4)", Cfg.Service, sourceID, w.Uuid, w.Content)
				if err != nil {
					return errors.New("Error (EV-050102): " + fmt.Sprintf("%s\n", err))
				}

				// ----- устанавливаем флаг получения ЭЦП -----
				_, err = DB.DB_COURIER.Exec("UPDATE [Document.Out] SET SignReceived = 1 WHERE (Service = @p1) AND (DocumentID = @p2)", Cfg.Service, sourceID)
				if err != nil {
					return errors.New("Error (EV-050103): " + fmt.Sprintf("%s\n", err))
				}

			}
		}
	}

	return nil
}
