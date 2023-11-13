#!/bin/bash
set -e

# Init
echo ""
echo "Init geth"
geth init "/root/files/genesis.json"
sleep 3

# Start geth
echo ""
echo "Start geth"
#geth --networkid=10086 --rpc --rpcapi "eth,net,web3,personal,admin,miner,debug,txpool" --rpcaddr "0.0.0.0" --rpcport "8545" --miner.threads 1 --mine --nat "extip:192.168.0.2" --allow-insecure-unlock &
#geth --networkid 10086 --rpc --rpcapi "eth,net,web3,personal,admin,miner,debug,txpool" --rpcaddr "0.0.0.0" --rpcport "8545" --miner.threads 1 --miner.etherbase "0x36587c80f8652875bcb4bb85de44409ef9a35245" --mine --nat "extip:192.168.0.2" --allow-insecure-unlock --http --http.corsdomain "https://remix.ethereum.org" --http.api "web3,eth,debug,personal,net" --snapshot=false --syncmode "fast" &
geth --networkid 10086 --http --http.api "eth,net,web3,personal,admin,miner,debug,txpool" --http.addr "0.0.0.0" --http.port "8545" --http.corsdomain "https://remix.ethereum.org" --ws --ws.addr "0.0.0.0" --ws.api "eth,net,web3" --ws.origins "*" --ws.port "8546" --miner.threads 1 --miner.etherbase "0x36587c80f8652875bcb4bb85de44409ef9a35245" --mine --nat "extip:192.168.0.2" --allow-insecure-unlock --gcmode archive --ignore-legacy-receipts --snapshot=false --syncmode "snap" &

sleep 10

while true; do
    sleep 1000000000
done
