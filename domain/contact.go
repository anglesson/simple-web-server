package domain

type Contact struct {
	Email *Email
	Phone *Phone
}

func NewContact(email, phone string) (*Contact, error) {
	voEmail, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	voPhone, err := NewPhone(phone)
	if err != nil {
		return nil, err
	}

	return &Contact{
		Email: voEmail,
		Phone: voPhone,
	}, nil
}
