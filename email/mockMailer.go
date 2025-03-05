package email

import "log"

// MockMailer is a mock Mailer
type MockMailer struct {
	SendFn      func(m Message) error
	SendInvoked bool
}

func logMessage(m Message) {
	log.Printf("MockMailer email sent:\nSubject: %s\nTo: %s <%s>\nContext: %#v\n", m.Subject, m.To.Name, m.To.Address, m.Content)
}

func NewMockMailer() *MockMailer {
	log.Println("ATTENTION: SMTP Mailer not configured => printing emails to stdout")
	return &MockMailer{
		SendFn: func(m Message) error {
			logMessage(m)
			return nil
		},
	}
}

func (s *MockMailer) Send(m Message) error {
	s.SendInvoked = true
	return s.SendFn(m)
}
