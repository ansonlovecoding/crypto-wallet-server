#!/usr/bin/env bash

shell_array=(
  stop_all.sh
  start_all.sh
)

for i in ${shell_array[*]}; do
  chmod +x $i
  ./$i
  #waiting port recycling
  sleep 1
  if [ $? -ne 0 ]; then
        exit -1
  fi
done