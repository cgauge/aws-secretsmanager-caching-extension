package secrets

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func GetSecretClient() *secretsmanager.SecretsManager {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		panic(err)
	}

	return secretsmanager.New(sess)
}

func GetSecret(name string) string {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(name),
	}

	println("Getting secret: ", name)

	result, err := GetSecretClient().GetSecretValue(input)
	if err != nil {
		panic(err.Error())
	}

	var secretString string
	if result.SecretString != nil {
		secretString = *result.SecretString
	}

	return secretString
}
