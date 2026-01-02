#!/bin/zsh
set -x
java -javaagent:crc-agent.jar -cp asm-9.6.jar:asm-commons-9.6.jar:rs317og.jar client 10 0 highmem members 32
