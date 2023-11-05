package landing

import (
	"github.com/tylermmorton/testmail/app/model"
	"github.com/tylermmorton/testmail/app/services/smtp"
	"github.com/tylermmorton/testmail/app/templates/html"
	"github.com/tylermmorton/tmpl"
	"github.com/tylermmorton/torque"
	"html/template"
	"log"
	"net/http"
)

// Template can be used to render the landing page.
var Template = tmpl.MustCompile(
	&LandingPage{},
	tmpl.UseFuncs(tmpl.FuncMap{
		"html": func(v string) template.HTML {
			return template.HTML(v)
		},
	}),
)

//tmpl:bind landing.tmpl.html
type LandingPage struct {
	// Page is a template for a base html page.
	// It exposes the `body` template slot.
	html.Page `tmpl:"page"` // <- name the template, so it can be used as a target

	// Emails is a list of emails to display in the inbox.
	Emails []*model.Email

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
	q, err := torque.DecodeQuery[model.FindEmailQuery](req)
	if err != nil || q == nil {
		return nil, err
	}

	var current *model.Email
	if q.ID != "" {
		current, err = rm.SmtpService.GetEmailByID(req.Context(), q.ID)
		if err != nil {
			return nil, err
		}
	}

	emails, err := rm.SmtpService.FindEmails(req.Context(), q)
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
	log.Printf("[landing] rendering landing page with %d emails\n", len(loaderData.Emails))

	return Template.Render(wr,
		&LandingPage{
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
			Emails:  loaderData.Emails,
			Current: loaderData.Current,
		},
		tmpl.WithName("body"),   // <- assign the landing page template to the `body` slot
		tmpl.WithTarget("page"), // <- render the `page` template, which contains the `body`
	)
}
