#!/bin/zsh
set -x
javac -cp asm-9.6.jar:asm-commons-9.6.jar CrcPrintAgent.java
rm -f crc-agent.jar
jar cvmf MANIFEST.MF crc-agent.jar CrcPrintAgent*.class
