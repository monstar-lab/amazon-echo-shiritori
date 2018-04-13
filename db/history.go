package db

import (
	"fmt"

	"../dataStructure"
	"../timeData"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

func Insert() {
	cred := credentials.NewStaticCredentials(ACCESS_KET_ID, SECRET_ACCESS_KEY, "") // 最後の引数は[セッショントークン]

	db := dynamo.New(session.New(), &aws.Config{
		Credentials: cred,
		Region:      aws.String("ap-northeast-1"), // "ap-northeast-1"等
	})

	table := db.Table("history")

	history := dataStructure.History{HistoryID: 1, Time: timeData.GetNowTimeFormat(), Flag: 3}
	//u := User{User_ID: "lambda test"}
	fmt.Println(history)

	if err := table.Put(history).Run(); err != nil {
		fmt.Println("err")
		panic(err.Error())
	}
}

// func main() {

// 	lambda.Start(Handler)
// }
