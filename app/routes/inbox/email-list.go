package inbox

import (
	"fmt"
	"github.com/tylermmorton/testmail/app/model"
	"github.com/tylermmorton/tmpl"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

var (
	EmailListTemplate = tmpl.MustCompile(&EmailList{}, tmpl.UseFuncs(tmpl.FuncMap{
		"hex": func(id primitive.ObjectID) string {
			return id.Hex()
		},
		"join": func(strs []string) string {
			return strings.Join(strs, ", ")
		},
		"formatTimeSince": formatTimeSince,
	}))
)

//tmpl:bind email-list.tmpl.html
type EmailList struct {
	// Emails is a list of emails to display in the inbox.
	Emails []*model.Email
	// Current is the currently selected email.
	Current *model.Email
}

func formatTimeSince(v time.Time) string {
	// round down to nearest int (floor)
	since := time.Since(v).Round(time.Hour)
	if val := since.Hours() / 24; val > 7 {
		return fmt.Sprintf("%dw", int(val/7))
	} else if since.Hours()/24 > 1 {
		return fmt.Sprintf("%dd", int(since.Hours()/24))
	} else if since.Hours() > 1 {
		return fmt.Sprintf("%dh", int(since.Hours()))
	} else if since.Minutes() > 1 {
		return fmt.Sprintf("%dm", int(since.Minutes()))
	} else if since.Seconds() > 1 {
		return fmt.Sprintf("%ds", int(since.Seconds()))
	} else {
		return "now"
	}
}
