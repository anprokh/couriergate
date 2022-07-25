package tickets

import (
	Cfg "couriergate/configs"
	Web "couriergate/internal/web"
	Models "couriergate/models"
	"encoding/json"
	"errors"
	"fmt"
)

// ----- получить список квитанций, на которые надо сформировать ответ -----
// 01
func Get_Ticket_Index() ([]Models.TicketIndexOptions, error) {

	var result []Models.TicketIndexOptions

	url := Cfg.ServiceURL + "/api/tickets"
	data, err := Web.SendGetRequest(url, "application/json", Cfg.TOKEN)
	if err != nil {
		return result, errors.New("Error (TU-010101): " + fmt.Sprintf("%s\n", err))
	}
	//fmt.Println(data)

	if err = json.Unmarshal([]byte(data), &result); err != nil {
		return result, errors.New("Error (TU-010102): " + fmt.Sprintf("%s\n", err))
	}

	return result, nil
}
