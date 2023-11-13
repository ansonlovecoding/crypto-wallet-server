#!/usr/bin/env bash
#Include shell font styles and some basic information
source ./style_info.cfg
source ./path_info.cfg

#Put config path to the ENV
export CONFIG_NAME=$config_path
day=`date +"%Y-%m-%d"`
echo "==========CONFIG_NAME:${CONFIG_NAME}==========="  >>../logs/wallet.log.${day} 2>&1 &

#Check if the service exists
#If it is exists,kill this process
check=`ps aux | grep -w ./${cron_service_name} | grep -v grep| wc -l`
if [ $check -ge 1 ]
then
oldPid=`ps aux | grep -w ./${cron_service_name} | grep -v grep|awk '{print $2}'`
 kill -9 $oldPid
fi
#Waiting port recycling
sleep 1

cd ${cron_service_binary_root}
for ((i = 0; i < ${cron_service_service_num}; i++)); do
      nohup ./${cron_service_name}    >>../logs/wallet.log.${day} 2>&1 &
done

#Check launched service process
check=`ps aux | grep -w ./${cron_service_name} | grep -v grep| wc -l`
if [ $check -ge 1 ]
then
newPid=`ps aux | grep -w ./${cron_service_name} | grep -v grep|awk '{print $2}'`
allPorts=""
    echo -e ${SKY_BLUE_PREFIX}"SERVICE START SUCCESS "${COLOR_SUFFIX}
    echo -e ${SKY_BLUE_PREFIX}"SERVICE_NAME: "${COLOR_SUFFIX}${YELLOW_PREFIX}${cron_service_name}${COLOR_SUFFIX}
    echo -e ${SKY_BLUE_PREFIX}"PID: "${COLOR_SUFFIX}${YELLOW_PREFIX}${newPid}${COLOR_SUFFIX}
    echo -e ${SKY_BLUE_PREFIX}"LISTENING_PORT: "${COLOR_SUFFIX}${YELLOW_PREFIX}${allPorts}${COLOR_SUFFIX}
else
    echo -e ${YELLOW_PREFIX}${cron_service_name}${RED_PREFIX}" start failure, please check wallet.log"${COLOR_SUFFIX}
fi
