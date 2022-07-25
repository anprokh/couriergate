package export

import (
	"bufio"
	Cfg "couriergate/configs"
	DB "couriergate/internal/db"
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
	"unicode/utf8"
)

// ----- экспорт подписанных документов из файлов в систему Courier
// 01
func Export_Signed_Documents_FromFiles() error {

	var fileNames []string

	files, err := ioutil.ReadDir(Cfg.ExPath)
	if err != nil {
		return errors.New("Error (EU-010102): " + fmt.Sprintf("%s\n", err))
	}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) == ".xml" {
			fileNames = append(fileNames, file.Name())
		}
	}

	for _, name := range fileNames {

		fileFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, name)
		fmt.Println(fileFullName)

		// ----- пропускаем файлы для которых найден файл ответа -----
		responseName := fmt.Sprintf("%s%s", name, ".Document")
		responseFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, responseName)

		if _, err := os.Stat(responseFullName); err == nil {
			fmt.Println("====================> Response file definitely exists!")
			continue
		}

		// ----- признак файла черновика документа, отправка без подписи -----
		var isDraft bool

		draftName := fmt.Sprintf("%s%s", name, ".Draft")
		draftFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, draftName)

		if _, err := os.Stat(draftFullName); err == nil {
			fmt.Println("Draft file definitely exists! ... send a document as a draft")
			isDraft = true
		}
		//fmt.Printf("isDraft: %t\n", isDraft)

		// ----- определяем имя файла подписи -----
		signatureName := fmt.Sprintf("%s%s", name, ".sgn")
		signatureFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, signatureName)

		//fmt.Println(signatureName)
		fmt.Println(signatureFullName)

		// пропускаем файлы для которых не найдена подпись
		if _, err := os.Stat(signatureFullName); os.IsNotExist(err) {
			if !isDraft {
				fmt.Println("Signature File definitely does not exist.")
				continue
			}
		}

		// ----- читаем содержимое файла -----
		file1, err := os.Open(fileFullName)
		if err != nil {
			continue
		}
		defer file1.Close()

		fileContent, err := ioutil.ReadAll(file1)
		if err != nil {
			continue
		}
		//fmt.Println(string(fileContent))

		// ----- кодируем содержимое в base64 -----
		data64 := base64.StdEncoding.EncodeToString(fileContent)
		//fmt.Println("data64: " + data64)

		// ---------- Определяем функцию документа из текста файла ----------

		// переместимся в начало файла документа
		_, err = file1.Seek(0, 0)
		if err != nil {
			continue
		}

		TypeCode := "UPD_SF_DOP"
		var nextLine string
		var i, j int

		r := charmap.Windows1251.NewDecoder().Reader(file1)
		scanner := bufio.NewScanner(r)

		for scanner.Scan() { // internally, it advances token based on sperator
			nextLine = scanner.Text()
			//			fmt.Println(nextLine)  // token in unicode-char

			i = strings.Index(nextLine, "<Документ")
			if i >= 0 {
				j = strings.Index(nextLine, "Функция=\"ДОП\"")
				if j >= 0 {
					TypeCode = "UPD_DOP"
					break
				}
			}
		}
		//		fmt.Println(TypeCode)
		// ------------------------------------------------------------------

		var file2 *os.File
		var signatureContent []byte
		var signature64 string

		// ----- для черновиков не обрабатываем файл подписи -----
		if isDraft {
			goto A
		}

		// ----- читаем содержимое подписи -----
		file2, err = os.Open(signatureFullName)
		if err != nil {
			continue
		}
		defer file2.Close()

		signatureContent, err = ioutil.ReadAll(file2)
		if err != nil {
			continue
		}
		signature64 = fmt.Sprintf("%s", signatureContent)
		//fmt.Printf("signature64:\n%s", signature64)

		// ----- сформируем DocumentCard -----
	A:
		documentCard, err := GetDocumentCardNew(name, data64, signature64, TypeCode)
		if err != nil {
			continue
		}
		//fmt.Printf("documentCard:\n%s\n", documentCard)

		// ----- определяем идентификатор получателя из имени файла -----
		var x1 = strings.Split(name, "_")
		destinationID := x1[2]
		//fmt.Printf("destinationID:\n%s\n", destinationID)

		//destinationID := "2BK-3437014240-49459"
		//destinationID = "2BM-7703270067-2013012406531229784200000000000"  // Ашан

		data, err := SendDocToCourier(destinationID, documentCard)
		if err != nil {
			err1 := CreateResponseFile(name, "ERROR", err.Error())
			if err1 != nil {
				return errors.New("Error (EU-010103): " + fmt.Sprintf("%s\n", err))
			}
			continue
		}

		fmt.Printf("%s\n", data)
		fmt.Printf("%s\n", err)

		k := strings.Index(name, ".xml")
		FileID := name[:k]

		// ----- записываем файл структуры Document при успешном выполнении запроса -----
		var Document Models.DocumentDetailsOptions
		if err = json.Unmarshal([]byte(data), &Document); err != nil {
			err := CreateResponseFile(name, "ERROR", data)
			if err != nil {
				return errors.New("Error (EU-010104): " + fmt.Sprintf("%s\n", err))
			}
		} else {
			err := CreateResponseFile(name, "Document", data)
			if err != nil {
				return errors.New("Error (EU-010105): " + fmt.Sprintf("%s\n", err))
			}

			// время передачи файла
			DATE_TIME := time.Now().String()
			DATE_TIME = fmt.Sprintf("%sT%s", DATE_TIME[:10], DATE_TIME[11:23])

			// ----- записываем в sql время передачи файла и Id документа -----
			_, err = DB.DB_COURIER.Exec("UPDATE [Document.Out] SET FileDelivered = CAST(@p1 AS datetime), DocumentID = @p2 WHERE (Service = @p3) AND (FileID = @p4)", DATE_TIME, Document.Id, Cfg.Service, FileID)
			if err != nil {
				return errors.New("Error (EU-010106): " + fmt.Sprintf("%s\n", err))
			}
		}

	}

	return nil
}

// ----- создать файл ответа (записать ответ в файл) -----
// 02
func CreateResponseFile(fileName string, extension string, fileContent string) error {

	fileName = fmt.Sprintf("%s.%s", fileName, extension)
	fileFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, fileName)

	f, err := os.Create(fileFullName)
	if err != nil {
		return errors.New("Error (EU-010201): " + fmt.Sprintf("%s\n", err))
	}
	defer f.Close()

	// кодируем в Windows1251
	enc := charmap.Windows1251.NewEncoder()

	fileContent1251, err := enc.String(fileContent)
	if err != nil {
		//fmt.Println("ERROR (EU-010202): " + err.Error())

		// вероятно встретилась некорректная руна, пытаемся преобразовать посимвольно
		tmp := ""
		for i := 0; i < len(fileContent); {
			r, size := utf8.DecodeRuneInString(fileContent[i:])
			i += size

			runeStr := fmt.Sprintf("%c", r)
			//fmt.Println(runeStr)

			// проверим, корректна ли очередная руна
			_, err = enc.String(runeStr)
			if err != nil {
				continue
			}
			tmp = tmp + runeStr
		}
		//fmt.Printf("tmp: %s\n", tmp)
		fileContent1251, _ = enc.String(tmp)
	}

	_, err = f.Write([]byte(fileContent1251))
	if err != nil {
		return errors.New("Error (EU-010202): " + fmt.Sprintf("%s\n", err))
	}

	return nil
}
