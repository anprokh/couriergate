package export

import (
	Cfg "couriergate/configs"
	DB "couriergate/internal/db"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"os"
)

// ----- создать файлы исходящих документов для подписи -----
// 01
func CreateFilesForSignature() error {

	// ----- формируем текст запроса -----
	//requestText := fmt.Sprintf("SELECT FileContent, ISNULL (FileID, '') FROM [Document.Out] (NOLOCK) WHERE (Service = '%s') AND (Account = '%s')", cfg.Service, cfg.Account)
	requestText := fmt.Sprintf("SELECT FileContent, ISNULL (FileID, '') FROM [Document.Out] (NOLOCK) WHERE (FileCreated = 0) AND (Service = '%s') AND (Account = '%s')", Cfg.Service, Cfg.Account)
	//fmt.Println(requestText)

	// выполнение запроса
	rows, err := DB.DB_COURIER.Query(requestText)
	if err != nil {
		return errors.New("Error (EU-010101): SQL request failure, " + fmt.Sprintf("%s\n", err))
	}
	defer rows.Close()

	// ----- выборка и обработка данных запроса -----
	var FileContent, FileID string

	for rows.Next() {

		err := rows.Scan(&FileContent, &FileID)
		if err != nil {
			return errors.New("Error (EU-010102): " + fmt.Sprintf("%s\n", err))
		}
		//fmt.Printf("%s : %s\n", FileID, FileContent)
		//fmt.Printf("%s : %s\n", FileID, "")

		// ----- записываем xml-файл документа -----
		fullFileName := fmt.Sprintf("%s\\%s.xml", Cfg.ExPath, FileID)
		fmt.Println(fullFileName)

		//err = ioutil.WriteFile(fullFileName, []byte("Hi\n"), 0666)
		//f, err := os.OpenFile(fullFileName, os.O_CREATE | os.O_RDWR | os.O_APPEND, 0777)

		f, err := os.Create(fullFileName)
		if err != nil {
			return errors.New("Error (EU-010103): " + fmt.Sprintf("%s\n", err))
		}
		defer f.Close()

		// кодируем в Windows1251
		enc := charmap.Windows1251.NewEncoder()
		fileContent1251, _ := enc.String(FileContent)

		_, err = f.Write([]byte(fileContent1251))
		if err != nil {
			return errors.New("Error (EU-010104): " + fmt.Sprintf("%s\n", err))
		}

		// ----- устанавливаем флаг создания файла -----
		_, err = DB.DB_COURIER.Exec("UPDATE [Document.Out] SET FileCreated = 1 WHERE (Service = @p1) AND (FileID = @p2)", Cfg.Service, FileID)
		if err != nil {
			return errors.New("Error (EU-010105): " + fmt.Sprintf("%s\n", err))
		}
	}

	return nil
}
