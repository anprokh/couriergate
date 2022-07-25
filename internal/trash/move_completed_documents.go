package trash

import (
	Cfg "couriergate/configs"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ----- переместить отработанные файлы в архивный каталог -----
// 01
func MoveCompletedDocumentsToArchiveFolder() error {

	name := time.Now().Format("2006-01-02")
	subdirFullName := fmt.Sprintf("%s\\COURIER-%s", Cfg.ExPath, name)
	//fmt.Println(subdirFullName)

	// создаем каталог если он не существует
	if _, err := os.Stat(subdirFullName); err != nil {
		err := os.Mkdir(subdirFullName, 0755)
		if err != nil {
			return errors.New("Error (TU-010101): " + fmt.Sprintf("%s\n", err))
		}
	}

	// проверим успешность создания
	if _, err := os.Stat(subdirFullName); err != nil {
		return errors.New("Error (TU-010102): " + fmt.Sprintf("%s\n", err))
	}

	// ----- перемещаем файлы исходящих УПД в архивный каталог -----
	var fileNames []string

	files, err := ioutil.ReadDir(Cfg.ExPath)
	if err != nil {
		return errors.New("Error (TU-010103): " + fmt.Sprintf("%s\n", err))
	}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) == ".Document" {
			fileNames = append(fileNames, file.Name())
		}
	}

	for _, name := range fileNames {

		responseFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, name)
		//fmt.Println(responseFullName)

		// ----- определяем имя файла подписи -----
		signatureName := strings.Replace(name, ".xml.Document", ".xml.sgn", -1)
		signatureFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, signatureName)

		//fmt.Println(signatureName)
		//fmt.Println(signatureFullName)

		// если файл подписи не существует, значит что-то пошло не так...
		if _, err := os.Stat(signatureFullName); err != nil {
			continue
		}

		// ----- определяем имя исходящего файла -----
		fileName := strings.Replace(name, ".xml.Document", ".xml", -1)
		fileFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, fileName)
		//fmt.Println(fileName)
		//fmt.Println(fileFullName)

		// если файл не существует, значит что-то пошло не так...
		if _, err := os.Stat(fileFullName); err != nil {
			continue
		}

		// перемещаем исходящий файл
		destination := fmt.Sprintf("%s\\%s", subdirFullName, fileName)
		err := os.Rename(fileFullName, destination)
		if err != nil {
			return errors.New("Error (TU-010104): " + fmt.Sprintf("%s\n", err))
		}

		// перемещаем файл подписи
		destination = fmt.Sprintf("%s\\%s", subdirFullName, signatureName)
		err = os.Rename(signatureFullName, destination)
		if err != nil {
			return errors.New("Error (TU-010105): " + fmt.Sprintf("%s\n", err))
		}

		// перемещаем файл ответа
		destination = fmt.Sprintf("%s\\%s", subdirFullName, name)
		err = os.Rename(responseFullName, destination)
		if err != nil {
			return errors.New("Error (TU-010106): " + fmt.Sprintf("%s\n", err))
		}

	}

	return nil
}
