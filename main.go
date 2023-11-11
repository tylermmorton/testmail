package main

import (
	"context"
	"embed"
	"github.com/256dpi/lungo"
	"github.com/tylermmorton/testmail/app/routes/inbox"
	"github.com/tylermmorton/testmail/app/routes/landing"
	"github.com/tylermmorton/testmail/app/services/smtp"
	"github.com/tylermmorton/torque"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

//go install github.com/tylermmorton/tmpl/cmd/tmpl@latest
//go:generate tmpl bind ./... --mode=embed

//go:embed .build/assets
var embeddedAssets embed.FS

func main() {
	wg := sync.WaitGroup{}
	staticAssets, err := fs.Sub(embeddedAssets, ".build/assets")
	if err != nil {
		log.Fatalf("failed to create static assets filesystem: %+v", err)
	}

	fsPath := os.Getenv("LUNGO_FILESTORE_DIR")
	if len(fsPath) == 0 {
		fsPath = ".lungo"
	}
	log.Printf("[testmail] using %s to store embedded mongodb", fsPath)

	opts := lungo.Options{
		Store:          lungo.NewFileStore(fsPath+"testmail.db", 0777),
		ExpireInterval: 0,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, engine, err := lungo.Open(ctx, opts)
	if err != nil {
		log.Fatalf("failed to open lungo in-memory database: %+v", err)
	}
	defer engine.Close()

	db := mongoClient.Database("testmail")

	// create a shared smtp service between client and server
	smtpService := smtp.New(db)

	r := torque.NewRouter(
		torque.WithRouteModule("/", &landing.RouteModule{SmtpService: smtpService}),
		torque.WithRouteModule("/{emailId}", &inbox.RouteModule{SmtpService: smtpService}),
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

	s, err := smtp.NewServer(smtpService)
	if err != nil {
		log.Fatalf("failed to create smtp server: %+v", err)
	}
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
