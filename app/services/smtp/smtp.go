package smtp

import (
	"context"
	"github.com/256dpi/lungo"
	"github.com/emersion/go-smtp"
	"github.com/tylermmorton/eventbus"
	"github.com/tylermmorton/testmail/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const (
	emailsCollection = "emails"
)

type Service interface {
	smtp.Backend

	GetEmailByID(ctx context.Context, id string) (*model.Email, error)
	DeleteEmailByID(ctx context.Context, id string) error
	FindEmails(ctx context.Context, q *model.FindEmailQuery) ([]*model.Email, error)

	WatchEmails(ctx context.Context) (<-chan *model.Email, error)
}

func New(db lungo.IDatabase) Service {
	return &baseLayer{
		emails:       db.Collection(emailsCollection),
		emailCreated: eventbus.New[string, *model.Email](),
	}
}

type baseLayer struct {
	emails       lungo.ICollection
	emailCreated eventbus.EventBus[string, *model.Email]
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

				err := l.handleIncomingEmail(ctx, email)
				if err != nil {
					log.Printf("[smtp] failed to handle incoming email: %+v\n", err)
					return
				}
			}
		}
	}(ctx)

	return &session{ctx: ctx, ch: ch}, nil
}

func (l *baseLayer) handleIncomingEmail(ctx context.Context, email *model.Email) error {
	email.ID = primitive.NewObjectID()
	email.CreatedAt = time.Now()

	res, err := l.emails.InsertOne(ctx, email)
	if err != nil {
		return err
	}

	if l.emailCreated != nil {
		l.emailCreated.Dispatch("*", email)
	}

	log.Printf("[smtp] saved email to database: %+v\n", res)
	return nil
}

func (l *baseLayer) GetEmailByID(ctx context.Context, id string) (*model.Email, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var email model.Email
	if err := l.emails.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&email); err != nil {
		return nil, err
	}

	return &email, nil
}

func (l *baseLayer) DeleteEmailByID(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	res, err := l.emails.DeleteOne(ctx, bson.D{{"_id", oid}})
	if err != nil {
		return err
	}

	log.Printf("[smtp] deleted emails from database: %+v\n", res)

	return nil
}

func (l *baseLayer) FindEmails(ctx context.Context, q *model.FindEmailQuery) ([]*model.Email, error) {
	opts := options.Find()
	opts.SetLimit(q.Limit)
	opts.SetSort(bson.D{{"createdAt", -1}})

	cur, err := l.emails.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}

	var emails []*model.Email
	for cur.Next(ctx) {
		var email model.Email
		if err := cur.Decode(&email); err != nil {
			return nil, err
		}
		emails = append(emails, &email)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return emails, nil
}

func (l *baseLayer) WatchEmails(ctx context.Context) (<-chan *model.Email, error) {
	ch := make(chan *model.Email)

	l.emailCreated.Subscribe("*", ch)

	go func(ctx context.Context, ch chan *model.Email) {
		<-ctx.Done()
		err := l.emailCreated.Unsubscribe("*", ch)
		if err != nil {
			log.Printf("failed to unsubscribe from emailCreated eventbus: %v", err)
		}
		close(ch)
	}(ctx, ch)

	return ch, nil
}
