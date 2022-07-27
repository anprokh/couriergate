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

// ----- получить список событий -----
// 01
func Get_Events_Index() error {

	id := 0
	// ----- определяем максимальный id таблицы событий -----
	rows, err := DB.DB_COURIER.Query("SELECT ISNULL (MAX(Id), 0) FROM [Document.Events] (NOLOCK) WHERE Service = @p1 AND Account = @p2", Cfg.Service, Cfg.Account)
	if err != nil {
		return errors.New("Error (EV-010101): " + fmt.Sprintf("%s\n", err))
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return errors.New("Error (EV-010102): " + fmt.Sprintf("%s\n", err))
		}
	}

	var DocumentEventsArr []Models.DocumentEventOptions

	url := Cfg.ServiceURL + "/api/events/" + fmt.Sprintf("%d", id) + "?count=20"
	fmt.Println(url)

	data, err := Web.SendPostRequest(url, "", "application/json", Cfg.TOKEN)
	if err != nil {
		return errors.New("Error (EV-010103): " + fmt.Sprintf("%s\n", err))
	}

	if err = json.Unmarshal([]byte(data), &DocumentEventsArr); err != nil {
		return errors.New("Error (EV-010104): " + fmt.Sprintf("%s\n", err))
	}

	for _, E := range DocumentEventsArr {
		//fmt.Printf("%d : %s\n", i, E)

		var ticketType string
		switch E.TicketType {
		case 1:
			ticketType = "ReceiveNotice"
		case 2:
			ticketType = "SendConfirmation"
		case 3:
			ticketType = "ReceiveConfirmation"
		case 4:
			ticketType = "ClarificationNotice"
		case 5:
			ticketType = "Invoice"
		case 6:
			ticketType = "CorrectionInvoice"
		case 7:
			ticketType = "Torg12SellerTitle"
		case 8:
			ticketType = "Torg12BuyerTitle"
		case 9:
			ticketType = "AcceptanceCertificatePerformerTitle"
		case 10:
			ticketType = "AcceptanceCertificateCustomerTitle"
		case 11:
			ticketType = "Document"
		case 12:
			ticketType = "RejectSignature"
		case 13:
			ticketType = "AvoidanceRequest"
		case 14:
			ticketType = "AcceptSignature"
		case 15:
			ticketType = "ReceiveNoticeRoaming"
		case 16:
			ticketType = "ClarificationNoticeRoaming"
		case 18:
			ticketType = "UpdInvoice"
		case 19:
			ticketType = "UpdInvoiceSellerInfo"
		case 20:
			ticketType = "UpdInvoiceBuyerInfo"
		case 21:
			ticketType = "UpdSellerInfo"
		case 22:
			ticketType = "UpdBuyerInfo"
		case 23:
			ticketType = "UpdCorrectionInvoice"
		case 24:
			ticketType = "UpdCorrectionInvoiceSellerInfo"
		case 25:
			ticketType = "UpdCorrectionInvoiceBuyerInfo"
		case 26:
			ticketType = "UpdCorrectionSellerInfo"
		case 27:
			ticketType = "UpdCorrectionBuyerInfo"
		case 28:
			ticketType = "DocumentOfTransferOfWorkResultsPerformerInfo"
		case 29:
			ticketType = "DocumentOfTransferOfWorkResultsCustomerInfo"
		case 30:
			ticketType = "DocumentOfTransferOfGoodsSellerInfo"
		case 31:
			ticketType = "DocumentOfTransferOfGoodsBuyerInfo"
		default:
			ticketType = ""
		}

		var eventType string
		switch E.EventType {
		case 1:
			eventType = "Created"
		case 2:
			eventType = "Delivered"
		case 3:
			eventType = "Accepted"
		case 4:
			eventType = "Signed"
		case 5:
			eventType = "Rejected"
		case 6:
			eventType = "TicketReceived"
		case 7:
			eventType = "Received"
		case 8:
			eventType = "TicketSended"
		case 9:
			eventType = "MovedToTrash"
		case 10:
			eventType = "Deleted"
		case 11:
			eventType = "Sended"
		case 12:
			eventType = "DeliveredForReciever"
		case 13:
			eventType = "Revoked"
		case 14:
			eventType = "ReceivedRequestReview"
		case 15:
			eventType = "RequestedAvoidance"
		case 16:
			eventType = "RejectedAvoidance"
		case 17:
			eventType = "AcceptedAvoidance"
		case 18:
			eventType = "SignatureReject"
		case 19:
			eventType = "RestoredFromTrash"
		default:
			eventType = ""
		}

		// ----- добавляем запись в sql таблицу событий -----
		_, err = DB.DB_COURIER.Exec("INSERT INTO [Document.Events] (Service, Account, Id, DocumentId, TicketType, EventType, Date) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7)",
			Cfg.Service, Cfg.Account, E.Id, E.DocumentId, ticketType, eventType, E.Date)
		if err != nil {
			return errors.New("Error (EV-010105): " + fmt.Sprintf("%s\n", err))
		}
	}

	return nil
}
