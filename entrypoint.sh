#!/bin/sh
hostip=$(ip route show | awk '/default/ {print $3}')
/app/wait_for_response "-url=$1" "-code=$2" "-timeout=$3" "-interval=$4" "-localhost=${hostip}"