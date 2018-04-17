package main

import (
	"errors"
	"log"
	"strings"

	"./alexa"
	"./constant"
	"./db"
	"./function"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	ErrInvalidIntent = errors.New("Invalid intent")
)

//最後に返答した単語を格納
var lastWord = ""

// OnLaunch is function-type
func OnLaunch(launchRequest alexa.RequestDetail) (alexa.Response, error) {
	return GetWelcomeResponse(), nil
}

// GetWelcomeResponse is function-type
func GetWelcomeResponse() alexa.Response {

	//ユーザー最後返答する単語を格納する変数を初期値に戻す
	lastWord = ""
	cardTitle := " しりとり"
	speechOutput := constant.GAME_START_MESSAGE
	repromptText := constant.GAME_START_MESSAGE
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
}

func getShiritoriWord(value string) (alexa.Response, error) {
	//ユーザー返答した単語をAPIと通信、ひらがなの取得
	value = function.GetAPIData(value)
	//空白文字を削除
	value = strings.TrimSpace(value)
	//文字列を分割して、末尾文字を取得
	arr := strings.Split(value, "")
	lastCharacter := arr[len(arr)-1]
	//始まり文字を取得
	firstCharacter := arr[0]

	//各変数の初期値を設定する
	//ユーザに返すレスポンス
	res := ""
	//エラーメッセージ
	errMes := ""
	//ユーザーに返答するメッセージ
	speechOutput := ""

	//末尾チェック
	if function.CheckN(lastCharacter) {
		//末尾が「ん」
		errMes = constant.LOSS_N_MESSAGE
	} else if function.CheckEndOfTheWordIsWrong(firstCharacter, lastWord) == true {
		//末尾が違う
		errMes = constant.WRONG_END_WORD
	} else {
		//データベースに登録
		db.PutWord(value, 1)
		log.Print(lastCharacter)
		log.Print(firstCharacter)
		//末尾文字を取得後データベースに参照、単語を取得して
		res = db.GetWordData(lastCharacter)
		//最後に返答した単語を値保持
		lastWord = res
		//ユーザーに単語をお知らせ
		if res == "" {
			errMes = constant.LOSS_GAME
		}
		log.Print(value + ": check")
	}

	//ユーザーに返すレスポンス設定
	cardTitle := " しりとりインテント"
	if errMes != "" {
		speechOutput = errMes
	} else {
		speechOutput = value + constant.ANSWER_MSG + res
	}
	repromptText := res
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
