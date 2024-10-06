package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
        slog.Error("could not load config", "error", err.Error())
    }

	svc := sqs.NewFromConfig(cfg)

	slog.Info("NotificationService starting")
	queueURL := "https://sqs.us-west-2.amazonaws.com/123456789012/MyQueue"
	
	listenForMessages(svc, queueURL)
}

func listenForMessages(svc *sqs.Client, queueURL string) {
    for {
        input := &sqs.ReceiveMessageInput{
            QueueUrl:            aws.String(queueURL),
            MaxNumberOfMessages: 10,
            WaitTimeSeconds:     20,
        }

        result, err := svc.ReceiveMessage(context.TODO(), input)
        if err != nil {
            slog.Error("could not receive message", "error", err.Error())
            time.Sleep(5 * time.Second) // Wait before retrying
            continue
        }

        for _, message := range result.Messages {
            fmt.Printf("Message received: %s\n", *message.Body)

            // Process the message
            processMessage(message)

            // Delete the message after processing
            deleteMessage(svc, queueURL, message.ReceiptHandle)
        }
    }
}

func processMessage(message sqs.Message) {
    // Add your message processing logic here
    fmt.Printf("Processing message: %s\n", *message.Body)
}

func deleteMessage(svc *sqs.Client, queueURL string, receiptHandle *string) {
    input := &sqs.DeleteMessageInput{
        QueueUrl:      aws.String(queueURL),
        ReceiptHandle: receiptHandle,
    }

    _, err := svc.DeleteMessage(context.TODO(), input)
    if err != nil {
        slog.Error("could not delete message", "error", err.Error())
    } else {
		slog.Info("message deleted")
    }
}