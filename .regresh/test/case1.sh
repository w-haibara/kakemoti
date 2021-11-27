#!/bin/bash

./kuirejo start-execution \
    --asl  "./workflows/HelloWorld/statemachine.asl.json" \
	--input "./workflows/HelloWorld/input1.json"
