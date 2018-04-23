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

//ゲーム履歴ID
var historyID = ""

// OnLaunch is function-type
func OnLaunch(launchRequest alexa.RequestDetail) (alexa.Response, error) {
	return GetWelcomeResponse(), nil
}

//ゲーム開始時のスタート単語を取得
func GetWelcomeStartWord() string {
	return db.GetGameStartWord(function.RandWordID())
}

// GetWelcomeResponse is function-type
func GetWelcomeResponse() alexa.Response {

	//スタート単語取得
	startWord := GetWelcomeStartWord()
	//ゲーム最新単語の更新
	lastWord = startWord
	//ユーザー最後返答する単語を格納する変数を初期値に戻す
	historyID = ""
	cardTitle := " しりとり"
	speechOutput := constant.GAME_START_MESSAGE + startWord
	repromptText := constant.GAME_START_MESSAGE
	shouldEndSession := true
	//DBにゲーム開始情報登録
	putGameInfo("LaunchRequest")
	//db.PutWord(value, 0)
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

//ゲーム開始情報登録
func putGameInfo(intent string) {
	//ゲーム開始第一回目フラグは0に登録
	if historyID == "" && intent == "LaunchRequest" {
		historyID = db.PutGameInfo(constant.FIRST_GAME_FLAG)
	} else {
		//ゲーム開始第二回目以後フラグは1に変更
		if historyID != "" {
			db.UpdateItem(constant.AFTER_GAME_FLAG, historyID)
		} else {
			lastWord = ""
			historyID = db.PutGameInfo(constant.FIRST_GAME_FLAG)
		}
	}
}

func getShiritoriWord(value string) (alexa.Response, error) {
	//DBにゲーム開始情報登録
	putGameInfo("ShiritoriIntent")
	//ユーザー返答した単語をAPIと通信、ひらがなの取得
	value = function.GetAPIData(value)
	//空白文字を削除
	value = strings.TrimSpace(value)
	//末尾文字を取得
	lastCharacter := function.ResLastCharacter(value)

	//始まり文字を取得
	firstCharacter := string([]rune(value)[:1])

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
	} else if function.IsExistWord(db.SearchWordCount(historyID, value)) {
		//ユーザー返答単語が重複しているかどうか
		errMes = constant.IS_EXIST_WORD
	} else {
		db.PutWord(value, historyID, constant.ANSWERER_USER, constant.NOT_LAST_ANSWERER)
		// db.PutGameInfo(1)
		log.Print(lastCharacter)
		log.Print(firstCharacter)
		//末尾文字を取得後データベースに参照、単語を取得して
		res = db.GetWordData(lastCharacter)
		//最後に返答した単語を値保持
		lastWord = res
		//ユーザーに結果をお知らせ
		if res == "" {
			errMes = constant.LOSS_GAME
			db.UpdateItem(constant.END_GAME_FLAG, historyID)
		} else {
			db.PutWord(res, historyID, constant.ANSWERER_ECHO, constant.LAST_ANSWERER)
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
