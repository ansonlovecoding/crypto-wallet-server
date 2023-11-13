#!/usr/bin/env bash
#fixme This script is the total startup script
#fixme The full name of the shell script that needs to be started is placed in the need_to_start_server_shell array

#fixme Put the shell script name here

source ./path_info.cfg
#Put config path to the ENV
export CONFIG_NAME=$config_path
export ETH_TOML_PATH=$eth_toml_path
export BTC_TOML_PATH=$btc_toml_path
export TRON_TOML_PATH=$tron_toml_path

need_to_start_server_shell=(
  start_rpc_service.sh
  start_cron_service.sh
)

#fixme The 10 second delay to start the project is for the docker-compose one-click to start wallet when the infrastructure dependencies are not started

sleep 10
time=`date +"%Y-%m-%d %H:%M:%S"`
day=`date +"%Y-%m-%d"`
echo "=========================================================="  >>../logs/wallet.log.${day} 2>&1 &
echo "=========================================================="  >>../logs/wallet.log.${day} 2>&1 &
echo "=========================================================="  >>../logs/wallet.log.${day} 2>&1 &
echo "==========server start time:${time}==========="  >>../logs/wallet.log.${day} 2>&1 &
echo "=========================================================="  >>../logs/wallet.log.${day} 2>&1 &
echo "=========================================================="  >>../logs/wallet.log.${day} 2>&1 &
echo "=========================================================="  >>../logs/wallet.log.${day} 2>&1 &
for i in ${need_to_start_server_shell[*]}; do
  chmod +x $i
  ./$i
done

sleep 15

#fixme prevents the wallet service exit after execution in the docker container
tail -f /dev/null
