# Prox

golang windows services and proxy server.

![terminal](./docs/terminal.png)

```bash

# make windows service 
./build.bat

sc.exe create "go-prox" "./bin/prox.exe"
sc.exe start
```

## Ref
- https://clickhouse.com/blog/storing-traces-and-spans-open-telemetry-in-clickhouse
- https://opentelemetry.io/docs/collector/quick-start/
- https://chatgpt.com/share/1ffa58a1-76c5-452f-8797-7b320eec1922
