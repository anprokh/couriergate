package tickets

import (
	Cfg "couriergate/configs"
	Web "couriergate/internal/web"
	Models "couriergate/models"
	"encoding/json"
	"errors"
	"fmt"
)

// ----- Возвращает информацию (и содержимое) квитанции -----
// 01
func Get_Ticket_Ticket(TI Models.TicketIndexOptions) (Models.TicketDetailsOptions, error) {

	var TicketData Models.TicketDetailsOptions

	url := Cfg.ServiceURL + "/api/tickets/ticket/" + fmt.Sprintf("%d", TI.Id)

	data, err := Web.SendGetRequest(url, "application/json", Cfg.TOKEN)
	if err != nil {
		return TicketData, errors.New("Error (TU-020101): " + fmt.Sprintf("%s\n", err))
	}

	if err = json.Unmarshal([]byte(data), &TicketData); err != nil {
		return TicketData, errors.New("Error (TU-020102): " + fmt.Sprintf("%s\n", err))
	}

	return TicketData, nil
}
