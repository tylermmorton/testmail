package inbox

import "github.com/tylermmorton/testmail/app/model"

//tmpl:bind email-view.tmpl.html
type EmailView struct {
	Current *model.Email
}
