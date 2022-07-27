package signatures

import (
	"bytes"
	Cfg "couriergate/configs"
	DB "couriergate/internal/db"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ----- создать файлы подписей ответных квитанций -----
// 01
func CreateTicketReplySignatures() error {

	var fileNames []string

	files, err := ioutil.ReadDir(Cfg.ExPath)
	if err != nil {
		return errors.New("Error (SU-020101): " + fmt.Sprintf("%s\n", err))
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

		// доп.проверка необходимости подписать квитанцию, по наличию записи в [Document.Out]
		//		allowedToSign, err := TicketAllowedToSign(ticketId)
		//		if err != nil {
		//			return errors.New("Error (SU-020102): " + fmt.Sprintf("%s\n", err))
		//		}
		//		if !allowedToSign {
		//			continue
		//		}

		// ----- определяем имя файла подписи -----
		signatureName := fmt.Sprintf("%s%s", name, ".sgn")
		signatureFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, signatureName)

		fmt.Println(signatureFullName)

		// пропускаем файлы для которых найденa подпись
		if _, err := os.Stat(signatureFullName); err == nil {
			fmt.Println("Ticket signature file definitely exists!")
			continue
		}

		// запускаем приложение cryptcp, подписываем файл
		appFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, "cryptcp.x64.exe")

		//		cmnd := exec.Command(appFullName, "-showstat")
		cmd := exec.Command(appFullName, "-sign", "-detached", "-dn", Cfg.Certificate, "-pin", Cfg.Pin, filepath, signatureFullName)

		var buf bytes.Buffer
		cmd.Stdout = &buf
		err = cmd.Start()
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
		err = cmd.Wait()

		fmt.Printf("Command finished with error: %v\n", err)

		i = strings.Index(name, "PDPOL")
		if i == -1 {
			// пропускаем квитанции не ПДП
			continue
		}

		// фиксируем время подписания файла
		if _, err := os.Stat(signatureFullName); err == nil {
			fmt.Println("Signature file successfully created ...")

			// время подписи
			DATE_TIME := time.Now().String()
			DATE_TIME = fmt.Sprintf("%sT%s", DATE_TIME[:10], DATE_TIME[11:23])

			// преобразовываем id в int64
			ticketId64, err := strconv.ParseInt(ticketId, 10, 64)
			if err != nil {
				return errors.New("Error (SU-020102): " + fmt.Sprintf("%s\n", err))
			}

			// ----- записываем в sql время подписи файла -----
			_, err = DB.DB_COURIER.Exec("UPDATE [Document.Out] SET TicketReplySigned = CAST(@p1 AS datetime) WHERE (Service = @p2) AND (TicketID = @p3)", DATE_TIME, Cfg.Service, ticketId64)
			if err != nil {
				return errors.New("Error (SU-020103): " + fmt.Sprintf("%s\n", err))
			}
		}
	}

	return nil
}

// проверяем необходимость подписать квитанцию
// 02
func TicketAllowedToSign(ticketId string) (bool, error) {

	// преобразовываем id в int64
	ticketId64, err := strconv.ParseInt(ticketId, 10, 64)
	if err != nil {
		return false, errors.New("Error (SU-020201): SQL request failure, " + fmt.Sprintf("%s\n", err))
	}

	// ПРОВЕРЯЕМ ПРОСТО, ЕСТЬ ЛИ В [Document.Out] ВООБЩЕ ЗАПИСЬ С НУЖНЫМ ticketId
	rows, err := DB.DB_COURIER.Query("SELECT ISNULL (TicketReplySigned, '') FROM [Document.Out] (NOLOCK) WHERE (Service = @p1) AND ((TicketID = @p2) OR (TicketUvutochID = @p2))", Cfg.Service, ticketId64)
	if err != nil {
		return false, errors.New("Error (SU-020202): SQL request failure, " + fmt.Sprintf("%s\n", err))
	}
	defer rows.Close()

	// ----- выборка и обработка данных запроса -----
	var TicketReplySigned string

	for rows.Next() {

		err := rows.Scan(&TicketReplySigned)
		if err != nil {
			return false, errors.New("Error (SU-020203): " + fmt.Sprintf("%s\n", err))
		}
		return true, nil
	}

	return false, nil
}
