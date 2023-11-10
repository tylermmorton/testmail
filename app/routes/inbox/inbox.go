package inbox

import (
	"github.com/tylermmorton/testmail/app/model"
	"github.com/tylermmorton/testmail/app/services/smtp"
	"github.com/tylermmorton/testmail/app/templates/html"
	"github.com/tylermmorton/tmpl"
	"github.com/tylermmorton/torque"
	"github.com/tylermmorton/torque/pkg/htmx"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// Template can be used to render the inbox page.
var Template = tmpl.MustCompile(
	&InboxPage{},
	tmpl.UseFuncs(tmpl.FuncMap{
		"html": func(v string) template.HTML {
			return template.HTML(v)
		},
		"hex": func(id primitive.ObjectID) string {
			return id.Hex()
		},
		"join": func(strs []string) string {
			return strings.Join(strs, ", ")
		},
		"formatTimeSince": formatTimeSince,
	}),
)

//tmpl:bind inbox.tmpl.html
type InboxPage struct {
	// Page is a template for a base html page.
	// It exposes the `body` template slot.
	html.Page `tmpl:"page"` // <- name the template, so it can be used as a target

	EmailList EmailList `tmpl:"email-list"`
	EmailView EmailView `tmpl:"email-view"`

	// Current is the currently selected email.
	Current *model.Email
}

// RouteModule is the torque route module for the landing page.
type RouteModule struct {
	SmtpService smtp.Service
}

var _ interface {
	torque.Action
	torque.Loader
	torque.Renderer
	torque.EventSource
} = &RouteModule{}

func (rm *RouteModule) Action(wr http.ResponseWriter, req *http.Request) error {
	action := torque.DecodeFormAction(req)
	switch action {
	case "delete":
		emailId := torque.RouteParam(req, "emailId")

		err := rm.SmtpService.DeleteEmailByID(req.Context(), emailId)
		if err != nil {
			return err
		}

		http.Redirect(wr, req, "/", http.StatusFound)
		return nil
	}

	return nil
}

func (rm *RouteModule) Load(req *http.Request) (any, error) {
	var current *model.Email
	var emailId string = torque.RouteParam(req, "emailId")

	query, err := torque.DecodeQuery[model.FindEmailQuery](req)
	if err != nil || query == nil {
		return nil, err
	}

	current, err = rm.SmtpService.GetEmailByID(req.Context(), emailId)
	if err != nil {
		return nil, err
	}

	emails, err := rm.SmtpService.FindEmails(req.Context(), query)
	if err != nil {
		return nil, err
	}

	return &LoaderData{
		Emails:  emails,
		Current: current,
	}, nil
}

type LoaderData struct {
	Emails  []*model.Email
	Current *model.Email
}

func (rm *RouteModule) Render(wr http.ResponseWriter, req *http.Request, ld any) error {
	loaderData := ld.(*LoaderData)

	return torque.VaryRender(wr, req, htmx.HxRequestHeader, map[any]torque.RenderFn{
		"true": func(wr http.ResponseWriter, req *http.Request) error {
			return Template.Render(wr, &InboxPage{
				Current:   loaderData.Current,
				EmailView: EmailView{Current: loaderData.Current},
			}, tmpl.WithTarget("email-view"))
		},

		torque.VaryDefault: func(wr http.ResponseWriter, req *http.Request) error {
			return Template.Render(wr,
				&InboxPage{
					Page: html.Page{
						TitlePrefix: "Welcome",
						Title:       "create-torque-app",
						Links: []html.Link{
							{Rel: "stylesheet", Href: "/s/app.css"},
						},
						Scripts: []html.Script{
							{Src: "https://unpkg.com/htmx.org@1.9.6"},
							{Src: "https://unpkg.com/htmx.org/dist/ext/sse.js"},
						},
					},
					EmailList: EmailList{Emails: loaderData.Emails, Current: loaderData.Current},
					EmailView: EmailView{Current: loaderData.Current},
					Current:   loaderData.Current,
				},
				tmpl.WithName("body"),
				tmpl.WithTarget("page"),
			)
		},
	})
}

func (rm *RouteModule) Subscribe(wr http.ResponseWriter, req *http.Request) error {
	return htmx.SSE(wr, req, htmx.EventSourceMap{
		"email-created": func(sse chan string) {
			ch, err := rm.SmtpService.WatchEmails(req.Context())
			if err != nil {
				panic(err)
			}

			for {
				select {
				case <-req.Context().Done():
					// the text/event-stream connection was closed
					// unsubscribe from real time updates for messages
					return

				case _, ok := <-ch:
					if !ok {
						return
					}

					loaderData, err := rm.Load(req)
					if err != nil {
						log.Printf("failed to load emails: %v", err)
					}

					err = EmailListTemplate.RenderToChan(sse, &EmailList{
						Current: loaderData.(*LoaderData).Current,
						Emails:  loaderData.(*LoaderData).Emails,
					})
					if err != nil {
						log.Printf("failed to render email list: %v", err)
					}
				}
			}
		},
	})
}
