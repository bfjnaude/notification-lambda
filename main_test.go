package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestCreateEmailFromMessage(t *testing.T) {

	subject := string("Test subject")
	sender := string("sender@example.com")
	recipient := string("recipiet@example.com")
	body := string("<h1>Test</h1>")

	message := events.SQSMessage{
		Body: body,
		MessageAttributes: map[string]events.SQSMessageAttribute{
			"subject":   events.SQSMessageAttribute{StringValue: &subject},
			"sender":    events.SQSMessageAttribute{StringValue: &sender},
			"recipient": events.SQSMessageAttribute{StringValue: &recipient},
		},
	}

	response, err := CreateEmailFromMessage(&message)

	assert.Equal(t, *response.Destination.ToAddresses[0], recipient)
	assert.Equal(t, *response.Source, sender)
	assert.Contains(t, *response.Message.Body.Html.Data, body)
	assert.Contains(t, *response.Message.Subject.Data, "Test subject")

	assert.Equal(t, err, nil)
}
