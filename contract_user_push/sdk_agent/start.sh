#!/bin/bash
dir=$(cd $(dirname $0); pwd)

function start_process(){
  name=$1
  full_process=${dir}/$name
  process=`ps aux | grep $full_process | grep -v grep | awk '{print $2}'`
  if [ "X${process}" != "X" ]; then
    echo "kill pid:" ${full_process}
    kill $process
  fi
  echo "start process:" $full_process
  nohup $full_process &
}

start_process sdk_agent
