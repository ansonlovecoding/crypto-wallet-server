#!/usr/bin/env bash

source ./style_info.cfg
source ./path_info.cfg
source ./function.sh
service_port_name=(
  #api port name
  adminApiPort
  walletApiPort
  #rpc port name
  adminRPCPort
  walletRPCPort
  ethRPCPort
  btcRPCPort
  tronRPCPort
  pushRPCPort
)

TOKEN="bot5536406203:AAEE3VdEDl_LyXT8EIGM78mfKz0qsKxb080"
chat_ID="-647119299"
URL="https://api.telegram.org/${TOKEN}/sendMessage"

day=`date +"%Y-%m-%d"`
switch=$(cat $config_path | grep demoswitch |awk -F '[:]' '{print $NF}')
for i in ${service_port_name[*]}; do
  message_text="‼️【DEV】服务器有服务 ${i} 异常中断,请前往检查,如在更新代码重启服务请忽视"
  list=$(cat $config_path | grep -w ${i} | awk -F '[:]' '{print $NF}')
  list_to_string $list
  for j in ${ports_array}; do
    port=$(netstat -netulp tcp | awk '{print $4}' | grep -w ${j} | awk -F '[:]' '{print $NF}')
    if [[ ${port} -ne ${j} ]]; then
      echo -e ${YELLOW_PREFIX}${i}${COLOR_SUFFIX}${RED_PREFIX}" service does not start normally,not initiated port is "${COLOR_SUFFIX}${YELLOW_PREFIX}${j}${COLOR_SUFFIX}
      echo -e ${RED_PREFIX}"please check ../logs/wallet.log "${COLOR_SUFFIX}
      curl -s -X POST $URL -d chat_id=${chat_ID} -d text="${message_text}"      
      exit -1
    else
      echo -e ${j}${GREEN_PREFIX}" port has been listening,belongs service is "${i}${COLOR_SUFFIX}
    fi
  done
done

#Check launched cronjon service process
check=$(ps aux | grep -w ./${cron_service_name} | grep -v grep | wc -l)
message_text2="‼️【DEV】服务器有服务 ${cron_service_name} 异常中断,请前往检查,如在更新代码重启服务请忽视"
if [ $check -eq ${cron_service_service_num} ]; then
  echo -e ${GREEN_PREFIX}"none  port has been listening,belongs service is CronJob service"${COLOR_SUFFIX}
else
  echo -e ${RED_PREFIX}"CronJob service does not start normally, num err"${COLOR_SUFFIX}
        echo -e ${RED_PREFIX}"please check ../logs/wallet.log "${COLOR_SUFFIX}
        curl -s -X POST $URL -d chat_id=${chat_ID} -d text="${message_text2}"
      exit -1
fi

echo -e ${YELLOW_PREFIX}"all services launch success"${COLOR_SUFFIX}
