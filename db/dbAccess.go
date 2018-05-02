package db

import (
	"encoding/json"
	"fmt"
	"strings"

	"../constant"
	"../dataStructure"
	"../timeData"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/guregu/dynamo"
)

//DBから単語を取得 何から始まる単語を word_v2
func GetDBWordList(keyword string) []string {
	cred := credentials.NewStaticCredentials(constant.ACCESS_KEY_ID, constant.SECRET_ACCESS_KEY, "") // 最後の引数は[セッショントークン]

	db := dynamodb.New(session.New(), &aws.Config{
		Credentials: cred,
		Region:      aws.String(constant.REGION), // constant.REGION等
	})

	getParams := &dynamodb.QueryInput{
		TableName: aws.String(constant.DB_WORD),

		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":wordID": {
				S: aws.String(keyword),
			},
			":word": { // :を付けるのがセオリーのようです
				S: aws.String(keyword),
			},
		},
		KeyConditionExpression: aws.String("word_id = :wordID and begins_with(word, :word)"),
	}

	getItem, getErr := db.Query(getParams)

	if getErr != nil {
		panic(getErr)
	}
	return makeWordList(getItem, keyword)

}

//単語全て取得きたら　単語だけのリストを作成
func makeWordList(output *dynamodb.QueryOutput, keyword string) []string {
	res := []string{}
	words := make([]*dataStructure.WordDB, 0)
	unMarshaListOfMapErr := dynamodbattribute.UnmarshalListOfMaps(output.Items, &words)
	if unMarshaListOfMapErr != nil {
		panic(fmt.Sprintf("failed to unmarshal Dynamodb Scan Items, %v", unMarshaListOfMapErr))
	}

	bytes, _ := json.Marshal(words)

	//変換されたデータ形をパースし、取得
	var data []*dataStructure.Words
	unMarshaErr := json.Unmarshal(bytes, &data)
	if unMarshaErr != nil {
		fmt.Println("error:", unMarshaErr)
	}

	for _, word := range data {
		fmt.Printf("word_id: %v, word: %v\n", word.WordID, word.Word)
		res = append(res, word.Word)
	}
	return res
}

//history_detail_v2から一回のゲームで回答した単語全て
func GetHistoryWord(historyDetailID string) string {
	cred := credentials.NewStaticCredentials(constant.ACCESS_KEY_ID, constant.SECRET_ACCESS_KEY, "") // 最後の引数は[セッショントークン]

	db := dynamodb.New(session.New(), &aws.Config{
		Credentials: cred,
		Region:      aws.String(constant.REGION), // constant.REGION等
	})

	getParams := &dynamodb.GetItemInput{
		TableName: aws.String("history_detail_v2"),

		Key: map[string]*dynamodb.AttributeValue{
			"history_id": {
				S: aws.String(historyDetailID),
			},
		},
		AttributesToGet: []*string{
			aws.String("answer"), // 欲しいデータの名前
		},
	}

	getItem, getErr := db.GetItem(getParams)

	if getErr != nil {
		panic(getErr)
	}
	return *getItem.Item["answer"].S
}

//DBから取得してきた単語と単語リストを比較し、重複してない単語最初の一件を返す
func ResNotUesWord(userWord string, wordList []string) string {
	//DBから取得きた単語を配列に変換
	oneRes := strings.Split(userWord, ";")
	for i := 0; i < len(oneRes)-1; i++ {
		oneWord := strings.Split(oneRes[i], ",")

		//重複した単語を削除
		for j := 0; j < len(oneWord)-1; j++ {
			wordList = delete_strings(wordList, oneWord[j])
		}
	}
	if len(wordList) != 0 {
		return wordList[0]
	}
	return ""
}

//重複したものを削除
func delete_strings(slice []string, s string) []string {
	ret := make([]string, len(slice))
	i := 0
	for _, x := range slice {
		if s != x {
			ret[i] = x
			i++
		}
	}
	return ret[:i]
}

//新規ゲーム登録
func PutHistoryDetailData(answer string, flag int) string {
	cred := credentials.NewStaticCredentials(constant.ACCESS_KEY_ID, constant.SECRET_ACCESS_KEY, "") // 最後の引数は[セッショントークン]

	db := dynamo.New(session.New(), &aws.Config{
		Credentials: cred,
		Region:      aws.String(constant.REGION), // constant.REGION等
	})

	table := db.Table(constant.DB_USE_WORD_HISTORY)
	//history idを作成
	historyID := timeData.GetNowTimeFormat(constant.DB_ID_FORMAT)

	history := dataStructure.HistoryDetail{HistoryID: historyID, Answer: answer, Flag: flag}

	if err := table.Put(history).Run(); err != nil {
		fmt.Println("err")
		panic(err.Error())
	}
	return historyID
}

//history_detail_v2 テーブルのフラグ変更
func UpdateHistoryFlag(flag int, historyID string) {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	db := dynamodb.New(sess, &aws.Config{
		Credentials: credentials.NewStaticCredentials(constant.ACCESS_KEY_ID, constant.SECRET_ACCESS_KEY, ""),
		Region:      aws.String(constant.REGION), // constant.REGION等
	})
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":flag": {
				N: aws.String(fmt.Sprintf("%v", flag)),
			},
		},
		TableName: aws.String(constant.DB_USE_WORD_HISTORY),
		Key: map[string]*dynamodb.AttributeValue{
			"history_id": {
				S: aws.String(historyID),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set flag = :flag"),
	}

	_, err = db.UpdateItem(input)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

//history_detail_v2 のanswerを変更
func UpdateHistoryDetailAnswer(answer string, historyID string) {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	db := dynamodb.New(sess, &aws.Config{
		Credentials: credentials.NewStaticCredentials(constant.ACCESS_KEY_ID, constant.SECRET_ACCESS_KEY, ""),
		Region:      aws.String(constant.REGION), // constant.REGION等
	})
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":answer": {
				S: aws.String(answer),
			},
		},
		TableName: aws.String(constant.DB_USE_WORD_HISTORY),
		Key: map[string]*dynamodb.AttributeValue{
			"history_id": {
				S: aws.String(historyID),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set answer = :answer"),
	}
	_, err = db.UpdateItem(input)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func GetGameStartWord(wordID int) string {
	cred := credentials.NewStaticCredentials(constant.ACCESS_KEY_ID, constant.SECRET_ACCESS_KEY, "") // 最後の引数は[セッショントークン]
	db := dynamo.New(session.New(), &aws.Config{
		Credentials: cred,
		Region:      aws.String(constant.REGION), // constant.REGION等
	})
	table := db.Table("word")
	result := []dataStructure.Word{}
	err := table.Get("word_id", wordID).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	return result[0].Word
}

func DeleteHistory(historyID string) {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	db := dynamodb.New(sess, &aws.Config{
		Credentials: credentials.NewStaticCredentials(constant.ACCESS_KEY_ID, constant.SECRET_ACCESS_KEY, ""),
		Region:      aws.String(constant.REGION), // constant.REGION等
	})

	params := &dynamodb.DeleteItemInput{
		TableName: aws.String(constant.DB_USE_WORD_HISTORY), // テーブル名

		Key: map[string]*dynamodb.AttributeValue{
			"history_id": { // キー名
				S: aws.String(historyID), // 削除するキーの値
			},
		},
		// //返ってくるデータの種類
		// ReturnConsumedCapacity:      aws.String("NONE"),
		// ReturnItemCollectionMetrics: aws.String("NONE"),
		// ReturnValues:                aws.String("NONE"),
	}

	_, err = db.DeleteItem(params)
	//fmt.Println(res)
	if err != nil {
		fmt.Println(err.Error())
	}
}
