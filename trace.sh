#!/bin/sh

set -e

./make.sh tracing

WRK=$HOME/Projects/thirdparty/wrk/

./shortener &
sleep 0.5
$WRK/wrk -t 1 -c 1 -d 5 --script $WRK/plaintext.lua http://localhost:7075/plaintext -- 16
killall shortener
