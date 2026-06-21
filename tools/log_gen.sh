#!/usr/bin/env bash
set -e

mkdir -p logs

# level: ERROR
# reason: Nostre Dama
# timestamp: 11-12-2013 22:11:02
#

info_log="INFO  generating random logs. Reason: nostre dama"
error_log="ERROR could not find #viu08we9aav.log file. Reason: Please check the volumes are mounted properly"
date_fmt=$(date '+%d-%m-%Y %H:%M:%S')

printf "[%s] %s\n" "$date_fmt" "$error_log" >> ./logs/test_logs.txt
