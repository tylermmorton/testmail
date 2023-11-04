package smtp

import (
	"context"
	"github.com/emersion/go-smtp"
	"github.com/tylermmorton/testmail/app/model"
	"log"
)

type Service interface {
	smtp.Backend
}

func New() Service {
	return &baseLayer{}
}

type baseLayer struct {
}

func (l *baseLayer) NewSession(c *smtp.Conn) (smtp.Session, error) {
	ctx := context.Background()
	ch := make(chan *model.Email)
	go func(ctx context.Context) {
		for {
			select {
			case email, ok := <-ch:
				if !ok {
					// channel closed
					return
				}
				if email == nil {
					continue
				}

				l.handleIncomingEmail(ctx, email)
			}
		}
	}(ctx)

	return &session{ctx: ctx, ch: ch}, nil
}

func (l *baseLayer) handleIncomingEmail(ctx context.Context, email *model.Email) {
	log.Printf("Received email: %+v\n", email)
}
