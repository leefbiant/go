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
    echo "Using " $0 " all/bbx_trade"
    return
  fi
  arg=$1
  if [ $1 == "all" ]; then
    start_process bbx_trade
    start_process bbx_depth
    start_process bbx_ticker
  else
    start_process $arg
  fi
}

start_arg_process $*
