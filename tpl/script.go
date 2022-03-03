/*
* @Author: Rumple
* @Email: ruipeng.wu@cyclone-robotics.com
* @DateTime: 2022/2/28 14:43
 */

package tpl

type DockerBuildParam struct {
	ContainerName string // app-crt
	ImageTag      string // app:latest
	BuildArg      string // --build-arg config=config
	ExportPort    string // 7788
}

const BuildScript = `
docker rm -f {{.ContainerName}}

docker build -t {{.ImageTag}} {{.BuildArg}} .
if ($?)
{
    Write-Host "构建成功"
}
else
{
    Write-Host "构建失败,退出"
    exit
}

docker run -itd --restart=always --name {{.ContainerName}} -p {{.ExportPort}}:8000 {{.ImageTag}}
docker logs -f {{.ContainerName}}
`

const GenSwagger = `swag init --parseDependency --parseInternal --parseDepth 3`

type DockerFileParam struct {
	BinName string
}

const DockerFile = `
FROM golang:latest as builder

WORKDIR /home/works

ADD . .
RUN go build -o bin/{{.BinName}} main.go

FROM alpine:latest

RUN echo "http://mirrors.aliyun.com/alpine/latest-stable/main/" > /etc/apk/repositories \
    && rm -rf /var/cache/apk/* \
    && rm -rf /tmp/* \
    && apk update --allow-untrusted \
    && apk add --no-cache -U tzdata ca-certificates libc6-compat libgcc libstdc++ --allow-untrusted \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata


WORKDIR /home/works

ARG config

EXPOSE 8000
ENV GIN_MODE=release
COPY --from=builder /home/works/bin/ ./
COPY ./config/${config}.yaml ./config/config.yaml
COPY ./static ./static
CMD ./{{.BinName}}
`

type ConfigYamlParam struct {
	Db     string
	DbName string
}

const ConfigYaml = `
{{if eq .Db "mysql"}}
mysql:
  user: "root"
  pass: "mysql"
  host: "localhost"
  port: 3306
  db: "{{.DbName}}"
{{else}}
postgres:
  user: "postgres"
  pass: "postgres"
  host: "localhost"
  port: 5432
  db: "{{.DbName}}"
{{end}}
`
