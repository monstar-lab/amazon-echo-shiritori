package getWordList

import (
	"fmt"

	"../constant"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

func Insert() {
	//取得してきた単語をデータベースに登録
	//下記コメントを外し、呼び出すと実行される
	//DBにおよそ57,600回のアクセス数になる！！！
	//putWord(getWordData.GetAllData())
}

type WordDB struct {
	WordID int    `dynamo:"word_id"`
	Word   string `dynamo:"word"`
}

func putWord(arr []string) {
	cred := credentials.NewStaticCredentials(constant.ACCESS_KEY_ID, constant.SECRET_ACCESS_KEY, "") // 最後の引数は[セッショントークン]
	db := dynamo.New(session.New(), &aws.Config{
		Credentials: cred,
		Region:      aws.String("ap-northeast-1"), // "ap-northeast-1"等
	})
	table := db.Table("word")
	for i, word := range arr {
		word := WordDB{WordID: i + 1, Word: word}
		fmt.Println(word)

		err := table.Put(word).Run()
		if err != nil {
			fmt.Println("err")
			panic(err.Error())
		}
	}
}
