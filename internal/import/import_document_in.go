package Import

import (
	"bufio"
	Cfg "couriergate/configs"
	DB "couriergate/internal/db"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"golang.org/x/text/encoding/charmap"
	"strconv"
	"strings"
)

// ----- получить методом Document.Details данные по документу и записать в Document.In -----
// 01
func Import_Document_In(sourceID string) error {

	documentData, err := GetDocumentDetails(sourceID)
	if err != nil {
		return errors.New("Error (IU-030101): " + fmt.Sprintf("%s\n", err))
	}
	//fmt.Println(documentData)

	status := strconv.Itoa(documentData.Status)
	switch documentData.Status {
	case 13:
		status = "Delivered"
	case 14:
		status = "Received"
	default:
	}
	//fmt.Println(status)

	data64 := documentData.Content.Content

	// ----- декодируем содержимое документа из Base64 и приводим к win-1251 -----
	sDec, err := base64.StdEncoding.DecodeString(data64) // string -> []byte
	if err != nil {
		return errors.New("Error (IU-030102): " + fmt.Sprintf("%s\n", err))
	}
	//fmt.Printf("sDec: %s\n", sDec)

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

	// ---------- формируем GUID из имени файла (последние 36 символов) ----------
	filename := documentData.Filename
	i := strings.Index(filename, ".")
	if i > 0 {
		filename = filename[:i]
	}
	//fmt.Println(filename)

	// произвольный GUID на случай отсутствия его в имени файла
	file_id := uuid.NewV4()
	file_id_string := file_id.String()
	//fmt.Println(file_id_string)

	l := len(filename)
	if l >= 36 {
		tmpS := filename[(l - 36):]
		//fmt.Println(tmpS)

		// если удалось создать GUID из имени файла - используем его
		u1, err := uuid.FromString(tmpS)
		if err == nil {
			file_id_string = u1.String()
		}
		//fmt.Println(file_id_string)
	}

	// ----- записываем в sql информацию о полученном документе -----
	_, err = DB.DB_COURIER.Exec("INSERT INTO [Document.In] (Service, Account, DocumentID, DocumentTypeCode, Description, Status, Number, TotalSum, SellerCode, BuyerCode, "+
		"Filename, FileContent, SenderName, Date, FileID) "+
		"VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12, @p13, @p14, @p15)",
		Cfg.Service,
		Cfg.Account,
		documentData.Id,
		documentData.DocumentTypeCode,
		documentData.Description,
		status,
		documentData.Number,
		documentData.TotalSum,
		documentData.SellerCode,
		documentData.BuyerCode,
		documentData.Filename,
		fileContent,
		documentData.SenderName,
		documentData.Date,
		file_id_string)

	if err != nil {
		return errors.New("Error (IU-030103): " + fmt.Sprintf("%s\n", err))
	}

	return nil
}
