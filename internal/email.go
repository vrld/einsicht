package internal

import (
	"io"
	"mime"
	"net/textproto"
	"slices"
	"strings"

	"github.com/emersion/go-message"
	_ "github.com/emersion/go-message/charset"
)

type Email struct {
	From        string
	ReplyTo     []string
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	Date        string
	MIMEHeader  textproto.MIMEHeader
	Text        string
	HTML        string
	Attachments []*Attachment
}

type HeaderDisplay struct {
	Key   string
	Value string
}

type Attachment struct {
	Filename    string
	ContentType string
	Header      textproto.MIMEHeader
	Content     []byte
	Cid         string
}

func ReadEmail(reader io.Reader) (*Email, error) {
	m, err := message.Read(reader)
	if err != nil {
		return nil, err
	}

	headers := decodeMIMEHeaders(m.Header.Map())

	email := Email{
		From:       getSingleValueFromHeaders(headers, "From"),
		ReplyTo:    headers.Values("ReplyTo"),
		To:         headers.Values("To"),
		Cc:         headers.Values("Cc"),
		Bcc:        headers.Values("Bcc"),
		Subject:    getSingleValueFromHeaders(headers, "Subject"),
		Date:       getSingleValueFromHeaders(headers, "Date"),
		MIMEHeader: headers,
	}

	err = readParts(m, &email)
	return &email, err
}

func readParts(entity *message.Entity, email *Email) error {
	mr := entity.MultipartReader()
	if mr == nil {
		return readSinglePart(entity, email)
	}

	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		} else if err = readSinglePart(p, email); err != nil {
			return err
		}
	}
	return nil
}

func readSinglePart(entity *message.Entity, email *Email) error {
	contentType, contentParams, err := entity.Header.ContentType()
	if err != nil {
		return err
	}

	switch contentType {
	case "", "text/plain":
		body, err := readAsString(entity.Body)
		if err != nil {
			return err
		}
		email.Text = body

	case "text/html":
		body, err := readAsString(entity.Body)
		if err != nil {
			return err
		}
		email.HTML = body

	case "multipart/related", "multipart/alternative":
		return readParts(entity, email)

	default:
		filename := contentParams["name"]
		if filename == "" {
			filename = entity.Header.Get("Content-Description")
		}
		if filename == "" {
			_, params, err := entity.Header.ContentDisposition()
			if err != nil {
				return err
			}
			filename = params["filename"]
		}

		attachment := Attachment{
			Filename:    filename,
			ContentType: contentType,
			Header:      entity.Header.Map(),
		}

		attachment.Content, err = io.ReadAll(entity.Body)
		if err != nil {
			return err
		}

		if cid := entity.Header.Get("Content-Id"); cid != "" {
			cid = strings.Trim(cid, "<>")
			attachment.Cid = cid
		}

		email.Attachments = append(email.Attachments, &attachment)

	}
	return nil
}

func readAsString(reader io.Reader) (string, error) {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (email *Email) HeaderDisplayCanonical() []HeaderDisplay {
	return []HeaderDisplay{
		{"Date", email.Date},
		{"From", email.From},
		{"To", strings.Join(email.To, ", ")},
		{"Cc", strings.Join(email.Cc, ", ")},
		{"Subject", email.Subject},
	}
}

func (email *Email) HeaderDisplayAdditional() []HeaderDisplay {
	selectedHeaders := []string{}
	for header := range email.MIMEHeader {
		switch header {
		case "Content-Type",
			"Content-Transfer-Encoding",
			"Mime-Version",
			"Date",
			"From",
			"To",
			"Cc",
			"Subject":
			// skip

		default:
			selectedHeaders = append(selectedHeaders, header)
		}
	}

	slices.Sort(selectedHeaders)

	res := []HeaderDisplay{}
	for _, header := range selectedHeaders {
		value := strings.Join(email.MIMEHeader.Values(header), ", ")
		res = append(res, HeaderDisplay{header, value})
	}

	return res
}

func (email *Email) HeaderDisplayAll() []HeaderDisplay {
	return append(email.HeaderDisplayCanonical(), email.HeaderDisplayAdditional()...)
}

func decodeMIMEHeaders(rawHeaders textproto.MIMEHeader) textproto.MIMEHeader {
	decoder := &mime.WordDecoder{}
	decoded := make(textproto.MIMEHeader)

	for key, values := range rawHeaders {
		decodedValues := make([]string, len(values))
		for i, value := range values {
			if decodedValue, err := decoder.DecodeHeader(value); err == nil {
				decodedValues[i] = decodedValue
			} else {
				decodedValues[i] = value
			}
		}
		decoded[key] = decodedValues
	}
	return decoded
}

func getSingleValueFromHeaders(headers textproto.MIMEHeader, key string) string {
	if values := headers.Values(key); len(values) > 0 {
		return values[0]
	}
	return ""
}
