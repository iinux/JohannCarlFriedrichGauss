#!/usr/bin/env bash

#服务停止
ps -ef | grep nsq| grep -v grep | awk '{print $2}' | xargs kill -2