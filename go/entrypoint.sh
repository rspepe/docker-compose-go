#!/usr/bin/env bash

# 開発モード（air でホットリロード）か本番モード（go run）かを判定
if [ "$MODE" = "release" ]
then
  # 本番モード：通常の go run
  go run /src/main.go
else
  # 開発モード：air でホットリロード
  air -c /src/.air.toml
fi
