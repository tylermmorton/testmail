package html

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

//go:embed link.tmpl.html
var LinkTmplFS embed.FS

func (t *Link) TemplateText() string {
	return _tmpl(LinkTmplFS, ".")
}

//go:embed page.tmpl.html
var PageTmplFS embed.FS

func (t *Page) TemplateText() string {
	return _tmpl(PageTmplFS, ".")
}

//go:embed script.tmpl.html
var ScriptTmplFS embed.FS

func (t *Script) TemplateText() string {
	return _tmpl(ScriptTmplFS, ".")
}
