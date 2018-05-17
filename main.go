package main

import (
	"errors"
	"fmt"
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
	stringChan := make(chan string)
	go func() {
		res := db.GetGameStartWord(function.RandWordID())
		log.Println(res)
		stringChan <- res
	}()
	return <-stringChan
}

// ゲームスタート
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
	shouldEndSession := false
	//DBにゲーム開始情報登録
	answer := lastWord + ",echo;"
	historyID = db.PutHistoryDetailData(answer, constant.FIRST_GAME_FLAG)
	return alexa.BuildResponse(alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession))
}

//intent開始
func OnIntent(intentRequest alexa.RequestDetail) (alexa.Response, error) {
	if intentRequest.Intent.Name == "ShiritoriIntent" {
		if intentRequest.Intent.Slots["shiritoriword"].Value == "" {
			return ResNotValue(), nil
		}
		return getShiritoriWord(intentRequest.Intent.Slots["shiritoriword"].Value)
	}
	return GetWelcomeResponse(), nil
}

func ResNotValue() alexa.Response {
	cardTitle := " しりとり"
	speechOutput := constant.NOT_FOUND_VALUE
	repromptText := constant.NOT_FOUND_VALUE
	shouldEndSession := false
	return alexa.BuildResponse(alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession))
}

func getShiritoriWord(value string) (alexa.Response, error) {
	//ユーザー返答した単語をAPIと通信、ひらがなの取得
	value = function.GetAPIData(value)
	fmt.Println("返答単語　" + value)
	//空白文字を削除
	value = strings.TrimSpace(value)
	//末尾文字を取得
	lastCharacter := function.ResLastCharacter(value)
	fmt.Println("lastCharacter" + lastCharacter)
	fmt.Println("lastWord " + function.ResLastCharacter(lastWord))
	//始まり文字を取得
	firstCharacter := string([]rune(value)[:1])

	//各変数の初期値を設定する
	//ユーザに返すレスポンス
	res := ""
	//エラーメッセージ
	errMes := ""
	//ユーザーに返答するメッセージ
	speechOutput := ""
	//今まで返答した単語
	useWord := db.GetHistoryWord(historyID)
	shouldEndSession := false
	//末尾チェック
	if function.CheckN(lastCharacter) {
		//末尾が「ん」
		errMes = constant.LOSS_N_MESSAGE
		db.DeleteHistory(historyID)
		shouldEndSession = true
	} else if function.CheckEndOfTheWordIsWrong(firstCharacter, function.ResLastCharacter(lastWord)) == true {
		//末尾が違う
		errMes = constant.WRONG_END_MES + function.ResLastCharacter(lastWord) + constant.WRONG_END_WORD_MES
	} else if function.IsExistWord(useWord, value) {
		//ユーザー返答単語が重複しているかどうか
		errMes = constant.IS_EXIST_WORD
	} else {
		useWord = function.MakeDBAnswer(useWord, value, constant.ANSWERER_USER)
		db.UpdateHistoryDetailAnswer(useWord, historyID)
		//末尾文字を取得後データベースに参照、単語を取得して
		res = db.ResNotUesWord(db.GetHistoryWord(historyID), db.GetDBWordList(lastCharacter))
		//最後に返答した単語を値保持
		lastWord = res
		//ユーザーに結果をお知らせ
		if res == "" {
			errMes = constant.LOSS_GAME
			db.DeleteHistory(historyID)
			shouldEndSession = true
		} else {
			useWord = function.MakeDBAnswer(useWord, res, constant.ANSWERER_ECHO)
			db.UpdateHistoryDetailAnswer(useWord, historyID)
		}
	}
	//ユーザーに返すレスポンス設定
	cardTitle := " shiritoriIntent"
	if errMes != "" {
		speechOutput = errMes
	} else {
		speechOutput = value + constant.ANSWER_MSG + res
	}
	repromptText := "user firstCharacter " + firstCharacter + " lastword " + lastWord
	return alexa.BuildResponse(alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession)), nil
}

//ゲームを中止
func onStopIntent() (alexa.Response, error) {
	return Stop(), nil
}

func Stop() alexa.Response {
	db.UpdateHistoryDetailFlag(historyID, constant.STOP_GAME_FLAG)
	cardTitle := " stop"
	speechOutput := constant.GAME_STOP_MEESAGE
	repromptText := constant.GAME_STOP_MEESAGE
	shouldEndSession := false
	return alexa.BuildResponse(alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession))
}

//ゲームを終了
func onCancelIntent() (alexa.Response, error) {
	return Cancel(), nil
}

func Cancel() alexa.Response {
	cardTitle := " cancel"
	speechOutput := constant.GAME_CANCEL_MEESAGE
	repromptText := constant.GAME_STOP_MEESAGE
	shouldEndSession := true
	db.DeleteHistory(historyID)
	historyID = ""
	return alexa.BuildResponse(alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession))
}

func GetResumeData() []string {
	getFlagData := make(chan []string)
	resumeID := make(chan string)
	getData := make(chan string)
	go func() {
		res := db.GetFlagData(constant.STOP_GAME_FLAG)
		log.Println(res)
		getFlagData <- res
	}()
	getOldData := <-getFlagData
	go func() {
		id := db.GetResumeData(getOldData)

		resumeID <- id
	}()
	id := <-resumeID
	go func() {

		word := function.GetHistoryLastWord(db.GetHistoryWord(id))
		fmt.Println(word)
		getData <- word
	}()
	resumeWord := <-getData
	return []string{id, resumeWord}
}

func onResumeIntent() (alexa.Response, error) {
	//再開したゲームのIDと最後返答した単語を取得
	resume := GetResumeData()
	// ゲーム再開した場合今開始したゲームを削除
	db.DeleteHistory(historyID)

	historyID = resume[0]
	lastWord = resume[1]
	return ResumeIntent(function.GetHistoryLastWord(resume[1])), nil
}

func ResumeIntent(word string) alexa.Response {

	cardTitle := "resume"
	speechOutput := constant.GAME_RESUME_MEESAGE + word
	repromptText := constant.GAME_RESUME_MEESAGE + word
	shouldEndSession := false
	return alexa.BuildResponse(alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession))
}

func Handler(event alexa.Request) (alexa.Response, error) {
	if !event.Session.New {
		oldGameID := db.GetFlagData(constant.FIRST_GAME_FLAG)
		for _, id := range oldGameID {
			if id != historyID {
				db.DeleteHistory(id)
			}
		}
	}
	eventRequestType := event.Request.Type
	if eventRequestType == "LaunchRequest" {
		return OnLaunch(event.Request)
	} else if eventRequestType == "IntentRequest" {
		intentName := event.Request.Intent.Name
		if intentName == "AMAZON.StopIntent" {
			return onStopIntent()
		} else if intentName == "AMAZON.CancelIntent" {
			return onCancelIntent()
		} else if intentName == "resumeIntent" {
			return onResumeIntent()
		} else if intentName == "newstartIntent" {
			db.DeleteHistory(historyID)
			return OnLaunch(event.Request)
		}
		return OnIntent(event.Request)
	} else if eventRequestType == "SessionEndedRequest" {

	}
	return alexa.Response{}, ErrInvalidIntent
}

func main() {
	lambda.Start(Handler)
}
