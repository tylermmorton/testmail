package smtp

import (
	"context"
	"github.com/emersion/go-smtp"
	"github.com/tylermmorton/testmail/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	FindEmails(ctx context.Context, q *model.FindEmailQuery) ([]*model.Email, error)
}

func New(db *mongo.Database) Service {
	// TODO: implement mongodb index creation
	return &baseLayer{emails: db.Collection(emailsCollection)}
}

type baseLayer struct {
	emails *mongo.Collection
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

func (l *baseLayer) FindEmails(ctx context.Context, q *model.FindEmailQuery) ([]*model.Email, error) {
	opts := options.Find()
	opts.SetLimit(100)
	opts.SetSort(bson.D{{"createdAt", -1}})

	filter := bson.D{}
	if len(q.ID) > 0 {
		oid, err := primitive.ObjectIDFromHex(q.ID)
		if err != nil {
			return nil, err
		}
		filter = append(filter, bson.E{Key: "_id", Value: oid})
	}

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
