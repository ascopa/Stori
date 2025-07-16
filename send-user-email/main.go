package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"send-user-email/internal/factory"
)

func main() {
	var f factory.Factory
	lambda.Start(f.Start)
}
