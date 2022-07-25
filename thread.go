package main

import (
	Cfg "couriergate/configs"
	Accept "couriergate/internal/accept"
	Auth "couriergate/internal/auth"
	Clarification "couriergate/internal/clarification"
	DB "couriergate/internal/db"
	Events "couriergate/internal/events"
	Export "couriergate/internal/export"
	Sign "couriergate/internal/signatures"
	Tickets "couriergate/internal/tickets"
	Trash "couriergate/internal/trash"
	"fmt"
	"github.com/fatih/color"
	"time"
)

func (p *program) ConnectDB() error {
	db, err := DB.GetDBConnection(p.cfg.ConnectionString)
	if err != nil {
		return err
	}
	DB.DB_COURIER = db
	return nil
}

func (p *program) GetAuthToken() error {
	token, err := Auth.GetAuthTokenByApiKey()
	if err != nil {
		return err
	}
	Cfg.TOKEN = token
	return nil
}

func (p *program) setEnv() {
	Cfg.Service = p.cfg.Service
	Cfg.ServiceURL = p.cfg.ServiceURL
	Cfg.Account = p.cfg.Account
	Cfg.ApiKey = p.cfg.ApiKey
	Cfg.Certificate = p.cfg.Certificate
	Cfg.Pin = p.cfg.Pin
}

func (p *program) exportSignedDocuments() {

	err := Export.CreateFilesForSignature() // 01
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s at %s\n", color.RedString("[error]"), color.CyanString(err.Error()), time.Now().Format("2006-01-02 15:04:05"))
	}

	err = Sign.CreateSignatures() // 02
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s at %s\n", color.RedString("[error]"), color.CyanString(err.Error()), time.Now().Format("2006-01-02 15:04:05"))
	}

	err = Export.Export_Signed_Documents_FromFiles() // 03 ТАКЖЕ ОТПРАВЛЯЕМ НЕПОДПИСАННЫЕ ЧЕРНОВИКИ !!
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s at %s\n", color.RedString("[error]"), color.CyanString(err.Error()), time.Now().Format("2006-01-02 15:04:05"))
	}

}

func (p *program) exportSignedTicketReply() {

	err := Tickets.CreateTicketFiles() // 05
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s at %s\n", color.RedString("[error]"), color.CyanString(err.Error()), time.Now().Format("2006-01-02 15:04:05"))
	}

	err = Tickets.CreateTicketReplyFiles() // 06
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s at %s\n", color.RedString("[error]"), color.CyanString(err.Error()), time.Now().Format("2006-01-02 15:04:05"))
	}

	err = Sign.CreateTicketReplySignatures() // 07
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s at %s\n", color.RedString("[error]"), color.CyanString(err.Error()), time.Now().Format("2006-01-02 15:04:05"))
	}

	err = Tickets.Export_Signed_TicketReply_FromFiles() // 08
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s at %s\n", color.RedString("[error]"), color.CyanString(err.Error()), time.Now().Format("2006-01-02 15:04:05"))
	}

}

func (p *program) processingDocumentEvents() {

	err := Events.Get_Events_Index() // 09
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s at %s\n", color.RedString("[error]"), color.CyanString(err.Error()), time.Now().Format("2006-01-02 15:04:05"))
	}

	err = Events.Processed_Document_Events() // 10
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s at %s\n", color.RedString("[error]"), color.CyanString(err.Error()), time.Now().Format("2006-01-02 15:04:05"))
	}

	err = Events.Processed_Document_Out_Pdf() // 11
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s at %s\n", color.RedString("[error]"), color.CyanString(err.Error()), time.Now().Format("2006-01-02 15:04:05"))
	}

	err = Events.Processed_Document_Out_Signatures() // 12
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s at %s\n", color.RedString("[error]"), color.CyanString(err.Error()), time.Now().Format("2006-01-02 15:04:05"))
	}

	// пропускаем функционал входящих УПД для незадействованных фирм
	if p.cfg.IncomingApply == false {
		return
	}

	err = Events.Processed_Document_In_Pdf() // 13
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s at %s\n", color.RedString("[error]"), color.CyanString(err.Error()), time.Now().Format("2006-01-02 15:04:05"))
	}

	err = Events.Processed_Document_In_Signatures() // 14
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s at %s\n", color.RedString("[error]"), color.CyanString(err.Error()), time.Now().Format("2006-01-02 15:04:05"))
	}

}

func (p *program) processingIncomingDocuments() {

	// ----- подписание проверенных входящих УПД -----
	err := Accept.CreateAcceptTitleFiles() // 15
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s at %s\n", color.RedString("[error]"), color.CyanString(err.Error()), time.Now().Format("2006-01-02 15:04:05"))
	}

	err = Accept.CreateAcceptTitleSignatures() // 16
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s at %s\n", color.RedString("[error]"), color.CyanString(err.Error()), time.Now().Format("2006-01-02 15:04:05"))
	}

	err = Accept.Export_Signed_AcceptTitles_FromFiles() // 17
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s at %s\n", color.RedString("[error]"), color.CyanString(err.Error()), time.Now().Format("2006-01-02 15:04:05"))
	}

	// **************************************** ОТКЛОНЕНИЕ ДОКУМЕНТОВ ****************************************

	err = Clarification.CreateClarificationFiles()
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s at %s\n", color.RedString("[error]"), color.CyanString(err.Error()), time.Now().Format("2006-01-02 15:04:05"))
	}

	err = Clarification.CreateClarificationSignatures()
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s at %s\n", color.RedString("[error]"), color.CyanString(err.Error()), time.Now().Format("2006-01-02 15:04:05"))
	}

	err = Clarification.Export_Signed_Clarifications_FromFiles()
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s at %s\n", color.RedString("[error]"), color.CyanString(err.Error()), time.Now().Format("2006-01-02 15:04:05"))
	}

}

func (p *program) moveCompletedDocuments() {

	// ----- удаление/архивирование отработанных файлов -----
	err := Trash.MoveCompletedDocumentsToArchiveFolder()
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s at %s\n", color.RedString("[error]"), color.CyanString(err.Error()), time.Now().Format("2006-01-02 15:04:05"))
	}

	err = Trash.DeleteTicketFiles()
	if err != nil {
		fmt.Fprintf(color.Output, "%s %s at %s\n", color.RedString("[error]"), color.CyanString(err.Error()), time.Now().Format("2006-01-02 15:04:05"))
	}

}
