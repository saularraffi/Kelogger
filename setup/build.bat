@echo off

set emailFrom=%1
set emailTo=%2
set emailPassword=%3
set reportInterval=%4
set sourceLocation=%5
set outputExe=%6

go build -ldflags "-H=windowsgui -X main.emailTo=%emailTo% -X main.emailFrom=%emailFrom% -X main.emailPassword=%emailPassword% -X main.reportInterval=%reportInterval%" -o %outputExe% %sourceLocation%