## Amazon echo でしりとりする

- [要件定義](https://github.com/monstar-lab/amazon-echo-shiritori/wiki/%E8%A6%81%E4%BB%B6%E5%AE%9A%E7%BE%A9)

- 進捗  　
  ![進捗](https://user-images.githubusercontent.com/38127805/39162741-fc5ec182-47b1-11e8-8b56-b9f94407f462.png)

- テスト仕様書  
  [テスト仕様書](https://github.com/monstar-lab/amazon-echo-shiritori/wiki/images/test.png)

- 手順
	- Alexaのスキル名を呼びゲーム開始  
      「しりとりゲーム」
    
	- しりとりする仕方  
		- echo は最初の単語を出す
		- 単語の続きにスキル名を言う[では、次は、etc]単語を言う

	- 中断  
	 「しりとりゲーム中断」

	- 再開   
	 「しりとりゲーム再開」

	- 終了  
	 「しりとりゲーム終了」

- オプション
	- スクレイピング
	  しりとりサイトをスクレイピングして、単語を取得しDBに登録
	  [ソースコード](https://github.com/monstar-lab/amazon-echo-shiritori/tree/inprocess/getWordList)

- 今困っているところ
	- DynamoDB Scan
	  全件取得については、制限がある
	  制限の解決で全件取得する

	- 重複してないword listを返す
  
- バグ
	- しりとりゲーム開始の時、時々問題が発生する
	- 単語のチェックが必要（正しい単語かどうかチェック）

- シーケンス図 