# Open API Proxy บน Golang

โปรเจกต์นี้เป็น Open API Proxy ที่พัฒนาขึ้นบน Golang ที่มีการ validate และ performance ที่สูง โดยใช้ library ต่อไปนี้:
- [kin-openapi](https://github.com/getkin/kin-openapi)
- [fasthttp](https://github.com/valyala/fasthttp)
- [gocache](https://github.com/eko/gocache)

## คุณสมบัติ

- การ validate ของ Open API Specification
- Performance สูงด้วย fasthttp
- Cache ผลลัพธ์ด้วย gocache

## How to use

1. ติดตั้ง dependencies ด้วย go mod:

```sh
go mod tidy
go run main.go
```

Contributing
หากคุณต้องการเพิ่มฟีเจอร์หรือแก้ไขบั๊ก โปรดส่ง Pull Request มาที่โปรเจกต์นี้

License
โปรเจกต์นี้ถูกลิขสิทธิ์ภายใต้ MIT License
