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
    echo "Using " $0 " all/huobi_trade/huobi_index/huobi_ticker/huobi_depth"
    return
  fi
  arg=$1
  if [ $1 == "all" ]; then
    start_process huobi_trade
    start_process huobi_index
    start_process huobi_ticker
    start_process huobi_depth
  else
    start_process $arg
  fi
}

start_arg_process $* &
