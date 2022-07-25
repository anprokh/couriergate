package models

// Config is the struct to hold configuration options
type Config struct {
	Service          string
	ServiceURL       string
	Account          string
	ApiKey           string
	Certificate      string
	Pin              string
	ConnectionString string
	IncomingApply    bool
}

type TrackingOptions struct {
	UserName        string
	Password        string
	RelationId      string
	TrackingId      string
	SystemEDI       string
	DocumentContent string
}

type TicketIndexOptions struct {
	Id           int64  `json:id`
	DocumentId   int64  `json:documentId`
	Uuid         string `json:uuid`
	ParentUuid   string `json:parentUuid`
	Direction    int    `json:direction`
	Type         int    `json:type`
	ExtendedType int    `json:extendedType`
	Date         string `json:date`
}

type SignedContentOptions struct {
	Content   string `json:content`
	FileName  string `json:filename`
	MimeType  string `json:mimeType`
	Signature string `json:signature`
}

type DocumentDetailsOptions struct {
	Id               int                  `json:id`
	DocumentTypeCode string               `json:documentTypeCode`
	Description      string               `json:description`
	Status           int                  `json:status`
	Number           string               `json:number`
	SellerCode       string               `json:sellerCode`
	BuyerCode        string               `json:buyerCode`
	Filename         string               `json:filename`
	Content          SignedContentOptions `json:content`
	SenderName       string               `json:senderName`
	TotalSum         float64              `json:totalSum`
	Date             string               `json:date`
}

type DocumentCardOptions struct {
	Number    string               `json:number`
	Date      string               `json:date`
	TypeCode  string               `json:typeCode`
	Content   SignedContentOptions `json:content`
	Signature SignatureOptions     `json:signature`
}

type SignatureOptions struct {
	Uuid       string `json:uuid`
	Content    string `json:content`
	SignerCode string `json:signerCode`
}

type TicketDetailsOptions struct {
	Id           int64                `json:id`
	DocumentId   int64                `json:documentId`
	Uuid         string               `json:uuid`
	ParentUuid   string               `json:parentUuid`
	Direction    int                  `json:direction`
	Type         int                  `json:type`
	ExtendedType int                  `json:extendedType`
	Date         string               `json:date`
	Content      SignedContentOptions `json:content`
	Signature    []SignatureOptions   `json:signature`
	NeedReply    bool                 `json:needReply`
}

type DocumentEventOptions struct {
	Id                   int64  `json:id`
	DocumentId           int64  `json:documentId`
	TicketId             int64  `json:ticketId`
	TicketType           int    `json:ticketType`
	ExtendedTicketType   int    `json:extendedTicketType`
	EventType            int    `json:eventType`
	Date                 string `json:date`
	DocumentRelationType int    `json:documentRelationType`
}

type LogonResponseOptions struct {
	Token string `json:token`
}
