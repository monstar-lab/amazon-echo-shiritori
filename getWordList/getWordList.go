package getWordList

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//取得したデータを配列で返す
func GetAllData() []string {
	arr := []string{}
	doc, err := goquery.NewDocument("http://siritori.net/line/")
	if err != nil {
		fmt.Print("url scarapping failed")
	}
	//全ての行のURLを取得
	doc.Find(".box > ul > li > a").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		if strings.HasPrefix(url, "/line/") {
			arr = append(arr, url)
		}
	})
	return get(arr)
}

//行ごとに取得したデータをまとめる
func get(arr []string) []string {
	words := []string{}
	lineWords := []string{}
	firstURL := ""
	for _, url := range arr {
		firstURL = "http://siritori.net" + url
		lineWords = line(getNumbers(firstURL), firstURL)
		words = append(words, lineWords...)
	}
	return words
}

//行ごとのページ数取得
func getNumbers(url string) []string {
	arr := []string{}
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Print("url scarapping failed")
	}
	doc.Find(".paginate > .numbers").Each(func(i int, s *goquery.Selection) {
		s.Find(".disabled").Remove()
		s.Find(".next").Remove()
		s.Find("span > a").Each(func(j int, m *goquery.Selection) {
			url, _ := m.Attr("href")
			arr = append(arr, url)
		})
	})
	return arr
}

//行ごとのloop
func line(lineURLS []string, firstURL string) []string {
	words := []string{}
	words = append(words, getWord(firstURL)...)
	for _, url := range lineURLS {
		words = append(words, getWord("http://siritori.net"+url)...)
	}
	return words
}

//行ごとの単語
func getWord(url string) []string {
	words := []string{}
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Print("url scarapping failed")
	}
	doc.Find(".pages > .linkClound > li > a").Each(func(i int, s *goquery.Selection) {
		word := s.Text()
		words = append(words, word)
	})
	return words
}
