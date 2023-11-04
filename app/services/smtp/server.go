package smtp

import (
	"fmt"
	"github.com/emersion/go-smtp"
	"github.com/kelseyhightower/envconfig"
	"time"
)

type serverConfig struct {
	Port   string `envconfig:"PORT"`
	Domain string `envconfig:"DOMAIN"`
}

func NewServer(svc Service) (*smtp.Server, error) {
	var s = smtp.NewServer(svc)
	var cfg = serverConfig{
		Port:   "1025",
		Domain: "localhost",
	}
	var err = envconfig.Process("SMTP", &cfg)
	if err != nil {
		return nil, err
	}
	s.Addr = fmt.Sprintf(":%s", cfg.Port)
	s.Domain = cfg.Domain
	s.WriteTimeout = 10 * time.Second
	s.ReadTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxLineLength = 1024 * 1024
	s.AllowInsecureAuth = true
	//s.Debug = os.Stdout
	return s, nil
}
