#!/usr/bin/env python3

import pymysql
import redis
import subprocess
import time
import requests
import datetime
import json
import socket
from behave import *


register_type(int=int)
register_type(str=lambda x: x if x != "N/A" else "")
register_type(bool=lambda x: True if x == "true" else False)


config = {
    "prefix": "output/go-tech",
    "service": {
        "port": 16063,
        "allowOrigins": [
            "http://127.0.0.1:4000"
            "http://127.0.0.1:4001"
        ],
    },
    "api": {
        "account": "test-go-account:16060",
    },
    "mysql": {
        "host": "test-mysql",
        "port": 3306,
        "user": "hatlonely",
        "password": "keaiduo1",
        "db": "article"
    },
    "redis": {
        "host": "test-redis",
        "port": 6379
    }
}


def wait_for_port(port, host="localhost", timeout=5.0):
    start_time = time.perf_counter()
    while True:
        try:
            with socket.create_connection((host, port), timeout=timeout):
                break
        except OSError as ex:
            time.sleep(0.01)
            if time.perf_counter() - start_time >= timeout:
                raise TimeoutError("Waited too long for the port {} on host {} to start accepting connections.".format(
                    port, host
                )) from ex


def deploy():
    fp = open("{}/configs/tech.json".format(config["prefix"]))
    cf = json.loads(fp.read())
    fp.close()
    cf["service"]["port"] = ":{}".format(config["service"]["port"])
    cf["service"]["allowOrigins"] = config["service"]["allowOrigins"]
    cf["api"]["account"] = config["api"]["account"]
    cf["mysql"]["uri"] = "{user}:{password}@tcp({host}:{port})/{db}?charset=utf8&parseTime=True&loc=Local".format(
        user=config["mysql"]["user"],
        password=config["mysql"]["password"],
        db=config["mysql"]["db"],
        host=config["mysql"]["host"],
        port=config["mysql"]["port"],
    )
    print(cf)
    fp = open("{}/configs/tech.json".format(config["prefix"]), "w")
    fp.write(json.dumps(cf, indent=4))
    fp.close()


def start():
    subprocess.Popen(
        "cd {} && nohup bin/tech &".format(config["prefix"]),  shell=True
    )

    wait_for_port(config["service"]["port"], timeout=5)


def stop():
    subprocess.getstatusoutput(
        "ps aux | grep bin/tech | grep -v grep | awk '{print $2}' | xargs kill"
    )


def before_all(context):
    config["url"] = "http://127.0.0.1:{}".format(config["service"]["port"])
    deploy()
    start()
    context.config = config
    context.mysql_conn = pymysql.connect(
        host=config["mysql"]["host"],
        user=config["mysql"]["user"],
        port=config["mysql"]["port"],
        password=config["mysql"]["password"],
        db=config["mysql"]["db"],
        charset="utf8",
        cursorclass=pymysql.cursors.DictCursor
    )
    context.redis_client = redis.Redis(
        config["redis"]["host"], port=6379, db=0
    )


def after_all(context):
    stop()
