package clarification

import (
	"bytes"
	Cfg "couriergate/configs"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// ----- создать файлы подписей УОУ (подписать УОУ) -----
// 01
func CreateClarificationSignatures() error {

	var fileNames []string

	files, err := ioutil.ReadDir(Cfg.ExPath)
	if err != nil {
		return errors.New("Error (CU-020101): " + fmt.Sprintf("%s\n", err))
	}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) == ".Clarification" {
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
	}

	return nil
}
