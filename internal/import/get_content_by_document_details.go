package Import

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"strings"
)

// ----- получить Содержимое документа используя метод Document.Details -----
// 01
func GetContentByDocumentDetails(sourceID string) (string, error) {

	documentData, err := GetDocumentDetails(sourceID)
	if err != nil {
		return "", errors.New("Error (IU-020101): " + fmt.Sprintf("%s\n", err))
	}

	data64 := documentData.Content.Content

	// ----- декодируем содержимое документа из Base64 и приводим к win-1251 -----
	sDec, err := base64.StdEncoding.DecodeString(data64) // string -> []byte
	if err != nil {
		return "", errors.New("Error (IU-020102): " + fmt.Sprintf("%s\n", err))
	}

	// строка в win-1251
	sr := strings.NewReader(string(sDec))

	r := charmap.Windows1251.NewDecoder().Reader(sr)
	scanner := bufio.NewScanner(r)

	var b1 strings.Builder
	for scanner.Scan() {
		b1.WriteString(scanner.Text())
	}
	fileContent := b1.String()
	//fmt.Printf("fileContent: %s\n", fileContent)

	return fileContent, nil
}
