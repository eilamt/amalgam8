#!/bin/bash
exec env JVM_ARGS="-Dhttp.proxyHost=127.0.0.1 -Dhttp.proxyPort=6379" /opt/ibm/wlp/bin/server run defaultServer
