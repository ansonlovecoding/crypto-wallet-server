# Wallet Server

## About Project

### Introduction

This project is a HD wallet, it use microservice architecture by GRPC. The project include the SDK and the server code. SDK code is in the folder "server/cmd/wallet-sdk-core", it can be generated by gmobile, and be provided to app team. Server code contains docker-compose files for different coin node and seprate RPC to deal different node flows.



### What is HD Wallet?

A hierarchical deterministic (HD) wallet is a digital wallet commonly used to store the digital keys for holders of cryptocurrencies such as Bitcoin and Ethereum. Anyone with a copy of both the public and password-like private key can control the cryptocurrency in the account.

![68747470733a2f2f692e696d6775722e636f6d2f666274324132712e706e67](README.assets/68747470733a2f2f692e696d6775722e636f6d2f666274324132712e706e67.png)



![v2-fdf6b9f54a1f51bfa4eb5e3b2ba88d22_r](README.assets/v2-fdf6b9f54a1f51bfa4eb5e3b2ba88d22_r.jpg)

### Code Structure

```
server
├─script
|   ├─build_all_service.sh
|   ├─check_all.sh
|   ├─docker_check_service.sh
|   ├─docker_start_all.sh
|   ├─env_check.sh
|   ├─function.sh
|   ├─mongo-init.sh		#after running the mongo container, need to execute this script to create a mongo user
|   ├─path_info.cfg
|   ├─restart_all.sh
|   ├─start_all.sh
|   ├─start_rpc_service.sh
|   ├─stop_all.sh
|   ├─style_info.cfg
|   └swagger.sh 	#generate the swagger docs
├─pkg
|  ├─coingrpc		#grpc classes for all coins
|  |    ├─ethgrpc		#eth
|  |    |     ├─api-interface.go	#api interface
|  |    |     └connection.go			#connect eth rpc
|  ├─utils
|  ├─swagger	#swagger files
|  |    ├─wallet_api
|  |    |     ├─docs.go
|  |    |     ├─swagger.json
|  |    |     └swagger.yaml
|  |    ├─admin_api
|  |    |     ├─docs.go
|  |    |     ├─swagger.json
|  |    |     └swagger.yaml
|  ├─struct
|  ├─proto
|  |   ├─auto_proto.sh	#execute the script if has any update in proto
|  |   ├─proto_dir.cfg	#file list that need generate
|  ├─grpc-etcdv3
|  ├───db
|  |   ├─model.go		#init mysql, redis and mongo
|  |   ├─redis
|  |   ├─mysql
|  |   ├─mongo
|  |   ├─local_db		#local database for sdk
|  ├─common
|  |   ├─token_verify
|  |   ├─log
|  |   ├─http
|  |   ├─constant
|  |   ├─config
|  |   |   └config.go		#init config.yaml
├─logs	#all the logs will be stored in this folder
├─internal
|    ├─wallet_api
|    ├─wallet-sdk-core		#sdk internal classes
|    ├─rpc
|    |  ├─wallet
|    |  ├─eth
|    |  ├─bitcoin
|    |  ├─admin
|    ├─admin_api
├─config
|   ├─config.example.yaml		#copy it and rename to "config.yaml", update the env settings in the file
|   ├─nginx
|   |   ├─conf.d
|   |   |   └wallet.conf		#check the domain here and add the domain to your hosts(/etc/hosts)
|   ├─wallet
|   |   └eth.toml		#eth node connection configuration
├─cmd		#main file for each service
|  ├─wallet_api
|  |     ├─Makefile
|  |     └main.go
|  ├─wallet-sdk-core
|  |        ├─Makefile
|  |        ├─hot-wallet
|  |        |     └hot-interface.go
|  |        ├─cold-wallet
|  |        |      └cold_interface.go
|  ├─rpc
|  |  ├─wallet_rpc
|  |  ├─eth_rpc
|  |  ├─bitcoin_rpc
|  |  ├─admin_rpc
|  ├─admin_api
|  |     ├─Makefile
|  |     └main.go
├─bin		#binary files will be here
├─docker	#main file for each service
|  ├─images
|  |     ├─geth-1.10.8
|  |     |   ├─bin
|  |     |   |  	└start.sh		#run geth node
|  |     |   ├─data	#store the data from geth node
|  |     |   ├─files
|  |     |   |  	└genesis.json		#genesis json for eth node
|  |     |   └Dockerfile
|  |     ├─bitcoin
|  |     |   ├─bitcoin-testnet-box
|  |     |   |  	├───1 #node 1	
|  |     |   |  	| 	└bitcoin.conf #configuration file for node1
|  |     |   |  	├───2 #node 2	
|  |     |   |  	| 	└bitcoin.conf #configuration file for node2
|  |     └share-wallet
|  |        └deploy.Dockerfile	#dockerfile for buidling server images
├─docker-compose.btc.private.yml  	#docker-compose for btc node in private net
├─docker-compose.eth.private.yml  	#docker-compose for eth node in private net
└docker-compose.yaml		#docker-compose for server

```



### How to set up environment by using the docker?

- [ ] Change Server Configuration (config/config.yaml)

  - Copy config.example.yaml and paste it in the config folder, change the file name to "config.yaml"
  - Change the etcd, mysql, mongo, redis, jpns setting according to your server environment
  - Change "swaggerEable" to false in PROD server, if false will unable to access the swagger url
  - Change "is_skip_database" to true if you want to run ethgrpc unit test, otherwise set it false

- [ ] Change Ethereum Node Connection Configuration (config/wallet/eth.toml)

  - Change the "host" to your node IP and change the "port" to your node port, default port is 8545
  - Change keydir to your keystore absolute path if you want to run the unit test, otherwise, keep it empty

- [ ] Change Bitcoin Node Connection Configuration (config/wallet/btc.toml)

  - Change the "host" to your node IP with port
  
- [ ] Install Docker

- [ ] Pull images first to make it build more quickly

  ```
  docker pull ubuntu
  doker pull golang
  docker pull ethereum/client-go:v1.10.8
  ```

- [ ] Change docker-compose file

  - Update your $GOPATH to the volumes of server_share in docker-compose.yaml

    `- /Users/anson/go:/go`

  - Run "docker-compose up" to start the server

  - Go inside the mongo container, run the "mongo-init.sh" to create the new user for project



### Use Different Chains in different environments

- [ ] ETH
  - [ ] Local/Dev Environment: Private Chain
  - [ ] Test Environment: Goerli Testnet Chain
  - [ ] Prod Environment: Mainnet Chain
- [ ] Bitcoin
  - [ ] Local/Dev Environment: regtest chain (private)
  - [ ] Test Environment: Chain
  - [ ] Prod Environment: Mainnet Chain

- [ ] 



### How to generate swagger docs? (v1.6.7)

After updating the APIs, need to generate the swagger docs by running the script "swagger.sh" in script folder



## About ETH

### How to run the ETH private chain?

- [ ] Update your $GOPATH to the volumes of eth_private in docker-compose.eth.private.yaml

  `- /Users/anson/go:/go`

- [ ] Run "docker-compose -f docker-compose.eth.private.yaml  up" to start the Eth node server

- [ ] Create some eth account to test

  ```
  #get into ipc
  geth attach /root/.ethereum/geth.ipc
  #create account
  personal.newAccount()
  #type password (very important)
  #the private key file will be store in /docker/geth-1.10.8/data/chain/keystore, do not delete the files
  ```

- [ ] Put the account's hex address in "alloc" section in genesis.json, and the first address in "alloc" will be the codebase of the node, and also you can set the balance of the addresses , then run again

  

### Try some commands in Geth

```
#get into ipc
geth attach /root/.ethereum/geth.ipc

#create account
personal.newAccount()
#type password (very important)
#the private key file will be store in /docker/geth-1.10.8/data/chain/keystore, do not delete the files

#check account list
eth.accounts

#Check balance by ether
web3.fromWei(eth.getBalance("0x048caa04b0976aa80f8a18616d0f6c13b27d4e5b"), "ether")

#transfer
personal.unlockAccount("0x048caa04b0976aa80f8a18616d0f6c13b27d4e5b")

#send 10eth = 10000000000000000000wei = 0x8ac7230489e80000(hexadecimal)
eth.sendTransaction({from:"0x048caa04b0976aa80f8a18616d0f6c13b27d4e5b", to:"0x0383ffCDBA4c07FA62Bac71eF16dF35C64E0Ba42", value:"0x8ac7230489e80000", gas:"0x5208", gasPrice:"0xa"})

eth.getTransaction("0x687067ee0617950dd8c979a7051ba3d569d7ee387a6535cfbaa9db68191656eb")

> eth.getBlockByNumber("19861")
{
  difficulty: "0x96c27",
  extraData: "0xd683010a1a846765746886676f312e3137856c696e7578",
  gasLimit: "0x1c9c380",
  gasUsed: "0x106f5",
  hash: "0x6a12abd38c3b9a92cecbb46aaaca04618a1e243aee729a8935128832797be46c",
  logsBloom: "0x00000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000008000400000000000000001000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000020000000000000000000000000000000000000000000000100000000200000000000000000080000000000000000",
  miner: "0x77f3b1b21c2d4ee6f4f1e877866166c9179373d5",
  mixHash: "0x58ed6868e208d9b2fff1ac54a0463d8b2a6d9db23589a1a6684e61509a0ddcd8",
  nonce: "0x6119a1752051e22e",
  number: "0x4d95",
  parentHash: "0xba4083e3cc77b52a62537426ac350d11d0dd79a100096ff585ebdc26f50883b6",
  receiptsRoot: "0x1172f2546199dbd2de1f375dc584eefb133a53d8942a38c195c0b83b5a40ab50",
  sha3Uncles: "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
  size: "0x339",
  stateRoot: "0x387f641a5ed0f44c33cd8b1b7804c2778db3e5c98bebe89c3b094683e728bfce",
  timestamp: "0x63d78892",
  totalDifficulty: "0x2850cdc13",
  transactions: ["0x7519038667f9991e21a9bc178a4557529421a6c0706e9885fa33eee29cf27649", "0xc3e440db6087bf683d873d50e9aec7d4d6e52e3f3f3f4e5e37bc147763436966"],
  transactionsRoot: "0x3dd9f3f473e6dc20794199d731058bc49ef18af9e69ccd8ce733fff439a222e7",
  uncles: []
}

> eth.getTransactionReceipt("0x7519038667f9991e21a9bc178a4557529421a6c0706e9885fa33eee29cf27649")
{
  blockHash: "0x6a12abd38c3b9a92cecbb46aaaca04618a1e243aee729a8935128832797be46c",
  blockNumber: 19861,
  contractAddress: null,
  cumulativeGasUsed: 21000,
  effectiveGasPrice: 1000000000,
  from: "0xecb8441d0081d5ef09e134cb527d7e748b4b5777",
  gasUsed: 21000,
  logs: [],
  logsBloom: "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
  status: "0x1",
  to: "0xda074b96a9f8cb4c4d53ece51cda457f3c00322e",
  transactionHash: "0x7519038667f9991e21a9bc178a4557529421a6c0706e9885fa33eee29cf27649",
  transactionIndex: 0,
  type: "0x0"
}

#get latest block number, so the comfirm block number for a tx is (latest_block_number - tx_block_number)
> eth.getBlockByNumber("latest")
{
  difficulty: "0x966c7",
  extraData: "0xd683010a1a846765746886676f312e3137856c696e7578",
  gasLimit: "0x1c9c380",
  gasUsed: "0x0",
  hash: "0x737a84de98a0ecb8b7f254995ac6a262bc069ca111c2f67821abdc13fed56ecb",
  logsBloom: "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
  miner: "0x77f3b1b21c2d4ee6f4f1e877866166c9179373d5",
  mixHash: "0x2982db8da699d5e80d2481d0a276c569f9bcbe5f331cbe2fb1c621724e088fd8",
  nonce: "0x4a2b8aa695ae9da7",
  number: "0x523f",
  parentHash: "0xa96464555b43212d7589f32d916c8d034e69ff092d24546022f66740aff84ce8",
  receiptsRoot: "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
  sha3Uncles: "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
  size: "0x219",
  stateRoot: "0x7294c939386c738c92caeee32d72f7a0bc54e9d3c99c4ad05a283d6ce0a079a4",
  timestamp: "0x63d7c38e",
  totalDifficulty: "0x2b1360366",
  transactions: [],
  transactionsRoot: "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
  uncles: []
}
```



### How to deploy USDT ERC20 to private chain by using Remix?

- [ ] Go etherscan and find the Tether contract, copy the source code, url: https://etherscan.io/address/0xdac17f958d2ee523a2206206994597c13d831ec7

  ![image-20221230145507717](README.assets/image-20221230145507717.png)

- [ ] Go remix ([https://remix.ethereum.org/](https://remix.ethereum.org/)) and create a tether.sol file in contracts floder, then paste the source code to the editor

  ![image-20221230150106892](README.assets/image-20221230150106892.png)

- [ ] Go to the complier page, and change the compiler version same with the contract "0.4.18+commit.9cf6e910", then complie tether.sol

  ![image-20221230150423629](README.assets/image-20221230150423629.png)

- [ ] Go to deploy page, you need to make sure your private chain is running and exposure the 8545 to http port, and alow the corsdomain when you start the geth (check the start.sh in docker/images/geth-1.10.8). Then you can choose the "External Http Provider" and choose a account has eth balance, choose the contract "Tether Token - contracts/tether.sol". Set the deploy parameters same with the picture below, click the "transact" button

  ![image-20221230151405053](README.assets/image-20221230151405053.png)

- [ ] It will show the error "creation of TetherToken errored: authentication needed: password or unlock" in the console, that means you have to unlock the ethereum account you choosed in the node，then click the "transact" button again

  `personal.unlockAccount("0x048caa04b0976aa80f8a18616d0f6c13b27d4e5b")`

- [ ] After success, you will get a contract address in the "Deployed Contracts"

  ![image-20221230152513513](README.assets/image-20221230152513513.png)

- [ ] The contract address is fixed and unique in blockchain, now you can use the contract address "0x7FB7a9356E8Ce2579de30638dc09cAA7Af7066A8" to transact. But in mainnet the USDT-ERC20 address is "0xdAC17F958D2ee523a2206206994597C13D831ec7"



### How to deploy USDT ERC20 to private chain by using script?

- [ ] Install node

  ```
  wget https://nodejs.org/download/release/v10.12.0/node-v10.12.0-linux-x64.tar.xz
  tar xvf node-v10.12.0-linux-x64.tar.xz 
  mv node-v10.12.0-linux-x64 /usr/local/nodejs
  
  vim /etc/profile
  #add nodejs path to profile
  export NODE_HOME=/usr/local/nodejs
  export PATH=$NODE_HOME/bin:$PATH
  
  #make the profile take effect
  source /etc/profile
  
  #check the version
  node -v
  npm -v
  ```

- [ ] Create the folder "contract" and copy contract code to the folder /data/contract

  ```
  mkdir /data/contract
  cd /data/contract
  vim usdt.sol
  
  #copy the contract code from https://etherscan.io/address/0xdac17f958d2ee523a2206206994597c13d831ec7#code then save it
  ```

- [ ] Create a deploy script to the folder /data/contract, and change the RPC address

  ```
  cd /data/contract
  vim deploy_usdt.js
  
  #paste the code below, and change the rpc address and sender address
  let Web3 = require("web3")
  let solc = require("solc")
  let fs = require("fs")
  
  // change the RPC address
  web3 = new Web3(new Web3.providers.HttpProvider("http://192.168.0.125:8545"));
  
  // read the contract code
  let source = fs.readFileSync("./usdt.sol", "utf8");
  
  // compile the contract
  let cacl = solc.compile(source, 1);
  
  console.log("=============================")
  console.log("version: " + web3.version);
  console.log("=============================")
  console.log(solc.version())
  console.log("=============================")
  console.log(cacl.contracts)
  
  let abi = JSON.parse(cacl.contracts[':TetherToken'].interface);
  console.log("=============================")
  console.log(abi)
  console.log("=============================")
  
  let bytecode = cacl.contracts[':TetherToken'].bytecode;
  console.log(bytecode)
  console.log("=============================")
  
  let sender = "0xda074b96a9f8cb4c4d53ece51cda457f3c00322e";
  
  // deploy the contract
  new web3.eth.Contract(abi).deploy({
  	data: "0x" + bytecode,
  	arguments: [2718282552619854,"Tether USD","USDT",6],
  }).send({
  	from: sender,
  	gas: 3070000,
  	gasPrice: 130000000000
  }, function (error) {
  	if (error != null) {
  		console.log("failed in deployment: " + error);
  	}
  	}).then(function (newContractInstance) {
  	// request the contract address
  	var newContractAddress = newContractInstance.options.address
  	console.log("deploy success, the contract address is：\n" + newContractAddress.toLowerCase());
  });
  ```

- [ ] Create the Install encoding@^0.1.0, web3, solc@0.4.18

  ```
  cd /data/contract
  npm install encoding@^0.1.0
  npm install web3
  npm install solc@0.4.18
  ```

- [ ] Run command `node deploy_usdt.js` to deploy the contract

- [ ] After success it will return the contract address




### Use Metamask to test your private chain

- [ ] Install Metamask to your Chrome: https://chrome.google.com/webstore/detail/metamask/nkbihfbeogaeaoehlefnkodbefgpgknn/related
- [ ] Add private network to your Metamask

![image-20230126114452534](README.assets/image-20230126114452534.png)

- [ ] Import your account to your metamask

  ![image-20230126114744468](README.assets/image-20230126114744468.png)

- [ ] After few minutes, your account will show in Metamask

  ![image-20230126114815440](README.assets/image-20230126114815440.png)

- [ ] Import the USDT-ERC20 token, the token address is different in different environment

  ![image-20230130134908236](README.assets/image-20230130134908236.png)

## About Bitcoin

### How to run the Bitcoin private chain?

- [ ] Run "docker-compose -f docker-compose.btc.private.yaml  up" to start the Bitcoin node server
- [ ] If need to clear the blockchain data, delete the "btc_private" folder in components



### Bitcoin Transaction flows

![msbt_0204](README.assets/msbt_0204.png)



### Try some commands in Bitcoin-core (Testing flows)

```
1.Returns details of the blockchain (check the chain name)
bitcoin-cli -datadir=1 getblockchaininfo
or
curl --data-binary '{"jsonrpc":"1.0","id":"1","method":"getblockchaininfo","params":[]}' http://admin1:123@127.0.0.1:19001/
return:
{
  "chain": "regtest",
  "blocks": 0,
  "headers": 0,
  "bestblockhash": "0f9188f13cb7b2c71f2a335e3a4fc328bf5beb436012afca590b1a11466e2206",
  "difficulty": 4.656542373906925e-10,
  "time": 1296688602,
  "mediantime": 1296688602,
  "verificationprogress": 1,
  "initialblockdownload": true,
  "chainwork": "0000000000000000000000000000000000000000000000000000000000000002",
  "size_on_disk": 293,
  "pruned": false,
  "warnings": ""
}

2.Creates and loads a new wallet.(descriptors should be false)
bitcoin-cli -datadir=1 createwallet anson false false 123456 true false
or
curl --data-binary '{"jsonrpc":"1.0","id":"1","method":"createwallet","params":["anson",false,false,"123456",true,false]}' http://admin1:123@127.0.0.1:19001/
return:
{
  "name": "anson",
  "warning": "Wallet created successfully. The legacy wallet type is being deprecated and support for creating and opening legacy wallets will be removed in the future."
}

3.generate blocks for the blockchain (it will generate 50btc from 101 blocks, and it will create an address for the wallet)
bitcoin-cli -datadir=1 -generate 101
return:
{
  "address": "bcrt1q20hcyyamkraj423nluzg95lqqj50hgyrq4pykf",
  "blocks": [
    "36c2a3b023c55572159779452a2475d5d6c6e25dc0c46e1e7ad981925214e830",
    "0c6d2d8fdfa34de9c40a03fd0579a9bd7d1897a1cae76b19fdb216a70ce4bd20",
    "79a1884b63fd55b8e8e1175bbdde679cc7844e3ad6a229a3d696689109529872",
    ...
    "28155f053ba5acc15b98d6d1bef91313e9bbb7a73a1393eb991eb8073ea3f372",
    "124dfe561055da14c5569dedd0c3b70b1f872ba8c994a208ad84deaa0e3a23cf"
  ]
}

4.check the total available balance.
bitcoin-cli -datadir=1 -rpcwallet=anson getbalance "*" 6
or
curl --data-binary '{"jsonrpc":"1.0","id":"1","method":"getbalance","params":["*", 6]}' http://admin1:123@127.0.0.1:19001/wallet/anson
return: 50.00000000

5.check an object with all balances in BTC.
bitcoin-cli -datadir=1 -rpcwallet=anson getbalances
or
curl --data-binary '{"jsonrpc":"1.0","id":"1","method":"getbalances","params":[]}' http://admin1:123@127.0.0.1:19001/wallet/anson
return:
{
  "mine": {
    "trusted": 50.00000000,
    "untrusted_pending": 0.00000000,
    "immature": 5000.00000000,
    "used": 0.00000000
  }
}

5.create another wallet to receive all bitcoins from first wallet
bitcoin-cli -datadir=1 createwallet anson1 false false 123456 true false
return:
{
  "name": "anson1",
  "warning": "Wallet created successfully. The legacy wallet type is being deprecated and support for creating and opening legacy wallets will be removed in the future."
}

6.create a new Bitcoin address for receiving payments.
bitcoin-cli -datadir=1 -rpcwallet=anson1 getnewaddress
or
curl --data-binary '{"jsonrpc":"1.0","id":"1","method":"getnewaddress","params":[]}' http://admin1:123@127.0.0.1:19001/wallet/anson1
return
{"result":"bcrt1qef9nev7sg0dwjfplx89ecqnv3950w2avkgm86z","error":null,"id":"1"}

7.for easy to test, send all coins to second wallet
sendtoaddress (fast send money, but will not use in project)
bitcoin-cli -datadir=1 -rpcwallet=anson sendtoaddress bcrt1qef9nev7sg0dwjfplx89ecqnv3950w2avkgm86z 40.0 "" "" false true 1 "unset" true
or
curl --data-binary '{"jsonrpc":"1.0","id":"1","method":"sendtoaddress","params":["bcrt1qef9nev7sg0dwjfplx89ecqnv3950w2avkgm86z",40.00,"","",false,true,1,"unset",true]}' http://worker:123456@127.0.0.1:18443/wallet/anson2
return: 760f4fd52553b84cdd709820e6da4a6f3421b46a9d4195bf55b288cd5db1203a

8.walletpassphrase (unlock the wallet)
bitcoin-cli -datadir=1 -rpcwallet=anson walletpassphrase 123456 60
or
curl --data-binary '{"jsonrpc":"1.0","id":"1","method":"walletpassphrase","params":["123456", 60]}' http://admin1:123@127.0.0.1:19001/wallet/anson

9.get array of unspent transaction outputs
bitcoin-cli -datadir=1 -rpcwallet=anson listunspent 0
or
curl --data-binary '{"jsonrpc":"1.0","id":"1","method":"listunspent","params":[0]}' http://admin1:123@127.0.0.1:19001/wallet/anson
return
[
  {
    "txid": "760f4fd52553b84cdd709820e6da4a6f3421b46a9d4195bf55b288cd5db1203a",
    "vout": 0,
    "address": "bcrt1q834u98eal72rgg6d74h0yfaentyh9q8rtt7qkx",
    "scriptPubKey": "00143c6bc29f3dff9434234df56ef227b99ac97280e3",
    "amount": 9.99998590,
    "confirmations": 0,
    "ancestorcount": 1,
    "ancestorsize": 141,
    "ancestorfees": 1410,
    "spendable": true,
    "solvable": true,
    "desc": "wpkh([05193e72/0'/1'/0']02c304c7bbb6ebf4a5c9e5dc076d85e860494a4aa924b8db6cda8e30c94778ec72)#nxms5dpx",
    "parent_descs": [
    ],
    "reused": false,
    "safe": true
  }
]

10.at this time the comfirmations is 0, so it will not complete, need to generate a block to comfirm, that's mining
bitcoin-cli -datadir=1 generateblock "bcrt1q20hcyyamkraj423nluzg95lqqj50hgyrq4pykf" '["760f4fd52553b84cdd709820e6da4a6f3421b46a9d4195bf55b288cd5db1203a"]'
return:
{
  "hash": "38a0d0b8c7ce36512195ac1737329d6b5f1e6c69376f85ea726b68dbf953df97"
}

11.then check the unspent for first wallet and second wallet
bitcoin-cli -datadir=1 -rpcwallet=anson listunspent 0
return:
[
  {
    "txid": "a91d6d9b88b448b3bfe42f55295f17495aa22af49921c6d3248f031c669da4d0",
    "vout": 0,
    "address": "bcrt1q20hcyyamkraj423nluzg95lqqj50hgyrq4pykf",
    "label": "",
    "scriptPubKey": "001453ef8213bbb0fb2aaa33ff0482d3e004a8fba083",
    "amount": 50.00000000,
    "confirmations": 101,
    "spendable": true,
    "solvable": true,
    "desc": "wpkh([05193e72/0'/0'/0']039c2490ae3b93f74243bde3a4c83efba44c2f1ac85c7519956fe6c07373754416)#pzkxq6zv",
    "parent_descs": [
    ],
    "reused": true,
    "safe": true
  },
  {
    "txid": "760f4fd52553b84cdd709820e6da4a6f3421b46a9d4195bf55b288cd5db1203a",
    "vout": 0,
    "address": "bcrt1q834u98eal72rgg6d74h0yfaentyh9q8rtt7qkx",
    "scriptPubKey": "00143c6bc29f3dff9434234df56ef227b99ac97280e3",
    "amount": 9.99998590,
    "confirmations": 1,
    "spendable": true,
    "solvable": true,
    "desc": "wpkh([05193e72/0'/1'/0']02c304c7bbb6ebf4a5c9e5dc076d85e860494a4aa924b8db6cda8e30c94778ec72)#nxms5dpx",
    "parent_descs": [
    ],
    "reused": false,
    "safe": true
  }
]

bitcoin-cli -datadir=1 -rpcwallet=anson1 listunspent 0
[
  {
    "txid": "760f4fd52553b84cdd709820e6da4a6f3421b46a9d4195bf55b288cd5db1203a",
    "vout": 1,
    "address": "bcrt1qef9nev7sg0dwjfplx89ecqnv3950w2avkgm86z",
    "label": "",
    "scriptPubKey": "0014ca4b3cb3d043dae9243f31cb9c026c8968f72bac",
    "amount": 40.00000000,
    "confirmations": 1,
    "spendable": true,
    "solvable": true,
    "desc": "wpkh([0881a97c/0'/0'/0']039f62c054d7cb7652c4c76723861ba458fdd43b738e329a384774f4bb23a4d656)#9gd4rjnw",
    "parent_descs": [
    ],
    "reused": false,
    "safe": true
  }
]

12.check the second wallet balance
bitcoin-cli -datadir=1 -rpcwallet=anson1 getbalances
return:
{
  "mine": {
    "trusted": 40.00000000,
    "untrusted_pending": 0.00000000,
    "immature": 0.00000000,
    "used": 0.00000000
  }
}

13.create third wallet to test raw transaction
bitcoin-cli -datadir=1 createwallet anson2 false false 123456 true false

14.create a new Bitcoin address for third wallet.
bitcoin-cli -datadir=1 -rpcwallet=anson2 getnewaddress
return: bcrt1qzrp4uhkehjperc6muwzgvxxvkjjyp06zgtkmj4

15.now try to send 10btc to third wallet

Create a transaction with the inputs from UTXO for sending wallet
bitcoin-cli -datadir=1 createrawtransaction "[{\"txid\":\"760f4fd52553b84cdd709820e6da4a6f3421b46a9d4195bf55b288cd5db1203a\",\"vout\":1}]" "[{\"bcrt1qzrp4uhkehjperc6muwzgvxxvkjjyp06zgtkmj4\":9.9997},{\"bcrt1qef9nev7sg0dwjfplx89ecqnv3950w2avkgm86z\":30.00}]"
or
curl --data-binary '{"jsonrpc":"1.0","id":"1","method":"createrawtransaction","params":[[{"txid":"760f4fd52553b84cdd709820e6da4a6f3421b46a9d4195bf55b288cd5db1203a","vout":1}], [{"bcrt1qzrp4uhkehjperc6muwzgvxxvkjjyp06zgtkmj4":9.9997},{"bcrt1qef9nev7sg0dwjfplx89ecqnv3950w2avkgm86z":30.00}]]}' http://admin1:123@127.0.0.1:19001/
return:
02000000013a20b15dcd88b255bf95419d6ab421346f4adae6209870dd4cb85325d54f0f760100000000fdffffff02d0549a3b0000000016001410c35e5ed9bc8391e35be3848618ccb4a440bf42005ed0b200000000160014ca4b3cb3d043dae9243f31cb9c026c8968f72bac00000000

16.unlock second wallet before doing sign
bitcoin-cli -datadir=1 -rpcwallet=anson1 walletpassphrase 123456 60

17.signrawtransactionwithwallet
Sign inputs for raw transaction (serialized, hex-encoded).
bitcoin-cli -datadir=1 -rpcwallet=anson1 signrawtransactionwithwallet 02000000013a20b15dcd88b255bf95419d6ab421346f4adae6209870dd4cb85325d54f0f760100000000fdffffff02d0549a3b0000000016001410c35e5ed9bc8391e35be3848618ccb4a440bf42005ed0b200000000160014ca4b3cb3d043dae9243f31cb9c026c8968f72bac00000000
or
curl --data-binary '{"jsonrpc":"1.0","id":"1","method":"signrawtransactionwithwallet","params":["02000000013a20b15dcd88b255bf95419d6ab421346f4adae6209870dd4cb85325d54f0f760100000000fdffffff02d0549a3b0000000016001410c35e5ed9bc8391e35be3848618ccb4a440bf42005ed0b200000000160014ca4b3cb3d043dae9243f31cb9c026c8968f72bac00000000"]}' http://admin1:123@127.0.0.1:19001/wallet/anson1
return:
{
  "hex": "020000000001013a20b15dcd88b255bf95419d6ab421346f4adae6209870dd4cb85325d54f0f760100000000fdffffff02d0549a3b0000000016001410c35e5ed9bc8391e35be3848618ccb4a440bf42005ed0b200000000160014ca4b3cb3d043dae9243f31cb9c026c8968f72bac024730440220764b58d9136c4142060a3c3c3c6ceffde25a11338fd36d6114635818e02df6eb022030af798a6fee98cd1d35ab0254f9c3e28ca045f791d234c12f81aa003321ea4a0121039f62c054d7cb7652c4c76723861ba458fdd43b738e329a384774f4bb23a4d65600000000",
  "complete": true
}

18.Submit a raw transaction (serialized, hex-encoded) to local node and network.
bitcoin-cli -datadir=1 sendrawtransaction 020000000001013a20b15dcd88b255bf95419d6ab421346f4adae6209870dd4cb85325d54f0f760100000000fdffffff02d0549a3b0000000016001410c35e5ed9bc8391e35be3848618ccb4a440bf42005ed0b200000000160014ca4b3cb3d043dae9243f31cb9c026c8968f72bac024730440220764b58d9136c4142060a3c3c3c6ceffde25a11338fd36d6114635818e02df6eb022030af798a6fee98cd1d35ab0254f9c3e28ca045f791d234c12f81aa003321ea4a0121039f62c054d7cb7652c4c76723861ba458fdd43b738e329a384774f4bb23a4d65600000000
or
curl --data-binary '{"jsonrpc":"1.0","id":"1","method":"sendrawtransaction","params":["020000000001013a20b15dcd88b255bf95419d6ab421346f4adae6209870dd4cb85325d54f0f760100000000fdffffff02d0549a3b0000000016001410c35e5ed9bc8391e35be3848618ccb4a440bf42005ed0b200000000160014ca4b3cb3d043dae9243f31cb9c026c8968f72bac024730440220764b58d9136c4142060a3c3c3c6ceffde25a11338fd36d6114635818e02df6eb022030af798a6fee98cd1d35ab0254f9c3e28ca045f791d234c12f81aa003321ea4a0121039f62c054d7cb7652c4c76723861ba458fdd43b738e329a384774f4bb23a4d65600000000"]}' http://admin1:123@127.0.0.1:19001/
return:
{"result":"c6cd113f9b96bc501fd3d903bdec87820db78050cfeb72fddf31f1924b49a328","error":null,"id":"1"}

19.check the balances for two wallets again
bitcoin-cli -datadir=1 -rpcwallet=anson1 getbalances
{
  "mine": {
    "trusted": 0.00000000,
    "untrusted_pending": 0.00000000,
    "immature": 0.00000000,
    "used": 30.00000000
  }
}

bitcoin-cli -datadir=1 -rpcwallet=anson2 getbalances
{
  "mine": {
    "trusted": 0.00000000,
    "untrusted_pending": 9.99970000,
    "immature": 0.00000000,
    "used": 0.00000000
  }
}

trusted balance for third wallet still 0, we need to generate a block to confirm the trasaction

20.generateblock
Mine a set of ordered transactions to a specified address or descriptor and return the block hash.
bitcoin-cli -datadir=1 generateblock "bcrt1qef9nev7sg0dwjfplx89ecqnv3950w2avkgm86z" '["c6cd113f9b96bc501fd3d903bdec87820db78050cfeb72fddf31f1924b49a328"]'
return:
{
  "hash": "4b9cdd200fec9105edf0cacedb411c0aa692ddecb613801b26e66381dcdec9f0"
}

21.now check balances again
bitcoin-cli -datadir=1 -rpcwallet=anson1 getbalances
{
  "mine": {
    "trusted": 0.00000000,
    "untrusted_pending": 0.00000000,
    "immature": 0.00000000,
    "used": 30.00000000
  }
}

bitcoin-cli -datadir=1 -rpcwallet=anson2 getbalances
{
  "mine": {
    "trusted": 9.99970000,
    "untrusted_pending": 0.00000000,
    "immature": 0.00000000,
    "used": 0.00000000
  }
}
```



## About TRON (version:4.7.0)

### How to run the TRON private chain?

- [ ] if you want to change the initial address for the node, you need to change the addresses in full.conf (quickstart/conf/full.conf), the section is "genesis.block"
- [ ] Run "docker-compose -f docker-compose.tron.private.yaml  up" to start the Tron node server



### Using tronWeb to test the server

- [ ] When the node is running, you can see the tronWeb in the folder "tron", run the command `./tronWeb` to start the command testing

- [ ] Anytime run the command `.exit` to exit the tronWeb

- [ ] There are two genesis address in the full.conf file, now we can use these address to test the transaction

  ```
  1. Get account balance
  > tronWeb.trx.getBalance('TTy7o4hXwuiztVe24EesKAB8haMcE5Keyo').then(result => console.log(result))
  > 95000000000000000
  
  > tronWeb.trx.getBalance('TS89rbaRNZPynZAotL1Zpu4PV3rwzooPCd').then(result => console.log(result))
  > 5000000000000000
    
  2. Create a raw transaction
  // Format
  tronWeb.transactionBuilder.sendTrx(to, amount, from, options);
  
  > tronWeb.transactionBuilder.sendTrx("TVDkStWyB1h7boGpiChry3LukXCwF54qwV", 100, "TS89rbaRNZPynZAotL1Zpu4PV3rwzooPCd").then(result => console.log(result))
  > { visible: false,
    txID:
     '39aa49d44035bb2bfe0c2342d2e4f0d6c5903171c29d690c5f23d7974a43606b',
    raw_data:
     { contract: [ [Object] ],
       ref_block_bytes: '0000',
       ref_block_hash: 'c5870571cda0d2e6',
       expiration: 60000,
       timestamp: 1674546678304 },
  	raw_data_hex:
    '0a0200002208c5870571cda0d2e640e0d4035a65080112610a2d747970652e676f6f676c65617069732e636f6d2f70726f746f636f6c2e5472616e73666572436f6e747261637412300a1541b1332150722e485727c1ee29f0a54769d3c8065b121541d32a766d903993f757d292da6ea4ed1fb3b3fc49186470a0948a96de30' }
     
  3. Sign the raw transaction
  // sign a transaction
  tronWeb.trx.sign(hex string, privateKey);
  
  4. send the signed transaction
  tronWeb.trx.sendRawTransaction(signedTransaction);
  
  5. get the transaction infomation
  tronWeb.trx.getTransaction(transaction id);
  
  6. get the transaction result
  tronWeb.trx.getTransaction("da7d6f7749119f9e76319d893f4a9573f1470bdba6235cc00ffcc8e1e1988eb7").then(console.log)
  
  
  ```



### Using Wallet-cli to test the RPC (default port: 50051)

- [ ] Go to wallet-cli folder and run the command `./gradlew build` to build it, after built, run the command `./gradlew run` to run it.



### Use TronLink to test your private chain

- [ ] Install TronLink to your chrome browser, link: https://chrome.google.com/webstore/detail/tronlink/ibnejdfjmmkpcnlpebklmnkoeoihofec

  ![image-20230303132327945](README.assets/image-20230303132327945.png)

- [ ] Add our private node to your TronLink in the node section of Settings 

  ![image-20230303132424629](README.assets/image-20230303132424629.png)

- [ ] Fill the node name, full node url and event server url

  Node name: TronLocal

  Full node: http://127.0.0.1:9090

  Event server: http://127.0.0.1:9090



### How to request Shasta TRX and Shasta USDT?

Go to discord, and join the tron official channel: https://discordapp.com/invite/hqKvyAM, then open the #faucet conversation, or join the telegram group (https://t.me/TronOfficialDevelopersGroupEn) then:

Sent !shasta [Trx address] to get TRX;

Sent !shasta_usdt [Trx address] to get USDT;



## Some addresses can be used in private chain

**Seed Phrase**

ozone together alpha follow miss wait grid company spider dolphin brief used 

**BIP39 Seed**

05de0a2482a940fc246a4b589103744ec31be822e1e1d189d68ecffc2a25894dd7d6cb7f219800ad989d2069570faed14aee3e5cc2d5a6a6ab80b59e9a0554d7

**BIP32 Root Key**

xprv9s21ZrQH143K3okF3dvf2KzonoSBj6qfftBAtnFksNXEcjMzaySCFQLx6MvzKGrf81y41vQ1HoTF5RMx6E2qT8wZyN5QBKG75PdraWUMrET

**BitCoin**

| Address                            | Public  Key                                                  | Private  Key                                         |
| ---------------------------------- | ------------------------------------------------------------ | ---------------------------------------------------- |
| 17zbWDtJpeHuhmfkMwVy3S8WJDwWurNf6H | 03706f6d3447279ccd8604bb624b110b49b4ed74ce480fe1608247197c9b4ab982 | KztieGKE7oKLoHDbRKU6iYWxBGH9EYPeAXiTUUk9Qx69wbX1edLY |
| 1ERz1kA7v21MEvu9gKnoyMTVd5nEMsWx13 | 03fcc82c8ba2355fdf82c9f21792163877f37491d23e1ed1afd4605c26d3922776 | KzYfmMxm88M3tsNbAhJ6sjHHFjrHjFQSnzY3nF1EtSzX8PVTVnnN |

**ETH**

| Address                                    | Public  Key                                                  | Private  Key                                                 |
| ------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| 0x0344E1f94A4667fc28ceb5d6DFff4CF9b9eE5a10 | 0x0279ef511e15a8eeb5c92d70a758d944878b350c29b242cc2ae7e17b502cd8eaa9 | 0xf7a62cbe00984123edc397511b63b8cd7700ae2c4ca1cb8ce4885ed7de4a0165 |
| 0xca2B0298AC86700f631280bCD21A67b08a6557Ca | 0x0285c08b47897cfafa9f5d7a1a0b3a43e880478caf42562611012aec8987419002 | 0x54e46095f2e0c5f66484fe1836fbecaff4b2b5e706c28eb11f2cb47e7021f299 |

**TRX**

| Address                            | Public  Key                                                  | Private  Key                                                 |
| ---------------------------------- | ------------------------------------------------------------ | ------------------------------------------------------------ |
| TVDkStWyB1h7boGpiChry3LukXCwF54qwV | 032dd6286786b4d9858eaef481af279e18634b684f19102e2e011102d42abb74a6 | 08a89fa29df60d4f0c48a45c58e1dd2664d6bf533ccddac11a0f7514ebc5ff66 |
| TS89rbaRNZPynZAotL1Zpu4PV3rwzooPCd | 039a0742e734ece3b71dc0cae51ebb5c7994af4f150cc795b1db46920571c8f9a6 | b3423dfc973c1986b13860d2079155813401d63e940baa408366d78e7559bcab |



## Reference

Seed Phrase Converter Tool: https://bip39.onekey.so/index.html

Go-Ethereum: https://github.com/ethereum/go-ethereum

Go-Ethereum Docs: https://geth.ethereum.org/docs/getting-started or https://geth.ethereum.org/

Geth RPC API Docs: https://etclabscore.github.io/core-geth/JSON-RPC-API/

Ethereum Development with Go(CN/EN): https://goethereumbook.org/en/

Add Eth to your account in Testnet Goerli: https://goerlifaucet.com/  or https://goerli-faucet.pk910.de/

Goerli Testnet Explorer: https://goerli.etherscan.io/

Testing ETH API tool: https://composer.alchemy.com/

Wei Converter: https://www.alchemy.com/gwei-calculator

USDT-ERC20 Contract Source Code: https://etherscan.io/address/0xdac17f958d2ee523a2206206994597c13d831ec7

BitCoin RPC Docs: https://bitcoincore.org/en/doc/24.0.0/rpc/

BitCoin Private Node: https://github.com/freewil/bitcoin-testnet-box

Tron Development Docs: https://developers.tron.network/docs/getting-start

Tron Wallet-cli: https://github.com/tronprotocol/wallet-cli

Tron tronWeb api: https://developers.tron.network/reference/getaccount

Metamask in chrome: https://chrome.google.com/webstore/detail/metamask/nkbihfbeogaeaoehlefnkodbefgpgknn/related

Decimal converter: https://www.rapidtables.com/convert/number/decimal-to-hex.html

Tron Trc20 contract source code: https://github.com/TRON-Developer-Hub/TRC20-Contract-Template

Getting Shasta TRX: https://www.trongrid.io/shasta/#request

Hex convert to Base58: https://www.btcschools.net/tron/tron_tool_base58check_hex.php

Tron Shasta Browser: https://shasta.tronscan.org/#/

Tron Mainnet Browser: https://tronscan.org/#/
