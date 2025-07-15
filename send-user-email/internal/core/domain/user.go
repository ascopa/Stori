package domain

type User struct {
	AccountId string `dynamodbav:"AccountId"`
	Name      string `dynamodbav:"Name"`
	Email     string `dynamodbav:"Email"`
}
