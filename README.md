## GoでAlexaスキル作成サンプル

### 概要
1. ユーザーがテストスキルを起動
2. Alexaが「名前は」と聞いてくれる
3. ユーザーが名前を答えたら
4. Hello「名前」と挨拶してスキル終了

### サンプルを動かす手順
- git clone (ブランチはgoSample)

- AWSに関数作成
	- `一から関数を作成`  
		<img src="https://github.com/monstar-lab/amazon-echo-shiritori/wiki/images/AWS/aws1.png" width= "400px"/>

	- `Alexa Skills Kitトリガーを追加`  
		<img src="https://github.com/monstar-lab/amazon-echo-shiritori/wiki/images/AWS/aws2.png" width= "400px"/>

	- `ARNをコピー`  
		<img src="https://github.com/monstar-lab/amazon-echo-shiritori/wiki/images/AWS/aws3.png" width= "400px"/>

	- `Alex開発コンソールにARNを貼り付けスキルIDをコピードリがーを設定`  
		<img src="https://github.com/monstar-lab/amazon-echo-shiritori/wiki/images/AWS/aws4.png" width= "400px"/>

	- `コードをZIP化してアップロード`  
		<img src="https://github.com/monstar-lab/amazon-echo-shiritori/wiki/images/AWS/aws5.png" width= "400px"/>


- Alexa 開発コンソールでスキル作成
	- `スキル名とスキル地域を選択後 -> 次へ`  
		<img src="https://github.com/monstar-lab/amazon-echo-shiritori/wiki/images/Alexa/Alexa1.png" width= "400px"/>

	- `カスタムスキルを作成`    
		<img src="https://github.com/monstar-lab/amazon-echo-shiritori/wiki/images/Alexa/Alexa2.png" width= "400px"/>

	- `スキルの呼び出す名を設定（既存のスキル名と設定すると別のスキルが起動してしまう可能性がある）`  
		<img src="https://github.com/monstar-lab/amazon-echo-shiritori/wiki/images/Alexa/Alexa3.png" width= "400px"/>

	- `カスタムインテントを作成`
		<img src="https://github.com/monstar-lab/amazon-echo-shiritori/wiki/images/Alexa/Alexa4.png" width= "400px"/>

	- `サンプル発話を追加`
	- `カスタムスロットを設定`
	- `インテントにスロットを追加`
	- `モデル保存して、ビルド忘れずに`
	-



  
