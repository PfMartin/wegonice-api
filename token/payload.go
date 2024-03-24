package token

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExpiredToken = fmt.Errorf("token has expired")
	ErrInvalidToken = fmt.Errorf("token is invalid")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserID    string    `json:"userId"`
	IssuedAt  int64     `json:"issuedAt"`
	ExpiresAt int64     `json:"expiresAt"`
}

func NewPayload(userID string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		UserID:    userID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(duration).Unix(),
	}

	return payload, nil
}

func (p *Payload) Validate() error {
	if time.Now().Unix() > p.ExpiresAt {
		return ErrExpiredToken
	}

	return nil
}
