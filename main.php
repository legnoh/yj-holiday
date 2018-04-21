<?php
require_once 'vendor/autoload.php';
date_default_timezone_set('Asia/Tokyo');

// iCalの基本情報定義
$vCalendar = new \Eluceo\iCal\Component\Calendar('yahoojapan.holiday.tool.legnoh.lkj.io');
$vCalendar->setCalendarColor('#FF2968');
$vCalendar->setName("YJ holidays");
$events_keymap = array();

// CSVを入手して配列化(ヘッダ行は不要なので削除)
$events = array_map('str_getcsv', file('http://www8.cao.go.jp/chosei/shukujitsu/syukujitsu_kyujitsu.csv'));
unset($events[0]);

foreach ($events as $event) {
  addEvent(new DateTime($event[0]), mb_convert_encoding($event[1], 'utf-8', 'shift-jis'));

  // 土曜日の場合は前営業日が休日となる
  if((new DateTime($event[0]))->format('w') == '6'){
    addEvent(getBeforeWorkday(new DateTime($event[0])), '休日(YJ)');
  }

  // 年末年始は12/29~1/4が休日となる
  if((new DateTime($event[0]))->format('m-d') == '01-01'){
    addEvent((new DateTime($event[0]))->add(new DateInterval('P1D')), '年末年始休暇(YJ)'); # 01/02
    addEvent((new DateTime($event[0]))->add(new DateInterval('P2D')), '年末年始休暇(YJ)'); # 01/03
    addEvent((new DateTime($event[0]))->add(new DateInterval('P3D')), '年末年始休暇(YJ)'); # 01/04
    addEvent((new DateTime($event[0]))->add(new DateInterval('P362D')), '年末年始休暇(YJ)'); # 12/29
    addEvent((new DateTime($event[0]))->add(new DateInterval('P363D')), '年末年始休暇(YJ)'); # 12/30
    addEvent((new DateTime($event[0]))->add(new DateInterval('P364D')), '年末年始休暇(YJ)'); # 12/31
  }
}
file_put_contents('htdocs/yahoojapan/holidays.ics', $vCalendar->render());
file_put_contents('htdocs/yahoojapan/holidays.json', json_encode($events_keymap));


function addEvent($datetime, $name){
  global $vCalendar, $events_keymap;
  $vEvent = new \Eluceo\iCal\Component\Event();
  $vEvent->setDtStart($datetime)
         ->setDtEnd($datetime)
         ->setNoTime(true)
         ->setSummary($name);
  $vCalendar->addComponent($vEvent);
  $events_keymap[] = array(
    "name" => $name,
    "date" => $datetime->format('Y-m-d'),
    "start_time" => $datetime->format('U'),
    "end_time" => $datetime->add(new DateInterval('PT86399S'))->format('U')
  );
}

function getBeforeWorkday($datetime){
  global $events;
  while(true){
    $yesterday = $datetime->sub(new DateInterval('P1D'));
    if( $yesterday->format('w') > '0'
      && $yesterday->format('w') < '6'
      && array_search($yesterday->format('Y-m-d'), array_column($events, 0)) == FALSE ){
      return $yesterday;
    }
  }
}
