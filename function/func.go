package function

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"../constant"
	"../dataStructure"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

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

//APIからひらがなのデータを取得
func GetAPIData(word string) string {
	//APIアクセスURL
	url := constant.URL + "?appid=" + constant.API_ACCESS_ID + "&sentence=" + word
	data := httpGet(url)

	//データセット
	result := dataStructure.ResultSet{}
	err := xml.Unmarshal([]byte(data), &result)
	furigana := ""
	if err != nil {
		fmt.Println("error: %v", err)
		return furigana
	}

	for _, word := range result.Result.WordList.Word {
		furigana += word.Furigana
		fmt.Println(word.Furigana)
		fmt.Println(furigana)
	}
	return furigana
}

//APIと通信
func httpGet(url string) string {
	response, _ := http.Get(url)
	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	return string(body)
}

//ランダムにword_idを出す
func RandWordID() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(11655)
}

var hirakanaList = map[rune]string{
	'ぁ': "あ",
	'ぃ': "い",
	'ぅ': "う",
	'ぇ': "え",
	'ぉ': "お",
	'っ': "つ",
	'ゃ': "や",
	'ゅ': "ゆ",
	'ょ': "よ",
	'ゎ': "わ",
}

func HiraToZenHira(value string) string {
	res := value
	for i, v := range hirakanaList {
		if string(i) == value {
			res = v
		}
	}
	return res
}

//末尾文字を返す
func IsLongVowel(value string) string {
	runes := []rune(value)
	strlen := utf8.RuneCountInString(value)
	if string(runes[strlen-1:strlen]) == "ー" {
		return string(runes[strlen-2 : strlen-1])
	}
	return string(runes[strlen-1 : strlen])
}

//返答単語を渡す
func ResLastCharacter(value string) string {
	lastCharacter := IsLongVowel(value)
	lastCharacter = HiraToZenHira(lastCharacter)
	return lastCharacter
}

func IsExistWord(count int) bool {
	if count >= 1 {
		return true
	}
	return false
}
