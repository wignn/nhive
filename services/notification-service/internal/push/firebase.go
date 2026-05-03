package push

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type FirebasePusher struct {
	client *messaging.Client
}

func NewFirebasePusher(credentialsPath string) (*FirebasePusher, error) {
	if credentialsPath == "" {
		return &FirebasePusher{client: nil}, nil 
	}

	opt := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %v", err)
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting messaging client: %v", err)
	}

	return &FirebasePusher{client: client}, nil
}

func (p *FirebasePusher) SendPushNotification(tokens []string, title, body string, data map[string]string) error {
	if p.client == nil || len(tokens) == 0 {
		return nil
	}

	message := &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
	}

	br, err := p.client.SendEachForMulticast(context.Background(), message)
	if err != nil {
		return err
	}

	if br.FailureCount > 0 {
		fmt.Printf("FCM failure count: %d\n", br.FailureCount)
	}

	return nil
}
