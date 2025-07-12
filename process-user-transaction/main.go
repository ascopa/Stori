package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"process-user-transaction/internal/factory"
)

func main() {
	var f factory.Factory
	lambda.Start(f.Start)
}
