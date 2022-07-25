package accept

import (
	Models "couriergate/models"

	Cfg "couriergate/configs"

	Web "couriergate/internal/web"

	DB "couriergate/internal/db"

	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ----- экспорт подписанных Титулов покупателя из файлов в систему Courier -----
// 01
func Export_Signed_AcceptTitles_FromFiles() error {

	var fileNames []string

	files, err := ioutil.ReadDir(Cfg.ExPath)
	if err != nil {
		return errors.New("Error (AU-040101): " + fmt.Sprintf("%s\n", err))
	}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) == ".AcceptTitle" {
			fileNames = append(fileNames, file.Name())
		}
	}

	for _, name := range fileNames {

		fileFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, name)
		fmt.Println(fileFullName)

		// ----- определяем имя файла подписи -----
		signatureName := fmt.Sprintf("%s%s", name, ".sgn")
		signatureFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, signatureName)

		//fmt.Println(signatureName)
		fmt.Println(signatureFullName)

		// пропускаем файлы для которых не найдена подпись
		if _, err := os.Stat(signatureFullName); os.IsNotExist(err) {
			fmt.Println("Signature File definitely does not exist.")
			continue
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

		var file2 *os.File
		var signatureContent []byte
		var signature64 string

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

		// ----- выделяем из имени файла ТП имя исходного файла -----
		var i = strings.Index(name, ".xml.")
		if i == -1 {
			// в имени файла нет ".xml.", что-то пошло не так
			continue
		}
		originalFileName := name[:(i + 4)]

		// ----- выделяем из имени файла id документа -----
		var j = strings.Index(name, ".AcceptTitle")
		if j == -1 {
			// в имени файла нет ".AcceptTitle", что-то пошло не так
			continue
		}
		documentID := name[(i + 5):j]

		//		data, err := SendAcceptToCourier(cfg, documentID)
		//		if err != nil {
		//			fmt.Printf("Error (AU-040102): %s\n", err)
		//			continue
		//		}
		//		fmt.Printf("data: %s\n", data)
		//continue

		// ----- сформируем SignedContent -----
		var signedContent Models.SignedContentOptions
		signedContent.Content = data64
		signedContent.Signature = signature64
		signedContent.FileName = originalFileName

		_, err = SendAcceptTitleToCourier(documentID, signedContent)
		if err != nil {
			fmt.Printf("Error (AU-040102): %s\n", err)
			continue
		}
		//fmt.Printf("data: %s\n", data)

		// ----- регистрируем время экспорта титула -----
		DATE_TIME := time.Now().String()
		DATE_TIME = fmt.Sprintf("%sT%s", DATE_TIME[:10], DATE_TIME[11:23])

		_, err = DB.DB_COURIER.Exec("UPDATE [Document.In] SET ActionProcessed = CAST(@p1 AS datetime) WHERE Service = @p2 AND DocumentID = @p3", DATE_TIME, Cfg.Service, documentID)
		if err != nil {
			return errors.New("Error (AU-040103): " + fmt.Sprintf("%s\n", err))
		}

		// ----- удаляем файлы отправленного ТП -----
		file1.Close()
		err = os.Remove(fileFullName)
		if err != nil {
			fmt.Println(err)
		}

		file2.Close()
		err = os.Remove(signatureFullName)
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

// ----- отправить подписанную квитанцию в "Courier" EDM -----
// 02
func SendAcceptTitleToCourier(documentID string, SignedContent Models.SignedContentOptions) (string, error) {

	//fmt.Println("*****************************************************************************************************************************")
	url := Cfg.ServiceURL + "/api/document/accept/" + documentID
	fmt.Println(url)

	value, err := json.Marshal(SignedContent)
	if err != nil {
		return "", errors.New("Error (AU-040201): " + fmt.Sprintf("%s\n", err))
	}
	//fmt.Println(string(value))

	data, err := Web.SendPostRequest(url, string(value), "application/json", Cfg.TOKEN)
	if err != nil {
		return "", errors.New("Error (AU-040202): " + fmt.Sprintf("%s\n", err))
	}
	//fmt.Println(data)

	return data, nil
}

// ----- принять документ (когда подписание в 2 этапа) -----
// 3
func SendAcceptToCourier(documentID string) (string, error) {

	url := Cfg.ServiceURL + "/api/document/accept/" + documentID
	fmt.Println(url)

	data, err := Web.SendPostRequest(url, "", "application/json", Cfg.TOKEN)
	if err != nil {
		return "", errors.New("Error (AU-040301): " + fmt.Sprintf("%s\n", err))
	}
	fmt.Println(data)

	return "", nil
}
