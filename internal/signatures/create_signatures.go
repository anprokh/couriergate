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
	"strings"
	"time"
)

// ----- создать файлы подписей (подписать документы) -----
// 01
func CreateSignatures() error {

	var fileNames []string

	files, err := ioutil.ReadDir(Cfg.ExPath)
	if err != nil {
		return errors.New("Error (SU-010101): " + fmt.Sprintf("%s\n", err))
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

		filepath := fmt.Sprintf("%s\\%s", Cfg.ExPath, name)
		fmt.Println(filepath)

		// ----- определяем имя файла подписи -----
		signatureName := fmt.Sprintf("%s%s", name, ".sgn")
		signatureFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, signatureName)

		//fmt.Println(signatureName)
		fmt.Println(signatureFullName)

		// пропускаем файлы для которых найденa подпись
		if _, err := os.Stat(signatureFullName); err == nil {
			fmt.Println("File definitely exists!")
			continue
		}

		// запускаем приложение cryptcp, подписываем файл
		appFullName := fmt.Sprintf("%s\\%s", Cfg.ExPath, "cryptcp.x64.exe")

		//		cmnd := exec.Command(appFullName, "-showstat")
		cmd := exec.Command(appFullName, "-sign", "-detached", "-dn", Cfg.Certificate, "-pin", Cfg.Pin, filepath, signatureFullName)

		var buf bytes.Buffer
		cmd.Stdout = &buf
		err := cmd.Start()
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
		err = cmd.Wait()

		fmt.Printf("Command finished with error: %v\n", err)
		//fmt.Printf("Command finished with output: %v\n", buf.String())

		i := strings.Index(name, ".xml")
		FileID := name[:i]

		// фиксируем время подписания файла
		if _, err := os.Stat(signatureFullName); err == nil {
			fmt.Println("Signature file successfully created ...")

			// время подписи
			DATE_TIME := time.Now().String()
			DATE_TIME = fmt.Sprintf("%sT%s", DATE_TIME[:10], DATE_TIME[11:23])

			// ----- записываем в sql время подписи файла -----
			_, err = DB.DB_COURIER.Exec("UPDATE [Document.Out] SET FileSigned = CAST(@p1 AS datetime) WHERE (Service = @p2) AND (FileID = @p3)", DATE_TIME, Cfg.Service, FileID)
			if err != nil {
				return errors.New("Error (SU-010102): " + fmt.Sprintf("%s\n", err))
			}
		}

	}

	return nil
}
