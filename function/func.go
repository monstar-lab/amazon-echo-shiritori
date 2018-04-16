package function

import (
	"encoding/json"
	"fmt"
	"strings"

	"../dataStructure"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func test() {

}

// func ResCount(output *dynamodb.ScanOutput) {
// 	count

// }
//ユーザ返答単語を返す
func ResWord(output *dynamodb.ScanOutput, keyword string) string {
	// DBから取得したデータのJSONの形を変換
	words := make([]*dataStructure.WordDB, 0)
	unMarshaListOfMapErr := dynamodbattribute.UnmarshalListOfMaps(output.Items, &words)
	if unMarshaListOfMapErr != nil {
		panic(fmt.Sprintf("failed to unmarshal Dynamodb Scan Items, %v", unMarshaListOfMapErr))
	}

	bytes, _ := json.Marshal(words)

	//変換されたデータ形をパースし、取得
	var data []dataStructure.Words
	unMarshaErr := json.Unmarshal(bytes, &data)
	if unMarshaErr != nil {
		fmt.Println("error:", unMarshaErr)
	}

	for _, word := range data {
		fmt.Printf("word_id: %v, word: %v\n", word.WordID, word.Word)
		isWord := CheckWord(word.Word, keyword)
		if isWord == true {
			return word.Word
		}
	}
	return ""
}

//文字列の先頭部分は末尾文字と一致しているかどうか
func CheckWord(value string, keyword string) bool {

	//fmt.Println(strings.HasPrefix("ナツ", "ツナミ"))
	fmt.Println(strings.HasPrefix(value, keyword))
	return strings.HasPrefix(value, keyword)
}

//「ん」のチェック
func CheckN(str string) bool {
	if str == "ん" || str == "ン" {
		fmt.Println("yes")
		return true
	}
	return false
}

//末尾文字の違いチェック
func CheckEndOfTheWordIsWrong(firstCharacter string, lastCharacter string) bool {
	if lastCharacter != "" {
		arr := strings.Split(lastCharacter, "")

		last := arr[len(arr)-1]

		if firstCharacter == last {
			return false
		} else {
			return true
		}

	} else {
		return false
	}

}
