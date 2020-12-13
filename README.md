# kyopro progress reporter
みんなちゃんと競プロやってるかわかるようになります
![Screenshot from 2020-11-22 07-16-29](https://user-images.githubusercontent.com/43411965/99888824-dbd5c780-2c92-11eb-830d-183887db647c.jpg)

## Usage
　使い方も何もないけど、`/report` で `GET` メソッドで http リクエスト を受け取ったら slack のチャンネルにメンバー全員の AC カウントを投稿します。  
　heroku ではスケジューラーで毎日0時にリクエスト飛ばしています。ただ月曜日以外はどうやっても反応しないのでそれは注意。


# Test
## 1. .env にシークレットキーなど書き込む。
S3 環境がなければ slack api で必要なものだけ記述して `members.json` をローカルに置けば ok  
`PORT` はわからなければ `8080` で ok

## 2. 起動する
```console
$ go run main.go member.go report.go slacklib.go
```

## 3. http リクエストをとばす
```console
$ curl localhost:8080/report
```

## 4. 〜完〜
あとは slack の設定などごにょごにょしてください。


# Deploy
## 2. heroku の環境変数を変える
`.env` 内のものを全て set する
```console
$ heroku config:set 環境変数の名前=値
```

## 3. heroku にデプロイ
```console
$ git push heroku main
```

## 4. 〜完〜
動かなかったら `heroku logs --tail` で確認して対処してください。
