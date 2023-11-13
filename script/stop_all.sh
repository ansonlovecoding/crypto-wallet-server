#!/usr/bin/env bash
#fixme This script is to stop the service

source ./style_info.cfg
source ./path_info.cfg

time=`date +"%Y-%m-%d %H:%M:%S"`
day=`date +"%Y-%m-%d"`
echo "==========================================================">>../logs/wallet.log.${day} 2>&1 &
echo "==========================================================">>../logs/wallet.log.${day} 2>&1 &
echo "==========================================================">>../logs/wallet.log.${day} 2>&1 &
echo "==========server stop time:${time}===========">>../logs/wallet.log.${day} 2>&1 &
echo "==========================================================">>../logs/wallet.log.${day} 2>&1 &
echo "==========================================================">>../logs/wallet.log.${day} 2>&1 &
echo "==========================================================">>../logs/wallet.log.${day} 2>&1 &

for i in ${service_names[*]}; do
  #Check whether the service exists
  name="ps aux |grep -w $i |grep -v grep"
  count="${name}| wc -l"
  if [ $(eval ${count}) -gt 0 ]; then
    pid="${name}| awk '{print \$2}'"
    echo -e "${SKY_BLUE_PREFIX}Killing service:$i pid:$(eval $pid)${COLOR_SUFFIX}">>../logs/wallet.log.${day} 2>&1 &
    #kill the service that existed
    kill -9 $(eval $pid)
    echo -e "${SKY_BLUE_PREFIX}service:$i was killed ${COLOR_SUFFIX}">>../logs/wallet.log.${day} 2>&1 &
  fi
done
