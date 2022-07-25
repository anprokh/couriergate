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

// ----- удалить ненужные отработанные файлы исходящих документов -----
// 01
func DeleteDocumentFiles() error {

	var fileNames []string

	files, err := ioutil.ReadDir(Cfg.ExPath)
	if err != nil {
		return errors.New("Error (TU-030101): " + fmt.Sprintf("%s\n", err))
	}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) == ".Document" {
			fileNames = append(fileNames, file.Name())
		}
	}

	for _, documentName := range fileNames {
		//fmt.Println(documentName)

		signatureName := strings.Replace(documentName, ".xml.Document", ".xml.sgn", -1)
		//fmt.Println(signatureName)

		fileName := strings.Replace(documentName, ".xml.Document", ".xml", -1)
		//fmt.Println(fileName)

		documentFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, documentName)
		//fmt.Println(documentFullName)

		signatureFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, signatureName)
		//fmt.Println(signatureFullName)

		fileFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, fileName)
		//fmt.Println(fileFullName)

		// пропускаем файлы для которых не найдена подпись
		if _, err := os.Stat(signatureFullName); os.IsNotExist(err) {
			fmt.Println("Signature File definitely does not exist.")
			continue
		}

		// пропускаем файлы для которых не найден исходный документ
		if _, err := os.Stat(fileFullName); os.IsNotExist(err) {
			fmt.Println("Original File definitely does not exist.")
			continue
		}

		newDocumentFullName := fmt.Sprintf("%s\\Archive\\%s", Cfg.ExPath, documentName)
		//fmt.Println(newDocumentFullName)

		newSignatureFullName := fmt.Sprintf("%s\\Archive\\%s", Cfg.ExPath, signatureName)
		//fmt.Println(newSignatureFullName)

		newFileFullName := fmt.Sprintf("%s\\Archive\\%s", Cfg.ExPath, fileName)
		//fmt.Println(newFileFullName)

		err := os.Rename(documentFullName, newDocumentFullName)
		if err != nil {
			return errors.New("Error (TU-030102): " + fmt.Sprintf("%s\n", err))
		}
		err = os.Rename(signatureFullName, newSignatureFullName)
		if err != nil {
			return errors.New("Error (TU-030103): " + fmt.Sprintf("%s\n", err))
		}
		err = os.Rename(fileFullName, newFileFullName)
		if err != nil {
			return errors.New("Error (TU-030104): " + fmt.Sprintf("%s\n", err))
		}

	}

	return nil
}
