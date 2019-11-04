# golang-docker-compose

* 手軽にgolangのコードを実行するためのリポジトリ。

## 実行に必要なもの

* docker-composeが動作する環境

## 実行方法

* 以下のコマンドでライブリロード状態でmain.goが実行される。

```
docker-compose up
```

### ライブリロードを体感

* go/src/main.goのHello worldを変更して保存してみてください。保存直後に再ビルドされ実行結果がコンソールに表示されます。

## 環境削除

* 停止シグナル後も死んだコンテナが残っているので以下のコマンドで削除ができます。

```
docker-compose down
```