package main

import (
	"embed"
	"github.com/tylermmorton/testmail/app/routes/landing"
	"github.com/tylermmorton/testmail/app/services/smtp"
	"github.com/tylermmorton/torque"
	"io/fs"
	"log"
	"net/http"
	"sync"
)

//go install github.com/tylermmorton/tmpl/cmd/tmpl@latest
//go:generate tmpl bind ./...

//go:embed .build/assets
var embeddedAssets embed.FS

func main() {
	wg := sync.WaitGroup{}
	staticAssets, err := fs.Sub(embeddedAssets, ".build/assets")
	if err != nil {
		log.Fatalf("failed to create static assets filesystem: %+v", err)
	}

	// create a shared smtp service between client and server
	smtpService := smtp.New()

	r := torque.NewRouter(
		torque.WithRouteModule("/", &landing.RouteModule{}),
		torque.WithFileSystemServer("/s", staticAssets),
	)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		log.Printf("[torque] Listening on http://localhost:8080")
		err = http.ListenAndServe(":8080", r)
		if err != nil {
			log.Fatalf("failed to start http server: %+v", err)
		}
	}(&wg)
	wg.Add(1)

	s := smtp.NewServer(smtpService)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		log.Printf("[smtp] Listening on http://localhost:1025")
		err = s.ListenAndServe()
		if err != nil {
			log.Fatalf("failed to start smtp server: %+v", err)
		}
	}(&wg)
	wg.Add(1)

	wg.Wait()
}
