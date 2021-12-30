# yj-holiday

[![Netlify Status](https://api.netlify.com/api/v1/badges/9a81f2ea-9d0b-4bd4-a0d2-c182849a0936/deploy-status)](https://app.netlify.com/sites/yj-holidays/deploys)

- https://event.home.lkj.io/yahoojapan/holidays.ics
- https://event.home.lkj.io/yahoojapan/holidays.json

ヤフー株式会社における土日以外の休日をまとめたiCal/JSONファイルを生成するスクリプトです。

## usage

```sh
git clone https://github.com/legnoh/yj-holiday.git && cd yj-holiday
go mod vendor
go run main.go
```

## appendix

- ヤフー株式会社は、完全週休2日制（土日）、かつ国民の祝日、年末年始（12月29日から1月4日まで）が休日となります。
  - [採用情報 - ヤフー株式会社](https://about.yahoo.co.jp/hr/)
- ヤフー株式会社は、祝日が土曜日にあたった場合、前労働日を振り替え特別休日とする"土曜日祝日振替特別休日"があります。
  - [福利厚生 - 制度・環境 - 採用情報 - ヤフー株式会社](https://about.yahoo.co.jp/hr/workplace/welfare/)
- 日本国の祝日の取得に、内閣府の提供するCSVデータを利用しています。
  - [国民の祝日について - 内閣府](https://www8.cao.go.jp/chosei/shukujitsu/gaiyou.html) - [CSV](https://www8.cao.go.jp/chosei/shukujitsu/syukujitsu.csv)
