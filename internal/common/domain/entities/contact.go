package common_entities

import common_vo "github.com/anglesson/simple-web-server/internal/common/domain/valueobjects"

type Contact struct {
	Email *common_vo.Email
	Phone *common_vo.Phone
}

func NewContact(email, phone string) (*Contact, error) {
	voEmail, err := common_vo.NewEmail(email)
	if err != nil {
		return nil, err
	}

	voPhone, err := common_vo.NewPhone(phone)
	if err != nil {
		return nil, err
	}

	return &Contact{
		Email: voEmail,
		Phone: voPhone,
	}, nil
}
