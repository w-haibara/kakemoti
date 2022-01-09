#!/bin/bash

outprefix="KAKEMOTI_OUT_"
errprefix="KAKEMOTI_ERR"

tmpfiledir="/tmp/kakemoti/script2"
tmpfile=${tmpfiledir}"/tmp"

if [[ -f "${tmpfile}" && `cat ${tmpfile}` = "" ]]
    then
        rm ${tmpfile}
    else
        mkdir -p ${tmpfiledir} && touch ${tmpfile}
        echo -n a >> ${tmpfile}
fi

if [[ `cat ${tmpfile}` = "aaaaa" ]]
    then
        echo ${outprefix}"Payload=OK"
        echo -n "" > ${tmpfile}
    else
        echo ${errprefix}"=SCRIPT.DUMMY.ERROR"
fi
