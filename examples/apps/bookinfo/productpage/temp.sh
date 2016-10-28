#!/bin/bash
exec env HTTP_PROXY=http://localhost:6379 python productpage.py 9080
