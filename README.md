## Amazon echo でしりとりする

- 仕様
	- スキル名  
     しりとりゲーム
     
   - しりとりする仕方  
   		- echo は最初の単語を出す
  		- 単語の続きにスキル名を言う[では、次は、etc]単語を言う
   
   		
    
- 実装予定
	- チェック	 
		- 重複チェック
		- 「ん」単語をチェック
		- 語尾が正しいかどうかのチェック 
	
	- DB
		- 返答単語の登録
		- しりとりマスタテーブルのデータ自動更新

	- お知らせ    
		- カウントダウン  
		- 自動終了知らせ
		- 間違い知らせ

- ライセンス
	- ヤフーの漢字をひらがなに変換するAPIを使用しています。[URL](https://developer.yahoo.co.jp/webapi/jlp/jim/v1/conversion.html)
	- [ライセンス詳細](https://developer.yahoo.co.jp/appendix/rate.html)