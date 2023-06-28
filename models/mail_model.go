package models

type Email struct {
	MessageID               string `json:"mensajeId"`
	Date                    string `json:"date"`
	From                    string `json:"from"`
	To                      string `json:"to"`
	Sent                    string `json:"sent"`
	Subject                 string `json:"subject"`
	MimeVersion             string `json:"mimeVersion"`
	ContentType             string `json:"contentType"`
	ContentTransferEncoding string `json:"contentTransferEncoding"`
	Cc                      string `json:"cc"`
	Bcc                     string `json:"bCc"`
	XFrom                   string `json:"xFrom"`
	XTo                     string `json:"xTo"`
	Xcc                     string `json:"xCc"`
	Xbcc                    string `json:"xBcc"`
	XFolder                 string `json:"xFolder"`
	XOrigin                 string `json:"xOrigin"`
	XFileName               string `json:"xFileName"`
	Content                 string `json:"content"`
}
