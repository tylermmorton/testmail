package smtp

import (
	"github.com/emersion/go-smtp"
	"io"
	"os"
)

type Service interface {
	smtp.Backend
}

func New() Service {
	return &baseLayer{}
}

func NewServer(svc Service) (s *smtp.Server) {
	s = smtp.NewServer(svc)
	s.Addr = ":1025"
	s.Domain = "localhost"
	s.AllowInsecureAuth = true
	s.Debug = os.Stdout
	return s
}

type baseLayer struct{}

func (l *baseLayer) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &session{}, nil
}

type session struct{}

func (s *session) AuthPlain(username, password string) error {
	return nil
}

func (s *session) Mail(from string, opts *smtp.MailOptions) error {
	return nil
}

func (s *session) Rcpt(to string, opts *smtp.RcptOptions) error {
	return nil
}

func (s *session) Data(r io.Reader) error {
	return nil
}

func (s *session) Reset() {}

func (s *session) Logout() error {
	return nil
}
