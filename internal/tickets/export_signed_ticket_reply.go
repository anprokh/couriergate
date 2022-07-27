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
	"golang.org/x/text/encoding/charmap"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ----- экспорт подписанных квитанций из файлов в систему Courier
// 01
func Export_Signed_TicketReply_FromFiles() error {

	var fileNames []string

	files, err := ioutil.ReadDir(Cfg.ExPath)
	if err != nil {
		return errors.New("Error (TU-050101): " + fmt.Sprintf("%s\n", err))
	}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) == ".TicketReply" {
			fileNames = append(fileNames, file.Name())
		}
	}

	for _, name := range fileNames {

		filepath := fmt.Sprintf("%s\\%s", Cfg.ExPath, name)
		fmt.Println(filepath)

		// ----- определяем имя файла подписи -----
		signatureName := fmt.Sprintf("%s%s", name, ".sgn")
		signatureFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, signatureName)

		fmt.Println(signatureFullName)

		// пропускаем файлы для которых не найдена подпись
		if _, err := os.Stat(signatureFullName); os.IsNotExist(err) {
			fmt.Println("File definitely does not exist.")
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
		fmt.Printf("ticketId: %s\n", ticketId)

		// ----- определяем тип квитанции -----
		var ticketType string
		i = strings.Index(name, "PDPOL")
		if i > 0 {
			ticketType = "PDPOL"
		}
		i = strings.Index(name, "UVUTOCH")
		if i > 0 {
			ticketType = "UVUTOCH"
		}
		i = strings.Index(name, "PDOTPR")
		if i > 0 {
			ticketType = "PDOTPR"
		}
		fmt.Printf("ticketType: %s\n", ticketType)

		replyNoticeName := fmt.Sprintf("%sNotice", name)
		replyNoticeFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, replyNoticeName)

		// ----- пропускаем файлы для которых найден файл ответа подтверждения квитанции -----
		if _, err := os.Stat(replyNoticeFullName); err == nil {
			fmt.Println("====================> Reply Notice file definitely exists!")
			continue
		}

		// ----- читаем содержимое файла -----
		file1, err := os.Open(filepath)
		if err != nil {
			continue
		}
		defer file1.Close()

		fileContent, err := ioutil.ReadAll(file1)
		if err != nil {
			continue
		}

		// ----- кодируем содержимое в base64 -----
		data64 := base64.StdEncoding.EncodeToString(fileContent)

		// ----- читаем содержимое подписи -----
		file2, err := os.Open(signatureFullName)
		if err != nil {
			continue
		}
		defer file2.Close()

		signatureContent, err := ioutil.ReadAll(file2)
		if err != nil {
			continue
		}
		signature64 := fmt.Sprintf("%s", signatureContent)
		//fmt.Printf("signature64:\n%s", signature64)

		// ----- сформируем SignedContent -----
		var signedContent Models.SignedContentOptions
		signedContent.Content = data64
		signedContent.Signature = signature64

		receiveNoticeContent, err := SendTicketToCourier(ticketId, signedContent)
		if err != nil {
			fmt.Printf("Error (TU-050102): %s\n", err)
			continue
		}

		// сохраним ответ подтверждения квитанции в файл
		f, err := os.Create(replyNoticeFullName)
		if err != nil {
			return errors.New("Error (TU-050103): " + fmt.Sprintf("%s\n", err))
		}
		defer f.Close()

		// кодируем в Windows1251
		enc := charmap.Windows1251.NewEncoder()
		fileContent1251, _ := enc.String(receiveNoticeContent)

		_, err = f.Write([]byte(fileContent1251))
		if err != nil {
			return errors.New("Error (TU-050104): " + fmt.Sprintf("%s\n", err))
		}

		// время передачи ответной квитанции
		DATE_TIME := time.Now().String()
		DATE_TIME = fmt.Sprintf("%sT%s", DATE_TIME[:10], DATE_TIME[11:23])

		// ----- записываем в sql время получения ответной квитанции -----
		switch ticketType {
		case "PDPOL":
			_, err = DB.DB_COURIER.Exec("UPDATE [Document.Out] SET TicketReplyDelivered = CAST(@p1 AS datetime) WHERE (Service = @p2) AND (TicketID = @p3)", DATE_TIME, Cfg.Service, ticketId)
		case "UVUTOCH":
			_, err = DB.DB_COURIER.Exec("UPDATE [Document.Out] SET TicketUvutochReplyDelivered = CAST(@p1 AS datetime) WHERE (Service = @p2) AND (TicketUvutochID = @p3)", DATE_TIME, Cfg.Service, ticketId)
		default:
			continue
		}
		if err != nil {
			return errors.New("Error (TU-050105): " + fmt.Sprintf("%s\n", err))
		}
	}

	return nil
}

// ----- отправить ответную квитанцию в "Courier" EDM -----
// 04
func SendTicketToCourier(id string, SignedContent Models.SignedContentOptions) (string, error) {

	url := Cfg.ServiceURL + "/api/tickets/add/" + id
	fmt.Println(url)

	value, err := json.Marshal(SignedContent)
	//	value, err := json.MarshalIndent(SignedContent, "", " ")
	if err != nil {
		return "", errors.New("Error (TU-050201): " + fmt.Sprintf("%s\n", err))
	}

	data, err := Web.SendPostRequest(url, string(value), "application/json", Cfg.TOKEN)
	if err != nil {
		return "", errors.New("Error (TU-050202): " + fmt.Sprintf("%s\n", err))
	}

	return data, nil
}
