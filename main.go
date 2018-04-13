package main

import (
	"errors"
	"log"
	"strings"

	"./alexa"
	"./db"
	"./timeData"
	"github.com/aws/aws-lambda-go/lambda"
)

const GAME_START_MESSAGE = `しりとりスキルへようこそ。「りんご」`

const RES_SHIRITORI_SLOT = "shiritoriword"

var (
	ErrInvalidIntent = errors.New("Invalid intent")
)

// OnLaunch is function-type
func OnLaunch(launchRequest alexa.RequestDetail) (alexa.Response, error) {
	return GetWelcomeResponse(), nil
}

// GetWelcomeResponse is function-type
func GetWelcomeResponse() alexa.Response {

	cardTitle := " しりとり"
	speechOutput := GAME_START_MESSAGE
	repromptText := GAME_START_MESSAGE
	shouldEndSession := true
	return alexa.BuildResponse(alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession))
}

//OnIntent is function-type
func OnIntent(intentRequest alexa.RequestDetail) (alexa.Response, error) {
	log.Print(intentRequest.Intent)
	log.Print(intentRequest.Intent.Slots)

	if intentRequest.Intent.Name == "ShiritoriIntent" {

		log.Print(intentRequest.Intent.Slots["shiritoriword"].Value)
		return getShiritoriWord(intentRequest.Intent.Slots["shiritoriword"].Value)
	}
	return GetWelcomeResponse(), nil
	//  else
	//  return GetWelcomeResponse()
	// cardTitle := " しりとり"
	// speechOutput := GAME_START_MESSAGE
	// repromptText := GAME_START_MESSAGE
	// shouldEndSession := true
	//return alexa.BuildResponse(alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession)), nil
	//	return alexa.Response{}, ErrInvalidIntent
	// return alexa.BuildResponse(alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession)), nil
}

func getShiritoriWord(value string) (alexa.Response, error) {
	value = strings.TrimSpace(value)
	//文字列を分割して、末尾文字を取得
	arr := strings.Split(value, "")
	lastCharacter := arr[len(arr)-1]
	//データベースに登録
	db.PutWord(value, 1)
	log.Print(lastCharacter)
	//末尾文字を取得後データベースに参照、単語を取得して
	res := db.GetWordData(lastCharacter)

	//ユーザーに単語をお知らせ

	if res == "" {
		res = "負けました。"
	}

	log.Print(value + ": check")

	cardTitle := " しりとりインテント"
	speechOutput := res
	repromptText := res + timeData.GetNowTimeFormat()
	shouldEndSession := true

	return alexa.BuildResponse(alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession)), nil
}

func Handler(event alexa.Request) (alexa.Response, error) {

	eventRequestType := event.Request.Type
	// if event.Session.New {
	// 	return OnSessionStarted(map[string]string{"requestId": event.Request.RequestID}, event.Session)
	// } else
	if eventRequestType == "LaunchRequest" {
		return OnLaunch(event.Request)
	} else if eventRequestType == "IntentRequest" {
		return OnIntent(event.Request)
	}
	return alexa.Response{}, ErrInvalidIntent
}

func main() {

	lambda.Start(Handler)
}

func test() {

}
