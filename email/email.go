package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jaytaylor/html2text"
	"github.com/spf13/viper"
	"github.com/vanng822/go-premailer/premailer"
	"gopkg.in/gomail.v2"
)

var (
	debug bool
)

type EmailService struct {
	client         *gomail.Dialer
	templates      *template.Template
	from, fromName string
}

func NewEmailService() (*EmailService, error) {
	templates, err := parseTemplates()
	if err != nil {
		return nil, err
	}

	smtpHost := viper.GetString("email_smtp_host")
	smtpPort := viper.GetInt("email_smtp_port")
	smtpUser := viper.GetString("email_smtp_user")
	smtpPass := viper.GetString("email_smtp_password")

	s := &EmailService{
		client:    gomail.NewPlainDialer(smtpHost, smtpPort, smtpUser, smtpPass),
		templates: templates,
		from:      viper.GetString("email_from_address"),
		fromName:  viper.GetString("email_from_name"),
	}

	d, err := s.client.Dial()
	if err != nil {
		log.Println("SMTP error:", err)
		log.Println("printing emails to stdout")
		debug = true
	} else {
		d.Close()
	}
	return s, nil
}

func (s *EmailService) send(msg *Message) error {
	buf := new(bytes.Buffer)
	if err := s.templates.ExecuteTemplate(buf, msg.template, msg.content); err != nil {
		return err
	}
	prem := premailer.NewPremailerFromString(buf.String(), premailer.NewOptions())
	html, err := prem.Transform()
	if err != nil {
		return err
	}

	text, err := html2text.FromString(html, html2text.Options{PrettyTables: true})
	if err != nil {
		return err
	}

	if debug {
		log.Println("To:", msg.to.Address)
		log.Println("Subject:", msg.subject)
		log.Println(text)
		return nil
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", msg.from.Address, msg.from.Name)
	m.SetAddressHeader("To", msg.to.Address, msg.to.Name)
	m.SetHeader("Subject", msg.subject)
	m.SetBody("text/plain", text)
	m.AddAlternative("text/html", html)

	if err := s.client.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

type Email struct {
	Name    string
	Address string
}

func NewEmail(name string, address string) *Email {
	return &Email{
		Name:    name,
		Address: address,
	}
}

type Message struct {
	from     *Email
	to       *Email
	subject  string
	template string
	content  interface{}
}

func parseTemplates() (*template.Template, error) {
	tmpl := template.New("").Funcs(fMap)
	err := filepath.Walk("./templates", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			_, err = tmpl.ParseFiles(path)
			return err
		}
		return err
	})
	return tmpl, err
}

var fMap = template.FuncMap{
	"formatAsDate":     formatAsDate,
	"formatAsDuration": formatAsDuration,
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d.%d.%d", day, month, year)
}

func formatAsDuration(t time.Time) string {
	dur := t.Sub(time.Now())
	hours := int(dur.Hours())
	mins := int(dur.Minutes())

	v := ""
	if hours != 0 {
		v += strconv.Itoa(hours) + " hours and "
	}
	v += strconv.Itoa(mins) + " minutes"
	return v
}
