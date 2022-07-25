package events

import (
	Cfg "couriergate/configs"
	DB "couriergate/internal/db"
	Import "couriergate/internal/import"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"strconv"
)

// ----- обработать новые события из списка -----
// 01
func Processed_Document_Events() error {

	rows, err := DB.DB_COURIER.Query("SELECT Id, DocumentId, TicketType, EventType FROM [Document.Events] (NOLOCK) WHERE ([Processed] = 0) AND (Service = @p1) AND (Account = @p2) ORDER BY [Id]", Cfg.Service, Cfg.Account)
	if err != nil {
		return errors.New("Error (EV-020101): " + fmt.Sprintf("%s\n", err))
	}
	defer rows.Close()

	// ----- выборка и обработка данных запроса -----
	var Id, DocumentId int
	var TicketType, EventType string

	for rows.Next() {

		err := rows.Scan(&Id, &DocumentId, &TicketType, &EventType)
		if err != nil {
			return errors.New("Error (EV-020102): " + fmt.Sprintf("%s\n", err))
		}
		fmt.Printf("%d     %d : %s   %s\n", Id, DocumentId, TicketType, EventType)

		// ----- сохраняем полученные УПД в таблице [Document.In] -----
		// обрабатываем UpdInvoiceSellerInfo, UpdInvoice, UpdSellerInfo
		if ((TicketType == "UpdInvoiceSellerInfo") && (EventType == "TicketReceived")) ||
			((TicketType == "UpdInvoice") && (EventType == "TicketReceived")) ||
			((TicketType == "UpdSellerInfo") && (EventType == "TicketReceived")) {
			id_str := strconv.Itoa(DocumentId)
			err = Import.Import_Document_In(id_str)
			if err != nil {
				return errors.New("Error (EV-020103): " + fmt.Sprintf("%s\n", err))
			}

			// получим и сохраним предварительные (не подписанные нами) pdf файлы входящих документов
			_, err = GetDocumentPdf(id_str, "Document.In", 0)
			if err != nil {
				fmt.Fprintf(color.Output, "%s Unsigned Pdf %s not received: %s", color.RedString("[warning]"), color.CyanString(id_str), err)
			} else {
				fmt.Fprintf(color.Output, "%s Unsigned Pdf %s successfully received\n", color.GreenString("[info]"), color.CyanString(id_str))
			}
		}

		// устанавливаем статус документа
		if (EventType == "Created") || (EventType == "Signed") || (EventType == "Rejected") || (EventType == "ClarificationRequested") {
			//fmt.Printf("   ---------------> %d : %s\n", DocumentId, EventType)
			_, err = DB.DB_COURIER.Exec("UPDATE [Document.Out] SET DocumentStatus = @p1, StatusProcessed = 0 WHERE (Service = @p2) AND (DocumentID = @p3)", EventType, Cfg.Service, DocumentId)
			if err != nil {
				return errors.New("Error (EV-020104): " + fmt.Sprintf("%s\n", err))
			}
		}

		// помечаем событие как обработанное
		_, err = DB.DB_COURIER.Exec("UPDATE [Document.Events] SET Processed = 1 WHERE (Service = @p1) AND (Id = @p2)", Cfg.Service, Id)
		if err != nil {
			return errors.New("Error (EV-020105): " + fmt.Sprintf("%s\n", err))
		}
	}

	return nil
}
