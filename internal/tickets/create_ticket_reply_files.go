package tickets

import (
	Cfg "couriergate/configs"
	DB "couriergate/internal/db"
	Web "couriergate/internal/web"
	Models "couriergate/models"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ----- создать файлы ответных квитанций -----
// 01
func CreateTicketReplyFiles() error {

	var fileNames []string

	files, err := ioutil.ReadDir(Cfg.ExPath)
	if err != nil {
		return errors.New("Error (TU-040101): " + fmt.Sprintf("%s\n", err))
	}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) == ".Ticket" {
			fileNames = append(fileNames, file.Name())
		}
	}

	for _, name := range fileNames {

		fmt.Println("---------------------------------------------------------------------------------------------------------------")
		filepath := fmt.Sprintf("%s\\%s", Cfg.ExPath, name)
		fmt.Println(filepath)

		replyName := fmt.Sprintf("%sReply", name)
		replyFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, replyName)
		fmt.Println(replyFullName)

		// пропускаем квитанции для которых найден ответный файл
		if _, err := os.Stat(replyFullName); err == nil {
			fmt.Println("Ticket reply file definitely exists!")
			continue
		}

		// ----- выделяем из имени файла Id квитанции -----
		var i = strings.Index(name, ".")
		var ticketId string

		if i == -1 {
			// в имени файла нет "."
			continue
		} else {
			var x1 = strings.Split(name, ".")
			ticketId = x1[0]
		}
		//fmt.Printf("ticketId: %s\n", ticketId)

		ticketReply, err := Get_Ticket_Reply(ticketId) // структура SignedContent
		if err != nil {
			fmt.Printf("Error (TU-040102): %s\n", err)
			continue
		}
		//fmt.Printf("ticketReply: %s\n", ticketReply)

		data64 := ticketReply.Content

		// ----- декодируем содержимое ответной квитанции из Base64, ожидаем что там win-1251 -----
		sDec, err := base64.StdEncoding.DecodeString(data64) // string -> []byte
		if err != nil {
			return errors.New("Error (TU-040103): " + fmt.Sprintf("%s\n", err))
		}
		//fmt.Printf("sDec: %s\n", sDec)

		f, err := os.Create(replyFullName)
		if err != nil {
			return errors.New("Error (TU-040104): " + fmt.Sprintf("%s\n", err))
		}
		defer f.Close()

		// кодируем в Windows1251
		//		enc := charmap.Windows1251.NewEncoder()
		//		fileContent1251, _ := enc.String(replyContent)

		_, err = f.Write([]byte(sDec))
		if err != nil {
			return errors.New("Error (TU-040105): " + fmt.Sprintf("%s\n", err))
		}

		i = strings.Index(name, "PDPOL")
		if i == -1 {
			// пропускаем квитанции не ПДП
			continue
		}

		// время получения ответной квитанции
		DATE_TIME := time.Now().String()
		DATE_TIME = fmt.Sprintf("%sT%s", DATE_TIME[:10], DATE_TIME[11:23])

		// преобразовываем id в int64
		ticketId64, err := strconv.ParseInt(ticketId, 10, 64)
		if err != nil {
			return errors.New("Error (TU-040106): " + fmt.Sprintf("%s\n", err))
		}

		// ----- записываем в sql время получения ответной квитанции -----
		_, err = DB.DB_COURIER.Exec("UPDATE [Document.Out] SET TicketReplyReceived = CAST(@p1 AS datetime) WHERE (Service = @p2) AND (TicketID = @p3)", DATE_TIME, Cfg.Service, ticketId64)
		if err != nil {
			return errors.New("Error (TU-040107): " + fmt.Sprintf("%s\n", err))
		}
	}

	return nil
}

// ----- получить ответную пользовательскую квитанцию -----
// 02
func Get_Ticket_Reply(ticketId string) (Models.SignedContentOptions, error) {

	var TicketReplyData Models.SignedContentOptions

	url := Cfg.ServiceURL + "/api/tickets/createreply/" + ticketId

	data, err := Web.SendPostRequest(url, "", "application/json", Cfg.TOKEN)
	if err != nil {
		return TicketReplyData, errors.New("Error (TU-040201): " + fmt.Sprintf("%s\n", err))
	}
	//fmt.Println(data)

	if err = json.Unmarshal([]byte(data), &TicketReplyData); err != nil {
		return TicketReplyData, errors.New("Error (TU-040202): " + fmt.Sprintf("%s\n", err))
	}
	//fmt.Printf("***** TicketReplyData: %s\n", TicketReplyData)

	return TicketReplyData, nil
}
