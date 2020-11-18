# kyopro progress reporter
みんなちゃんと競プロやってるかわかるようになります

# Test
書くのめんどくさいからとりあえず slack api 追加したところから書きます・・。
## 1. サーバーを起動する
```console
$ go run main.go
2020/11/19 00:16:17 [INFO] Server listening..
```

## 2. ngrok で Web サーバーを外部に公開する
```console
$ ngrok http 8080
```
で色々でてくる。

```console
Forwarding                    http://********.ngrok.io -> http://localhost:8080
Forwarding                    https://********.ngrok.io -> http://localhost:8080
```

下のほうの `https://********.ngrok.io` をコピーして Slack api の Event Subscriptions の Request URL に貼る。このときパターン(`/slack/events`)も追加することを忘れないこと。  
こんな感じになれば OK
![Screenshot from 2020-11-19 00-26-15](https://user-images.githubusercontent.com/43411965/99550605-3c000b80-29fe-11eb-828f-410f96efb995.png)


## 3. 