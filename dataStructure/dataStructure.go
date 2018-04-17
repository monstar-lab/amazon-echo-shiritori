package dataStructure

//DB
type Word struct {
	WordID int    `dynamo:"word_id"`
	Word   string `dynamo:"word"`
}

type History struct {
	HistoryID int    `dynamo:"history_id"`
	Time      string `dynamo:"time"`
	Flag      int    `dynamo:"flag"`
}

type HistoryDetail struct {
	HistoryDetailID int    `dynamo:"history_detail_id"`
	HistoryID       int    `dynamo:"history_id"`
	Time            string `dynamo:"time"`
	Answerer        string `dynamo:"Answerer"`
	answer          string `dynamo:"answer"`
}

type WordDB struct {
	WordID int    `json:"word_id" dynamodbav:"word_id"`
	Word   string `json:"word" dynamodbav:"word"`
}

type Words struct {
	WordID int    `json:"word_id"`
	Word   string `json:"word"`
}

//漢字をふりがなに変換するAPI
type XML struct {
	ResultSet ResultSet `xml: "ResultSet"`
}
type ResultSet struct {
	Result Result `xml: "Result"`
}

type Result struct {
	WordList WordList `xml: "WordList"`
}

type WordList struct {
	Word []WordXML `xml:"Word"`
}

type WordXML struct {
	Surface  string `xml:"Surface"`
	Furigana string `xml:"Furigana"`
	Roman    string `xml:"Roman"`
}
