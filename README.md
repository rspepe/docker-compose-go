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

---

## runcによるコンテナ実行（Docker/Podmanが使えない環境向け）

Docker/Podmanが利用できない環境（カーネル制限、外部レジストリアクセス制限など）でも、runcを使ってGoアプリケーションをコンテナとして実行できます。

### 前提条件

- runcがインストールされていること
- Goがインストールされていること（静的バイナリビルド用）
- root権限またはそれに相当する権限

### 実行手順

#### 1. OCIバンドルの準備

```bash
# OCIバンドル用ディレクトリ作成
mkdir -p /tmp/mycontainer/rootfs
cd /tmp/mycontainer

# config.jsonの作成（runc spec の代わりに手動で作成）
# 注: runcがインストールされている場合は `runc spec` で自動生成できます
```

#### 2. 静的バイナリのビルド

```bash
# プロジェクトディレクトリに移動
cd go/src

# 依存関係の整理
go mod tidy

# 静的リンクバイナリのビルド
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /tmp/mycontainer/rootfs/server main.go
```

#### 3. config.jsonの編集

以下の項目を設定：
- `process.terminal`: `false` に設定（非対話モード）
- `process.args`: `["/server"]` に設定（実行するコマンド）
- `root.readonly`: `false` に設定（書き込み可能）
- `linux.namespaces`: `network` タイプを削除（ホストネットワーク共有）

#### 4. コンテナの実行

```bash
# runcディレクトリの作成（必要な場合）
mkdir -p /run/runc

# コンテナをバックグラウンドで起動
cd /tmp/mycontainer
runc run -d myapp

# コンテナの状態確認
runc list

# 動作確認
curl http://127.0.0.1:8080/
curl http://127.0.0.1:8080/health
curl http://127.0.0.1:8080/hello/YourName
```

#### 5. コンテナの停止と削除

```bash
# コンテナの停止
runc kill myapp SIGTERM

# コンテナの削除
runc delete myapp
```

### 利用可能なエンドポイント

- `GET /` - Hello World メッセージを返す
- `GET /health` - ヘルスチェック（JSON形式）
- `GET /hello/{name}` - パラメータ付きの挨拶メッセージ

### トラブルシューティング

- runcがインストールされていない場合は、GitHubからソースをクローンしてビルドできます：
  ```bash
  cd /tmp
  git clone https://github.com/opencontainers/runc.git
  cd runc
  apt-get install -y libseccomp-dev  # 依存関係のインストール
  make runc
  cp runc /usr/local/bin/
  ```

- `/run/runc: no such file or directory` エラーが出る場合：
  ```bash
  mkdir -p /run/runc
  ```