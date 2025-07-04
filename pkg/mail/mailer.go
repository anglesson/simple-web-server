package mail

type Mailer interface {
	From(email string)
	To(email string)
	Subject(subject string)
	Body(body string)
	Send()
}
