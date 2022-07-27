package accept

import (
	"bytes"
	Cfg "couriergate/configs"
	Sign "couriergate/internal/signatures"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ----- создать файлы подписей Титул покупателя -----
// 01
func CreateAcceptTitleSignatures() error {

	var fileNames []string

	files, err := ioutil.ReadDir(Cfg.ExPath)
	if err != nil {
		return errors.New("Error (AU-030101): " + fmt.Sprintf("%s\n", err))
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

		filepath := fmt.Sprintf("%s\\%s", Cfg.ExPath, name)

		// ----- определяем имя файла подписи -----
		signatureName := fmt.Sprintf("%s%s", name, ".sgn")
		signatureFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, signatureName)

		fmt.Println(signatureFullName)

		// пропускаем файлы для которых найденa подпись
		if _, err := os.Stat(signatureFullName); err == nil {
			fmt.Println("File definitely exists!")
			continue
		}

		// ----- выделяем из имени файла ТП имя исходного файла -----
		var i = strings.Index(name, ".xml.")
		if i == -1 {
			// в имени файла нет ".xml.", что-то пошло не так
			continue
		}

		// ----- выделяем из имени файла id документа -----
		var j = strings.Index(name, ".AcceptTitle")
		if j == -1 {
			// в имени файла нет ".AcceptTitle", что-то пошло не так
			continue
		}
		documentID := name[(i + 5):j]

		certificateName, err := Sign.GetCertificateNameByDocumentID(documentID)
		if err != nil {
			fmt.Println(err)
			continue
		}

		rdn, err := Sign.GetRDNByCertificateName(certificateName)
		if err != nil {
			fmt.Println(err)
			continue
		}

		pin, err := Sign.GetPinByCertificateName(certificateName)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// запускаем приложение cryptcp, подписываем файл
		appFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, "cryptcp.x64.exe")

		//		cmnd := exec.Command(appFullName, "-showstat")
		cmd := exec.Command(appFullName, "-sign", "-detached", "-dn", rdn, "-pin", pin, filepath, signatureFullName)

		var buf bytes.Buffer
		cmd.Stdout = &buf
		err = cmd.Start()
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
		err = cmd.Wait()

		fmt.Printf("Command finished with error: %v\n", err)
		//fmt.Printf("Command finished with output: %v\n", buf.String())
	}

	return nil
}
