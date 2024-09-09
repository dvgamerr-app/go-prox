@echo off
git rev-parse --short main > %TMP%/BUILD.txt
set /p BUILD=<%TMP%/BUILD.txt && rm -f %TMP%/BUILD.txt
set APP=go-prox
go build -ldflags "-s -w -X 'prox/envs.AppName=%APP%' -X 'prox/envs.Version=v0.0.0' -X 'prox/envs.BuildTime=%BUILD%'" -o ./bin/prox.exe