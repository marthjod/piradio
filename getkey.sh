#!/bin/bash

# the input device (a numpad)
INPUT_DEVICE=/dev/input/event0

# find these out from "echo $KEY" below
KEY_0="0062"
KEY_1="0059"
KEY_2="005a"
KEY_3="005b"
KEY_4="005c"
KEY_5="005d"
KEY_6="005e"
KEY_7="005f"
KEY_8="0060"
KEY_9="0061"
KEY_COMMA="0063"
KEY_ENTER="0058"
KEY_PLUS="0057"
KEY_MINUS="0056"
KEY_BACKSPACE="002a"
KEY_SPLAT="0055"
KEY_SLASH="0054"

# we'll exit signalling the key that has been pressed "encoded" in the exit status code
function getkey() {
	KEY=$(dd if=$INPUT_DEVICE bs=16 count=1 2> /dev/null | od -x | grep -v 0000020 | cut -b 39-42)
	case $KEY in
		$KEY_0) 
			exit 0
		;;
		$KEY_1)
			exit 1
		;;
		$KEY_2)
			exit 2
		;;
		$KEY_3)
			exit 3
		;;
		$KEY_4)
			exit 4
		;;
		$KEY_5)
			exit 5
		;;
		$KEY_6)
			exit 6
		;;
		$KEY_7)
			exit 7
		;;	
		$KEY_8)
			exit 8
		;;
		$KEY_9)
			exit 9
		;;											
		$KEY_ENTER)
			exit 10
		;;	
		$KEY_PLUS)
			exit 11
		;;
		$KEY_MINUS)
			exit 12
		;;
		$KEY_BACKSPACE)
			exit 13
		;;			
		$KEY_SPLAT)
			exit 14
		;;	
		$KEY_SLASH)
			exit 15
		;;
		$KEY_COMMA)
			exit 16
		;;											
	esac
}

while :;
do
	getkey
done
