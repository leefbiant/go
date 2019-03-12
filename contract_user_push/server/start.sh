#!/bin/bash
dir=$(cd $(dirname $0); pwd)

function start_process(){
  name=$1
  full_process=${dir}/$name
  process=`ps aux | grep $full_process | grep -v grep | awk '{print $2}'`
  if [ "X${process}" != "X" ]; then
    echo "kill pid:" ${process}
    kill $process
  fi
  echo "start process:" $full_process
  nohup $full_process &
}


function start_arg_process() {
  if [ $# -eq 0 ]; then 
    echo "Using " $0 " all/spot/change/bigtrade/bigorder/premium/blastingorder/deliverydate/everydaymarket/user_push/bbx_notify"
    return
  fi
  arg=$1
  if [ $1 == "all" ]; then
    # 定点
    start_process spot
    # 涨跌
    start_process change
    # 大额挂单
    start_process bigtrade
    # 大单买卖
    start_process bigorder
    # 期现溢价
    start_process premium
    # 大额爆仓
    start_process blastingorder
    # 日报
    start_process deliverydate
    # 日报
    start_process everydaymarket
    # push
    start_process user_push
    # bbx notify
    start_process bbx_notify
  else
    start_process $arg
  fi
}

start_arg_process $* &
