#!/usr/bin/env bash

begin_path=$PWD
echo "========== Start generating swag doc for admin server ==========="

#cd $begin_path
#cd ../cmd/admin_api
swag init -d ../cmd/admin_api --parseDependency --parseInternal -o ../pkg/swagger/admin_api

echo "========== Generate admin server success!! ==========="

echo "========== Start generating swag doc for wallet server ==========="

#cd $begin_path
#cd ../cmd/wallet_api
swag init -d ../cmd/wallet_api --parseDependency --parseInternal -o ../pkg/swagger/wallet_api

echo "========== Generate wallet server success!! ==========="