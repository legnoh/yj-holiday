# yj-holiday

- http://event.home.lkj.io/yahoojapan/holidays.ics
- http://event.home.lkj.io/yahoojapan/holidays.json

ヤフー株式会社における土日以外の休日をまとめたiCal/JSONファイルを生成するスクリプトです。

## install
```sh
$ git clone https://github.com/legnoh/yj-holiday.git
$ composer install
$ cf push
```

## appendix

- ヤフー株式会社は、完全週休2日制（土日 ※1）、かつ国民の祝日、年末年始（12月29日から1月4日まで）が休日となります。
  - [募集要項 - 採用情報 - ヤフー株式会社](https://about.yahoo.co.jp/hr/guideline/)
- ヤフー株式会社は、祝日が土曜日にあたった場合、前労働日を振り替え特別休日とする"土曜日祝日振替特別休日"があります。
  - [休暇 - 制度・環境 - 採用情報 - ヤフー株式会社](https://about.yahoo.co.jp/hr/workplace/vacation.html)
- 日本国の祝日の取得に、内閣府の提供するCSVデータを利用しています。
  - [国民の祝日について - 内閣府](http://www8.cao.go.jp/chosei/shukujitsu/gaiyou.html) - [CSV](http://www8.cao.go.jp/chosei/shukujitsu/syukujitsu_kyujitsu.csv)
