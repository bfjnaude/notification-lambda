package main

import (
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// CreateEmailFromMessage converts an events.SQSMessage to an ses.SendEmailInput
// The SQSMessage should have the email recipient, sender and subject defined as MessageAttributes
// The SQSMessage should have the email body defined as the message body
func CreateEmailFromMessage(message *events.SQSMessage) (*ses.SendEmailInput, error) {
	// Recipient
	attr, b := message.MessageAttributes["recipient"]
	if !b {
		return nil, errors.New("'recipient' missing from message")
	}
	recipient := *attr.StringValue

	// Sender
	attr, b = message.MessageAttributes["sender"]
	if !b {
		return nil, errors.New("'sender' missing from message")
	}
	sender := *attr.StringValue

	// Character set
	charset := "UTF-8" // default to utf-8
	attr, b = message.MessageAttributes["charset"]
	if b {
		charset = *attr.StringValue
	}

	// Body
	htmlBody := message.Body

	// Subject
	attr, b = message.MessageAttributes["subject"]
	if !b {
		return nil, errors.New("'subject' missing from message")
	}
	subject := *attr.StringValue

	// Create email
	email := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charset),
					Data:    aws.String(htmlBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charset),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(sender),
	}
	return email, nil
}

var sesSession *session.Session
var sesService *ses.SES

func init() {
	sesSession = session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("us-east-1"),
		},
	}))

	sesService = ses.New(sesSession)
}

// Handler processes sqs events to send emails
func Handler(sqsEvent events.SQSEvent) error {
	// Get messages
	for _, message := range sqsEvent.Records {

		// Extract email fields from message
		email, err := CreateEmailFromMessage(&message)
		if err != nil {
			fmt.Printf("Could create email %s\n", err)
			return err
		}

		// Send it
		sendResult, err := sesService.SendEmail(email)

		// handle errors
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case ses.ErrCodeMessageRejected:
					fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
				case ses.ErrCodeMailFromDomainNotVerifiedException:
					fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
				case ses.ErrCodeConfigurationSetDoesNotExistException:
					fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}

			return err
		}

		fmt.Println("Email sent", sendResult)
	}
	return nil
}

func main() {
	lambda.Start(Handler)
}
