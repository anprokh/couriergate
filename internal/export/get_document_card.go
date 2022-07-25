package export

import (
	Models "couriergate/models"
	"fmt"
	"time"
)

// ----- заполнить структуру DocumentCard -----
// 01
func GetDocumentCardNew(fileName string, content64 string, signature64 string, TypeCode string) (Models.DocumentCardOptions, error) {

	var Content Models.SignedContentOptions
	Content.FileName = fileName
	Content.MimeType = "text/xml"
	Content.Content = content64

	var documentCard Models.DocumentCardOptions

	// Временна́я метка отправляемого пакета данных
	t := time.Now()
	DOCNO := fmt.Sprintf("%d", t.Unix())
	//fmt.Println(DOCNO)

	DATE_TIME := t.String()
	DATE_TIME = fmt.Sprintf("%sT%s.9632729+03:00", DATE_TIME[:10], DATE_TIME[11:19])
	//fmt.Printf("DATE_TIME: >%s<\n", DATE_TIME)

	documentCard.Number = DOCNO
	documentCard.Date = DATE_TIME
	documentCard.TypeCode = TypeCode
	documentCard.Content = Content

	if len(signature64) > 0 {
		var Signature Models.SignatureOptions
		Signature.Content = signature64
		documentCard.Signature = Signature
	}

	return documentCard, nil
}

func GetDocumentCard(fileName string, content64 string, signature64 string, TypeCode string) (string, error) {

	// Временна́я метка отправляемого пакета данных
	t := time.Now()
	DOCNO := fmt.Sprintf("%d", t.Unix())
	//fmt.Println(DOCNO)

	DATE_TIME := t.String()
	DATE_TIME = fmt.Sprintf("%sT%s.9632729+03:00", DATE_TIME[:10], DATE_TIME[11:19])
	//fmt.Printf("DATE_TIME: >%s<\n", DATE_TIME)

	documentCard := "<DocumentCard>\n" +
		"<Number>" + DOCNO + "</Number>\n" +
		"<Date>" + DATE_TIME + "</Date>\n" +
		"<TypeCode>" + TypeCode + "</TypeCode>\n" +
		"<Content>\n" +
		"  <Filename>" + fileName + "</Filename>\n" +
		"  <MimeType>text/xml</MimeType>\n" +
		"  <Content>" + content64 + "</Content>\n" +
		"</Content>\n"

	if len(signature64) > 0 {
		documentCard = documentCard +
			"<Signature>\n" +
			"  <Content>" + signature64 + "</Content>\n" +
			"</Signature>\n"
	}

	documentCard = documentCard +
		"</DocumentCard>"

	//fmt.Println("----------------------------------------------------------------------------------------------------------------------")
	//fmt.Println(documentCard)

	return documentCard, nil
}
