# yj-holiday

- http://holiday.tool.legnoh.lkj.io/yahoojapan.ics
- http://holiday.tool.legnoh.lkj.io/yahoojapan.json

ヤフー株式会社における土日以外の休日をまとめたiCal/JSONファイルを生成するスクリプトです。

## install
```sh
$ git clone https://github.com/legnoh/yj-holiday.git
$ composer install
$ cf push
```

## FYI

- 日本国の祝日の取得に、内閣府の提供するCSVデータを利用しています。
  - [国民の祝日について - 内閣府](http://www8.cao.go.jp/chosei/shukujitsu/gaiyou.html) - [CSV](http://www8.cao.go.jp/chosei/shukujitsu/syukujitsu_kyujitsu.csv)
