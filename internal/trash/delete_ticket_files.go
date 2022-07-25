package trash

import (
	Cfg "couriergate/configs"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// ----- удалить ненужные отработанные файлы квитанций -----
// 01
func DeleteTicketFiles() error {

	var fileNames []string

	files, err := ioutil.ReadDir(Cfg.ExPath)
	if err != nil {
		return errors.New("Error (TU-020101): " + fmt.Sprintf("%s\n", err))
	}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) == ".TicketReplyNotice" {
			fileNames = append(fileNames, file.Name())
		}
	}

	for _, noticeName := range fileNames {

		noticeFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, noticeName)
		//fmt.Println(noticeFullName)

		// ----- определяем имя файла подписи -----
		signatureName := strings.Replace(noticeName, ".TicketReplyNotice", ".TicketReply.sgn", -1)
		signatureFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, signatureName)
		//fmt.Println(signatureFullName)

		// если файл подписи не существует, значит что-то пошло не так...
		if _, err := os.Stat(signatureFullName); err != nil {
			continue
		}

		// ----- определяем имя файла ответной квитанции -----
		replyName := strings.Replace(noticeName, ".TicketReplyNotice", ".TicketReply", -1)
		replyFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, replyName)
		//fmt.Println(replyFullName)

		// если файл ответной квитанции не существует, значит что-то пошло не так...
		if _, err := os.Stat(replyFullName); err != nil {
			continue
		}

		// ----- определяем имя файла квитанции -----
		ticketName := strings.Replace(noticeName, ".TicketReplyNotice", ".Ticket", -1)
		ticketFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, ticketName)
		//fmt.Println(ticketFullName)

		// если файл квитанции не существует, значит что-то пошло не так...
		if _, err := os.Stat(ticketFullName); err != nil {
			continue
		}

		// ----- удаляем найденные файлы -----
		err := os.Remove(ticketFullName)
		if err != nil {
			continue
		}

		err = os.Remove(replyFullName)
		if err != nil {
			continue
		}

		err = os.Remove(signatureFullName)
		if err != nil {
			continue
		}

		err = os.Remove(noticeFullName)
		if err != nil {
			continue
		}

	}

	return nil
}
