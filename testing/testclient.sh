#!/bin/sh

sleep 1 # wait for pumba to start

while :
do 
    go run ./cmd -mode=client -ip=$1 < testing/testfile.txt -output=./testing/output.txt
    if diff ./testing/testfile.txt ./testing/output.txt ; then 
        echo "############ Test successful ############"
        sleep 1
    else 
        echo "############ Test failed ############"
        exit 1
    fi
done
