package accept

import (
	"encoding/json"
	"errors"
	"fmt"
)

type PersonOptions struct {
	LastName   string `json:"lastName"`
	FirstName  string `json:"firstName"`
	MiddleName string `json:"middleName"`
}

type SignerOptions struct {
	SignerType int           `json:"signerType"`
	INN        string        `json:"inn"`
	Name       string        `json:"name"`
	Person     PersonOptions `json:"person"`
	Title      string        `json:"title"`
}

type SignerAuthorityOptions struct {
	AuthorityScope int    `json:"authorityScope"`
	Condition      int    `json:"condition"`
	Authority      string `json:"authority"`
}

type DocumentSignerOptions struct {
	Signer          SignerOptions          `json:"signer"`
	SignerAuthority SignerAuthorityOptions `json:"signerAuthority"`
}

type UniversalDocumentAcceptInfoOptions struct {
	DealSummary    string                  `json:"dealSummary"`
	DocumentSigner []DocumentSignerOptions `json:"documentSigner"`
}

// ----- заполнить структуру UniversalDocumentAcceptInfo -----
//
func GetUniversalDocumentAcceptInfo(certificateName string) (string, error) {

	var person PersonOptions
	var signer SignerOptions

	switch certificateName {
	case "Медведев":

		person.LastName = "Медведев"
		person.FirstName = "Вадим"
		person.MiddleName = "Витальевич"

		signer.SignerType = 1
		signer.INN = "7825444144"
		signer.Name = "ООО \"ТДСЗ\""
		signer.Person = person
		signer.Title = "Директор по закупкам"

	case "Позднеева":

		person.LastName = "Позднеева"
		person.FirstName = "Елена"
		person.MiddleName = "Витальевна"

		signer.SignerType = 1
		signer.INN = "7825444144"
		signer.Name = "ООО \"ТДСЗ\""
		signer.Person = person
		signer.Title = "Операционный директор"

	case "Красилова":

		person.LastName = "Красилова"
		person.FirstName = "Наталья"
		person.MiddleName = "Борисовна"

		signer.SignerType = 1
		signer.INN = "7825444144"
		signer.Name = "ООО \"ТДСЗ\""
		signer.Person = person
		signer.Title = "Бухгалтер"

	case "Анненкова":

		person.LastName = "Анненкова"
		person.FirstName = "Ольга"
		person.MiddleName = "Александровна"

		signer.SignerType = 1
		signer.INN = "7825444144"
		signer.Name = "ООО \"ТДСЗ\""
		signer.Person = person
		signer.Title = "Заместитель директора"

	case "Варзунов":

		person.LastName = "Варзунов"
		person.FirstName = "Андрей"
		person.MiddleName = "Викторович"

		signer.SignerType = 1
		signer.INN = "7825444144"
		signer.Name = "ООО \"ТДСЗ\""
		signer.Person = person
		signer.Title = "Директор по закупкам и маркетингу СТМ"

	default:
		return "", errors.New("Error (AU-020101): неизвестный сертификат " + fmt.Sprintf("%s\n", certificateName))
	}

	var signerAuthority SignerAuthorityOptions
	signerAuthority.AuthorityScope = 1
	signerAuthority.Condition = 5
	signerAuthority.Authority = "Должностные обязанности"

	var documentSigner DocumentSignerOptions
	documentSigner.Signer = signer
	documentSigner.SignerAuthority = signerAuthority

	var documentSignerArr []DocumentSignerOptions
	documentSignerArr = append(documentSignerArr, documentSigner)

	var universalDocumentAcceptInfo UniversalDocumentAcceptInfoOptions
	universalDocumentAcceptInfo.DealSummary = "Товары приняты"
	universalDocumentAcceptInfo.DocumentSigner = documentSignerArr

	data, err := json.Marshal(universalDocumentAcceptInfo)
	if err != nil {
		return "", errors.New("Error (AU-020102): " + fmt.Sprintf("%s\n", err))
	}

	return fmt.Sprintf("%s", data), nil
}
