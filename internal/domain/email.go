package domain

type Email struct {
	Sender        string
	SenderName    string
	Recipient     string
	RecipientName string
	Subject       string
	HtmlContent   string
}
