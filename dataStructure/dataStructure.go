package dataStructure

//DB
type Word struct {
	WordID int    `dynamo:"word_id"`
	Word   string `dynamo:"word"`
}

type HistoryDB struct {
	HistoryID string `json:"historyID" dynamodbav:"history_id"`
}

type History struct {
	HistoryID string `json:"historyID"`
}

type HistoryDetail struct {
	HistoryID string `dynamo:"history_id"`
	Answer    string `dynamo:"answer"`
	Flag      int    `dynamo:"flag"`
}

type WordDB struct {
	WordID string `json:"word_id" dynamodbav:"word_id"`
	Word   string `json:"word" dynamodbav:"word"`
}

type Words struct {
	WordID string `json:"word_id"`
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
