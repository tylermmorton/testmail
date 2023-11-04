package smtp

import (
	"context"
	"encoding/base64"
	"github.com/emersion/go-smtp"
	"github.com/tylermmorton/testmail/app/model"
	"io"
	"strings"
)

type session struct {
	ctx context.Context
	ch  chan *model.Email
}

func (s *session) Data(r io.Reader) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	// bytesRead is the number of bytes read so far
	bytesRead := 0
	headers := map[string]string{}
	lines := strings.Split(string(buf), "\n")
	for _, line := range lines {
		bytesRead += len(line) + 1
		if len(line) == 0 {
			continue
		}

		if line == "\r" {
			// we are at the end of the headers
			break
		}

		kv := strings.SplitN(line, ":", 2)
		if len(kv) != 2 {
			continue
		}

		headers[kv[0]] = strings.TrimSpace(kv[1])
	}

	// if the body is encoded in some way, decode it
	body := string(buf[bytesRead:])
	switch headers["Content-Transfer-Encoding"] {
	case "base64":
		byt, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			return err
		}
		body = string(byt)
	}

	email := model.Email{}

	// promote some headers for search indexing
	email.From = headers["From"]
	delete(headers, "From")

	email.To = strings.Split(headers["To"], ",")
	delete(headers, "To")

	email.Subject = headers["Subject"]
	delete(headers, "Subject")

	email.Headers = headers
	email.Body = body

	s.ch <- &email

	return nil
}

func (s *session) Logout() error {
	close(s.ch)
	return nil
}

// Don't need to implement these lifecycle methods -- testmail accepts all incoming mail

func (s *session) AuthPlain(username, password string) error      { return nil }
func (s *session) Mail(from string, opts *smtp.MailOptions) error { return nil }
func (s *session) Rcpt(to string, opts *smtp.RcptOptions) error   { return nil }
func (s *session) Reset()                                         {}
