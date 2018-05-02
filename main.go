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
	putGameInfo("LaunchRequest")
	//db.PutWord(value, 0)
	return alexa.BuildResponse(alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession))
}

//OnIntent is function-type
func OnIntent(intentRequest alexa.RequestDetail) (alexa.Response, error) {
	log.Print(intentRequest.Intent)
	log.Print(intentRequest.Intent.Name)

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
		answer := lastWord + ",echo;"
		historyID = db.PutHistoryDetailData(answer, constant.FIRST_GAME_FLAG)
	} else {
		//ゲーム開始第二回目以後フラグは1に変更
		if historyID != "" {
			db.UpdateHistoryFlag(constant.AFTER_GAME_FLAG, historyID)
		} else {
			//スキルを起動せずに しりとりゲーム開始した場合
			// lastWord = ""
			// historyID = db.PutHistoryDetailData(constant.FIRST_GAME_FLAG)
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
		errMes = constant.WRONG_END_WORD + "古い返答した単語は" + lastWord + "。最新返答した単語は" + firstCharacter
	} else if function.IsExistWord(useWord, value) {
		//ユーザー返答単語が重複しているかどうか
		errMes = constant.IS_EXIST_WORD
	} else {
		useWord = function.MakeDBAnswer(useWord, value, constant.ANSWERER_USER)
		fmt.Println("ユーザーの単語が問題ない " + useWord)
		db.UpdateHistoryDetailAnswer(useWord, historyID)
		// db.PutGameInfo(1)
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
	cardTitle := " しりとりインテント"
	if errMes != "" {
		speechOutput = errMes
	} else {
		speechOutput = value + constant.ANSWER_MSG + res
	}
	repromptText := res

	// go test()

	// var cd CountDown = CountDown{10, "残りは", "負け"}
	// aKun := make(chan *CountDown)
	// bSan := make(chan *CountDown)

	// var wg sync.WaitGroup
	// wg.Add(2)

	// go func() {
	// 	defer wg.Done()
	// 	for {
	// 		cd, ok := <-aKun // aKunから読み込み(書き込みを待つ)
	// 		if !ok {
	// 			break
	// 		}
	// 		time.Sleep(time.Second)
	// 		// fmt.Printf("　A君「%d！！！」\n", cd.Count)
	// 		cd.Count-- // データの書き換え
	// 		if cd.Count == 0 {
	// 			break
	// 		}
	// 		if cd.Count > 5 {
	// 			fmt.Printf("　A君「%d！！！」\n", cd.Count)
	// 		}
	// 		bSan <- cd // bSanへ書き込み
	// 		//}

	// 	}
	// 	close(bSan)
	// }()

	// go func() {
	// 	defer wg.Done()
	// 	speechOutput := ""
	// 	for {
	// 		cd, ok := <-bSan // bSanから読み込み(書き込みを待つ)
	// 		if !ok {
	// 			break
	// 		}
	// 		//			time.Sleep(time.Second)
	// 		if cd.Count <= 5 {

	// 			//ユーザーに返すレスポンス設定
	// 			cardTitle := " しりとりインテント"

	// 			speechOutput = "残りあとは" + strconv.Itoa(cd.Count) + "秒"
	// 			repromptText := "残りあとは" + strconv.Itoa(cd.Count) + "秒"
	// 			shouldEndSession := true
	// 			alexa.BuildResponse(alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession))

	// 			fmt.Printf("Bさん「%d！！！」\n", cd.Count)
	// 			//cd.Count-- // データの書き換え

	// 		}
	// 		if cd.Count == 0 {
	// 			break
	// 		}

	// 		aKun <- cd // aKunへ書き込み
	// 	}
	// 	close(aKun)
	// }()

	// aKun <- &cd // 最初の書き込み
	// wg.Wait()   // 終了を待つ
	// time.Sleep(time.Second)
	// fmt.Printf("A君・Bさん「「%s！！！」」\n", cd.LostWord)

	return alexa.BuildResponse(alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession)), nil
}

//ゲームを終了
func onCancelIntent() (alexa.Response, error) {
	return Cancel(), nil
}

func Cancel() alexa.Response {
	cardTitle := " しりとり"
	speechOutput := constant.GAME_STOP_MEESAGE
	repromptText := constant.GAME_STOP_MEESAGE
	shouldEndSession := true
	db.DeleteHistory(historyID)
	historyID = ""
	return alexa.BuildResponse(alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession))
}

func Handler(event alexa.Request) (alexa.Response, error) {

	eventRequestType := event.Request.Type

	// if event.Session.New {
	// 	return OnSessionStarted(map[string]string{"requestId": event.Request.RequestID}, event.Session)
	// } else
	if eventRequestType == "LaunchRequest" {
		return OnLaunch(event.Request)
	} else if eventRequestType == "IntentRequest" {
		intentName := event.Request.Intent.Name
		if intentName == "AMAZON.StopIntent" {
			//return onStopIntent()
		} else if intentName == "AMAZON.CancelIntent" {
			return onCancelIntent()
		}
		return OnIntent(event.Request)
	}
	return alexa.Response{}, ErrInvalidIntent
}

func main() {

	lambda.Start(Handler)
}

type CountDown struct {
	Count    int
	countMes string
	LostWord string
}

func test() {

}
