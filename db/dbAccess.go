package db

import (
	"fmt"
	"log"

	"../constant"
	"../dataStructure"
	"../function"
	"../timeData"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
)

//DBから単語を取得
func GetWordData(keyword string) string {
	log.Print(keyword)
	cred := credentials.NewStaticCredentials(constant.ACCESS_KEY_ID, constant.SECRET_ACCESS_KEY, "") // 最後の引数は[セッショントークン]

	db := dynamodb.New(session.New(), &aws.Config{
		Credentials: cred,
		Region:      aws.String(constant.REGION), // constant.REGION等
	})

	getParams := &dynamodb.ScanInput{
		TableName: aws.String("word"),

		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":wordID": {
				N: aws.String("1"),
			},
			":word": { // :を付けるのがセオリーのようです
				S: aws.String(keyword),
			},
		},
		// word_id が 1以上と「keyword」から始まりの単語
		FilterExpression: aws.String("word_id >= :wordID and begins_with(word, :word)"),
	}

	getItem, getErr := db.Scan(getParams)

	if getErr != nil {
		panic(getErr)
	}
	fmt.Println(getItem)

	return function.ResWord(getItem, keyword)

}

func PutGameInfo(flag int) string {
	cred := credentials.NewStaticCredentials(constant.ACCESS_KEY_ID, constant.SECRET_ACCESS_KEY, "") // 最後の引数は[セッショントークン]

	db := dynamo.New(session.New(), &aws.Config{
		Credentials: cred,
		Region:      aws.String(constant.REGION), // constant.REGION等
	})

	table := db.Table("history")

	//history idを作成
	historyID := timeData.GetNowTimeFormat(constant.DB_ID_FORMAT)

	history := dataStructure.History{HistoryID: historyID, Time: timeData.GetNowTimeFormat(constant.DB_INSERT_TIME_FORMAT), Flag: flag}
	//u := User{User_ID: "lambda test"}
	fmt.Println(history)

	if err := table.Put(history).Run(); err != nil {
		fmt.Println("err")
		panic(err.Error())
	}
	return historyID
}

//単語をデータベースに登録
func PutWord(word string, historyID string, answerer string, flag int) {
	log.Println("db -> " + word)
	cred := credentials.NewStaticCredentials(constant.ACCESS_KEY_ID, constant.SECRET_ACCESS_KEY, "") // 最後の引数は[セッショントークン]

	db := dynamo.New(session.New(), &aws.Config{
		Credentials: cred,
		Region:      aws.String(constant.REGION), // constant.REGION等
	})

	table := db.Table("history_detail")

	historyDetail := dataStructure.HistoryDetail{HistoryDetailID: timeData.GetNowTimeFormat(constant.DB_ID_FORMAT), HistoryID: historyID, Time: timeData.GetNowTimeFormat(constant.DB_INSERT_TIME_FORMAT), Answerer: answerer, Answer: word, Flag: flag}
	fmt.Println(historyDetail)

	err := table.Put(historyDetail).Run()
	if err != nil {
		fmt.Println("err")
		panic(err.Error())
	}
}

//historyテーブルのフラグ変更
func UpdateItem(flag int, historyID string) {
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
		TableName: aws.String("history"),
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

func SearchWordCount(historyID string, word string) int {
	cred := credentials.NewStaticCredentials(constant.ACCESS_KEY_ID, constant.SECRET_ACCESS_KEY, "") // 最後の引数は[セッショントークン]

	db := dynamodb.New(session.New(), &aws.Config{
		Credentials: cred,
		Region:      aws.String(constant.REGION), // constant.REGION等
	})

	getParams := &dynamodb.ScanInput{
		TableName: aws.String("history_detail"),

		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":historyID": {
				S: aws.String(historyID),
			},
			":answer": { // :を付けるのがセオリーのようです
				S: aws.String(word),
			},
		},
		// word_id が 1以上と「keyword」から始まりの単語
		FilterExpression: aws.String("history_id = :historyID and answer = :answer"),
	}

	getItem, getErr := db.Scan(getParams)

	if getErr != nil {
		panic(getErr)
	}

	return len(getItem.Items)
}

/* データベース上に登録されたMax IDを取得
* Scanは1MBのデータしか取得できない
* 解決する必要がる
* LastEvaluatedKeyを利用
 */
// func GetMaxID() {
// 	sess, err := session.NewSession()
// 	if err != nil {
// 		panic(err)
// 	}

// 	db := dynamodb.New(sess, &aws.Config{
// 		Credentials: credentials.NewStaticCredentials(constant.ACCESS_KEY_ID, constant.SECRET_ACCESS_KEY, ""),
// 		Region:      aws.String(constant.REGION), // constant.REGION等
// 	})

// 	getParams := &dynamodb.ScanInput{
// 		TableName: aws.String("word"),
// 		Select:    aws.String("COUNT"),
// 	}
// 	getItem, getErr := db.Scan(getParams)
// 	if getErr != nil {
// 		fmt.Println(getErr)
// 		return
// 	}

// 	fmt.Println(string(getItem.Count))
// }
