package Import

import (
	Cfg "couriergate/configs"
	Web "couriergate/internal/web"
	Models "couriergate/models"
	"encoding/json"
	"errors"
	"fmt"
)

// ----- получить список квитанций с подписями для документа в массиве структур TicketDetails -----
// 01
func GetTicketsTickets(sourceID string) ([]Models.TicketDetailsOptions, error) {

	var TicketDetailsArr []Models.TicketDetailsOptions

	url := Cfg.ServiceURL + "/api/tickets/tickets/" + sourceID
	data, err := Web.SendGetRequest(url, "application/json", Cfg.TOKEN)
	if err != nil {
		return TicketDetailsArr, errors.New("Error (IU-040101): " + fmt.Sprintf("%s\n", err))
	}
	//fmt.Println(data)

	if err = json.Unmarshal([]byte(data), &TicketDetailsArr); err != nil {
		return TicketDetailsArr, errors.New("Error (IU-040102): " + fmt.Sprintf("%s\n", err))
	}
	//fmt.Printf("TicketDetailsArr: %s\n", TicketDetailsArr)

	return TicketDetailsArr, nil
}
