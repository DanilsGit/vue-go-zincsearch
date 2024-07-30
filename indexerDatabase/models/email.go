package models

import "time"

// Email struct
type Email struct {
	MessageID               string    `json:"message_id"`
	Date                    time.Time `json:"date"`
	From                    string    `json:"from"`
	To                      string    `json:"to"`
	Subject                 string    `json:"subject"`
	MimeVersion             string    `json:"mime_version"`
	ContentType             string    `json:"content_type"`
	ContentTransferEncoding string    `json:"content_transfer_encoding"`
	XFrom                   string    `json:"x_from"`
	XTo                     string    `json:"x_to"`
	Xcc                     string    `json:"x_cc"`
	Xbcc                    string    `json:"x_bcc"`
	XFolder                 string    `json:"x_folder"`
	XOrigin                 string    `json:"x_origin"`
	XFilename               string    `json:"x_filename"`
	Content                 string    `json:"content"`
}

// Bulk struct
type Bulk struct {
	Index   string  `json:"index"`
	Records []Email `json:"records"`
}
