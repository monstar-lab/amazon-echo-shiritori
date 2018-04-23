## Amazon echo でしりとりする

- [要件定義](https://github.com/monstar-lab/amazon-echo-shiritori/wiki/%E8%A6%81%E4%BB%B6%E5%AE%9A%E7%BE%A9)

- 進捗　
![4月23日進捗](https://github.com/monstar-lab/amazon-echo-shiritori/files/1936561/default.pdf "進捗")

- 手順
	- Alexaのスキル名を呼びゲーム開始  
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
	   		  
- 未解決
	- DynamoDB Scan
	  全件取得については、制限がある
	  制限の解決で全件取得する

	- 全件取得テスト
  