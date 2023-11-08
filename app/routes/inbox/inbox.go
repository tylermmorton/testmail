package inbox

import (
	"github.com/tylermmorton/testmail/app/model"
	"github.com/tylermmorton/testmail/app/services/smtp"
	"github.com/tylermmorton/testmail/app/templates/html"
	"github.com/tylermmorton/tmpl"
	"github.com/tylermmorton/torque"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
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

	EmailList `tmpl:"email-list"`

	// Current is the currently selected email.
	Current *model.Email
}

// RouteModule is the torque route module for the landing page.
type RouteModule struct {
	SmtpService smtp.Service
}

var _ interface {
	torque.Loader
	torque.Renderer
} = &RouteModule{}

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
				},
			},
			EmailList: EmailList{Emails: loaderData.Emails, Current: loaderData.Current},
			Current:   loaderData.Current,
		},
		tmpl.WithName("body"),
		tmpl.WithTarget("page"),
	)
}
