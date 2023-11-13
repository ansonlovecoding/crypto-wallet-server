#!/usr/bin/env bash
echo "docker-compose ps..........................."
docker-compose ps

echo "check wallet, waiting 3s...................."
sleep 3

echo "check wallet................................"
./check_all.sh


