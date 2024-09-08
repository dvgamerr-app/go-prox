@echo off
go build -ldflags "-s -w -X 'app/config.Version=0.0.0' -X 'app/config.BuildTime=0'" -o ./bin/prox.exe