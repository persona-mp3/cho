#!/usr/bin/env bash
set -e

info_log="INFO  generating random logs. Reason: nostre dama"
error_log="ERROR could not find #viu08we9aav.log file. Please check the volumes are mounted properly"
date_fmt=$(date '+%d-%m-%Y %H:%M:%S')

printf "[%s] %s\n" "$date_fmt" "$info_log" >> test_logs.txt
