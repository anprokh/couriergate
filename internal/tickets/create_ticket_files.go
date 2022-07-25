package tickets

import (
	Cfg "couriergate/configs"
	DB "couriergate/internal/db"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"time"
)

// ----- создать файлы квитанций -----
// 01
func CreateTicketFiles() error {

	data, err := Get_Ticket_Index() // массив TicketIndexOptions
	if err != nil {
		return errors.New("Error (TU-030101): " + fmt.Sprintf("%s\n", err))
	}

	//data1, _ := json.Marshal(data)
	//fmt.Println(string(data1))

	for _, TI := range data {
		//fmt.Printf("%d : %s\n", i, TI)

		var ticketType string
		switch TI.Type {
		case 2:
			ticketType = "PDOTPR"
		case 3:
			ticketType = "PDPOL"
		case 4:
			ticketType = "UVUTOCH"
		case 18:
			ticketType = "UPD"
		case 19:
			ticketType = "SELLERINFO"
		case 20:
			ticketType = "BUYERINFO"
		default:
			ticketType = "UNKNOWN"
		}

		ticketName := fmt.Sprintf("%d.%d.%s.Ticket", TI.Id, TI.DocumentId, ticketType)
		ticketFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, ticketName)
		fmt.Println(ticketFullName)

		// пропускаем квитанции для которых найден файл
		if _, err := os.Stat(ticketFullName); err == nil {
			fmt.Println("Ticket file definitely exists!")
			continue
		}

		Ticket, err := Get_Ticket_Ticket(TI) // структура TicketDetails
		if err != nil {
			return errors.New("Error (TU-030102): " + fmt.Sprintf("%s\n", err))
		}

		data64 := Ticket.Content.Content

		// ----- декодируем содержимое квитанции из Base64, ожидаем что там win-1251 -----
		sDec, err := base64.StdEncoding.DecodeString(data64) // string -> []byte
		if err != nil {
			return errors.New("Error (TU-030103): " + fmt.Sprintf("%s\n", err))
		}
		//fmt.Printf("sDec: %s\n", sDec)

		f, err := os.Create(ticketFullName)
		if err != nil {
			return errors.New("Error (TU-030104): " + fmt.Sprintf("%s\n", err))
		}
		defer f.Close()

		// кодируем в Windows1251
		//		enc := charmap.Windows1251.NewEncoder()
		//		fileContent1251, _ := enc.String(ticketContent)

		_, err = f.Write([]byte(sDec))
		if err != nil {
			return errors.New("Error (TU-030105): " + fmt.Sprintf("%s\n", err))
		}

		// время получения квитанции
		DATE_TIME := time.Now().String()
		DATE_TIME = fmt.Sprintf("%sT%s", DATE_TIME[:10], DATE_TIME[11:23])

		// ----- записываем в sql время получения квитанции и ее Id -----
		switch ticketType {
		case "PDPOL":
			_, err = DB.DB_COURIER.Exec("UPDATE [Document.Out] SET TicketReceived = CAST(@p1 AS datetime), TicketID = @p2 WHERE (Service = @p3) AND (DocumentID = @p4)", DATE_TIME, TI.Id, Cfg.Service, TI.DocumentId)
		case "UVUTOCH":
			_, err = DB.DB_COURIER.Exec("UPDATE [Document.Out] SET TicketUvutochID = @p1 WHERE (Service = @p2) AND (DocumentID = @p3)", TI.Id, Cfg.Service, TI.DocumentId)
		default:
		}
		if err != nil {
			return errors.New("Error (TU-030108): " + fmt.Sprintf("%s\n", err))
		}
	}

	return nil
}
