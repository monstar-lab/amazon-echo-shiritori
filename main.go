package main

import (
	"errors"

	"./alexa"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	ErrInvalidIntent = errors.New("Invalid intent")
)

func HandleRequest(event alexa.Request) (alexa.Response, error) {
	eventRequestType := event.Request.Type
	if eventRequestType == "LaunchRequest" {
		return OnLaunch(event.Request)
	} else if eventRequestType == "IntentRequest" {
		return OnIntent(event.Request)
	}
	return alexa.Response{}, ErrInvalidIntent
}

// OnLaunch is function-type
func OnLaunch(launchRequest alexa.RequestDetail) (alexa.Response, error) {
	return GetWelcomeResponse(), nil
}

// GetWelcomeResponse is function-type
func GetWelcomeResponse() alexa.Response {
	cardTitle := "Hello"
	speechOutput := "Alexaです。お名前は"
	repromptText := "Hello スキル起動"
	shouldEndSession := false
	return alexa.BuildResponse(alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession))
}

//OnIntent is function-type
func OnIntent(intentRequest alexa.RequestDetail) (alexa.Response, error) {
	if intentRequest.Intent.Name == "HelloIntent" {
		return HelloIntent(intentRequest.Intent.Slots["userName"].Value), nil
	}
	return GetWelcomeResponse(), nil
}

func HelloIntent(name string) alexa.Response {
	cardTitle := "HelloIntent"
	speechOutput := "Hello" + name
	repromptText := "Hello" + name
	shouldEndSession := true
	return alexa.BuildResponse(alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession))
}

func main() {
	lambda.Start(HandleRequest)
}
