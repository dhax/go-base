package main

import (
	"flag"
	"fmt"
	"github.com/unrolled/render"
	"github.com/vanng822/go-premailer/premailer"
	"github.com/vanng822/r2router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var (
		host string
		port int
	)

	flag.StringVar(&host, "h", "127.0.0.1", "Host to listen on")
	flag.IntVar(&port, "p", 8080, "Port number to listen on")
	flag.Parse()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Kill, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR2)
	app := r2router.NewSeeforRouter()

	r := render.New()

	app.Get("/", func(w http.ResponseWriter, req *http.Request, _ r2router.Params) {
		r.JSON(w, http.StatusOK, r2router.M{"usage": "POST / html=HTML&cssToAttributes=boolean&removeClasses=boolean"})
	})
	app.Post("/", func(w http.ResponseWriter, req *http.Request, _ r2router.Params) {
		req.ParseForm()
		html := req.Form.Get("html")
		cssToAttributes := req.Form.Get("cssToAttributes")
		removeClasses := req.Form.Get("removeClasses")
		var result string
		if html != "" {
			options := premailer.NewOptions()
			if removeClasses == "true" {
				options.RemoveClasses = true
			}
			if cssToAttributes == "false" {
				options.CssToAttributes = false
			}
			pre := premailer.NewPremailerFromString(html, options)
			result, _ = pre.Transform()
		} else {
			result = ""
		}
		r.JSON(w, http.StatusOK, r2router.M{"result": result})
	})

	log.Printf("listening to address %s:%d", host, port)
	go http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), app)
	sig := <-sigc
	log.Printf("Got signal: %s", sig)

}
