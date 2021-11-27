#!/bin/bash

./kuirejo start-execution \
	--asl  "./workflows/HelloWorld2/statemachine.asl.json" \
	--input "./workflows/HelloWorld2/input1.json"
