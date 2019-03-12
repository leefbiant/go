#!/bin/bash
dir=$(cd $(dirname $0); pwd)
cd $dir

process_list="spot|change|bigtrade|bigorder|premium|blastingorder|everydaymarket|deliverydate|user_push|bbx_notify"
time=$(date "+%Y-%m-%d %H:%M:%S")

function check_process() {
  array=($process_list)
  IFS="|"
  for var in ${array[@]}
  do
       full_process=${dir}/$var
       echo "$time  check process $full_process"
       process=`ps aux | grep $full_process | grep -v grep | awk '{print $2}'`
       if [ "X${process}" == "X" ]; then 
         echo "$time start process:" $full_process
         nohup $full_process &
       fi
  done
}

check_process &
