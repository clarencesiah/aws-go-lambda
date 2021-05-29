package main

// snippet-start:[dynamodb.go.load_items.imports]
import (
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"encoding/json"
	"io/ioutil"
	"log"
)

// Movie entity
type Movie struct {
	ID   string
	Name string
}

func findAll(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	size, err := strconv.Atoi(request.Headers["Count"])
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Count Header should be a number",
		}, nil
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	req, res := svc.ScanRequest(&dynamodb.ScanInput{
		TableName: aws.String("movies"),
		Limit:     aws.Int64(int64(size)),
	})

	err = req.Send()

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusTeapot,
			Body:       "Count Header should be a number",
		}, nil
	}

	response, err := json.Marshal(res.Items)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(response),
	}, nil
}

func main() {
	lambda.Start(findAll)
}
