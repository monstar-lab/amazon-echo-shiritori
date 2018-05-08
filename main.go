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
	//time.Sleep(300 * time.Millisecond)
	stringChan := make(chan string)
	go func() {
		res := db.GetGameStartWord(function.RandWordID())
		log.Println(res)
		stringChan <- res
	}()

	log.Println(stringChan)
	return <-stringChan
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
	shouldEndSession := false
	//DBにゲーム開始情報登録
	answer := lastWord + ",echo;"
	historyID = db.PutHistoryDetailData(answer)
	return alexa.BuildResponse(alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession))
}

//OnIntent is function-type
func OnIntent(intentRequest alexa.RequestDetail) (alexa.Response, error) {
	log.Print(intentRequest.Intent)
	log.Print(intentRequest.Intent.Name)

	if intentRequest.Intent.Name == "ShiritoriIntent" {

		log.Print(intentRequest.Intent.Slots["shiritoriword"].Value)
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

// //ゲーム開始情報登録
// func putGameInfo(intent string) {
// 	//ゲーム開始第一回目フラグは0に登録
// 	if historyID == "" && intent == "LaunchRequest" {
// 		answer := lastWord + ",echo;"
// 		historyID = db.PutHistoryDetailData(answer, constant.FIRST_GAME_FLAG)
// 	} else {
// 		//ゲーム開始第二回目以後フラグは1に変更
// 		if historyID != "" {
// 			//db.UpdateHistoryFlag(constant.AFTER_GAME_FLAG, historyID)
// 		} else {
// 			//スキルを起動せずに しりとりゲーム開始した場合
// 			// lastWord = ""
// 			// historyID = db.PutHistoryDetailData(constant.FIRST_GAME_FLAG)
// 		}
// 	}
// }

func getShiritoriWord(value string) (alexa.Response, error) {
	// //DBにゲーム開始情報登録
	// putGameInfo("ShiritoriIntent")

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
		fmt.Println("ユーザーの単語が問題ない " + useWord)
		db.UpdateHistoryDetailAnswer(useWord, historyID)

		log.Print(lastCharacter)
		log.Print(firstCharacter)
		//末尾文字を取得後データベースに参照、単語を取得して
		res = db.ResNotUesWord(db.GetHistoryWord(historyID), db.GetDBWordList(lastCharacter))
		//最後に返答した単語を値保持
		lastWord = res
		//ユーザーに結果をお知らせ
		if res == "" {
			errMes = constant.LOSS_GAME
			db.DeleteHistory(historyID)
			shouldEndSession = true
			//db.UpdateHistoryFlag(constant.END_GAME_FLAG, historyID)
		} else {
			useWord = function.MakeDBAnswer(useWord, res, constant.ANSWERER_ECHO)
			fmt.Println("echoが返答した単語　 " + useWord)
			db.UpdateHistoryDetailAnswer(useWord, historyID)
		}
		log.Print(value + ": check")
	}

	//ユーザーに返すレスポンス設定
	cardTitle := " shiritoriIntent"
	if errMes != "" {
		speechOutput = errMes
	} else {
		speechOutput = value + constant.ANSWER_MSG + res
	}
	repromptText := res
	return alexa.BuildResponse(alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession)), nil
}

//ゲームを中止
func onStopIntent() (alexa.Response, error) {
	return Stop(), nil
}

func Stop() alexa.Response {
	cardTitle := " stop"
	speechOutput := constant.GAME_STOP_MEESAGE
	repromptText := constant.GAME_STOP_MEESAGE
	shouldEndSession := false
	//db.UpdateHistoryFlag(constant.STOP_GAME_FLAG, historyID)
	//historyID = ""
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

// func OnSessionStarted() (alexa.Response, error) {
// 	cardTitle := " sessionstarted"
// 	speechOutput := "ある" + historyID + lastWord
// 	repromptText := "ある" + historyID + lastWord
// 	shouldEndSession := false
// 	return alexa.BuildResponse(alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession)), nil
// }

func onResumeIntent() (alexa.Response, error) {
	cardTitle := "resume"
	speechOutput := constant.GAME_RESUME_MEESAGE + lastWord
	repromptText := constant.GAME_RESUME_MEESAGE + lastWord
	shouldEndSession := false
	return alexa.BuildResponse(alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession)), nil
}

func Handler(event alexa.Request) (alexa.Response, error) {

	eventRequestType := event.Request.Type

	// if eventRequestType == "SessionEndedRequest" {
	// 	return Test()
	// } else
	if eventRequestType == "LaunchRequest" {
		// if !event.Session.New {
		// 	return OnSessionStarted()
		// }
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
	}
	return alexa.Response{}, ErrInvalidIntent
}

func main() {

	lambda.Start(Handler)
}

// type CountDown struct {
// 	Count    int
// 	countMes string
// 	LostWord string
// }

func Test() (alexa.Response, error) {
	cardTitle := " SessionEndedRequest"
	speechOutput := "SessionEndedRequest" + historyID + lastWord
	repromptText := "SessionEndedRequest" + historyID + lastWord
	shouldEndSession := false

	//historyIDとlastwordが存在していた場合
	//
	//dbにアクセス 2
	return alexa.BuildResponse(alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession)), nil
}
