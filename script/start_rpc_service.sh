#!/usr/bin/env bash

source ./style_info.cfg
source ./path_info.cfg
source ./function.sh

#service filename
service_filename=(
  #api
  admin_api
  wallet_api
  #rpc
  admin_rpc
  wallet_rpc
  eth_rpc
  btc_rpc
  tron_rpc
  push_rpc
)

#service config port name
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
day=`date +"%Y-%m-%d"`
for ((i = 0; i < ${#service_filename[*]}; i++)); do
  #Check whether the service exists
  service_name="ps aux |grep -w ${service_filename[$i]} |grep -v grep"
  count="${service_name}| wc -l"

  if [ $(eval ${count}) -gt 0 ]; then
    pid="${service_name}| awk '{print \$2}'"
    echo  "${service_filename[$i]} service has been started,pid:$(eval $pid)"
    echo  "killing the service ${service_filename[$i]} pid:$(eval $pid)"
    #kill the service that existed
    kill -9 $(eval $pid)
    sleep 0.5
  fi
  cd ../bin
  #Get the rpc port in the configuration file
  portList=$(cat $config_path | grep ${service_port_name[$i]} | awk -F '[:]' '{print $NF}')
  list_to_string ${portList}
  #Start related rpc services based on the number of ports
  for j in ${ports_array}; do
    #Start the service in the background
    #    ./${service_filename[$i]} -port $j &
    nohup ./${service_filename[$i]} -port $j   >>../logs/wallet.log.${day} 2>&1 &
    sleep 1
    pid=$(eval "netstat -ntlp tcp|grep ${j} |awk '{printf \$7}'|cut -d/ -f1")
    if [ -n "${pid}" ]; then
      echo -e "${GREEN_PREFIX}${service_filename[$i]}${RED_PREFIX} start success,port number:${j} pid:${pid} ${COLOR_SUFFIX}"
    else
      echo -e "${YELLOW_PREFIX}${service_filename[$i]}${RED_PREFIX} start failure, please check wallet.log ${COLOR_SUFFIX}"
    fi

  done
done
