package inbox

// /!\ THIS FILE IS GENERATED DO NOT EDIT /!\

import (
	"embed"
	"path/filepath"
	"strings"
)

func _tmpl(fsys embed.FS, path string) string {
	builder := &strings.Builder{}
	entries, err := fsys.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			builder.WriteString(_tmpl(fsys, filepath.Join(path, entry.Name())))
		} else {
			byt, err := fsys.ReadFile(filepath.Join(path, entry.Name()))
			if err != nil {
				panic(err)
			}
			builder.Write(byt)
		}
	}
	return builder.String()
}

//go:embed email-list.tmpl.html
var EmailListTmplFS embed.FS

func (t *EmailList) TemplateText() string {
	return _tmpl(EmailListTmplFS, ".")
}

//go:embed email-view.tmpl.html
var EmailViewTmplFS embed.FS

func (t *EmailView) TemplateText() string {
	return _tmpl(EmailViewTmplFS, ".")
}

//go:embed inbox.tmpl.html
var InboxPageTmplFS embed.FS

func (t *InboxPage) TemplateText() string {
	return _tmpl(InboxPageTmplFS, ".")
}
