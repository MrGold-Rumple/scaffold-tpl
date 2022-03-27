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

const PowerBuildScript = `docker rm -f {{.ContainerName}}

docker build -t {{.ImageTag}} {{.BuildArg}} .
if ($?)
{
    Write-Host "build success ~~"
}
else
{
    Write-Host "build failed !!"
    exit
}

docker run -itd --restart=always --name {{.ContainerName}} -p {{.ExportPort}}:8000 {{.ImageTag}}
docker logs -f {{.ContainerName}}
`

const BashBuildScript = `docker rm -f {{.ContainerName}}

docker build -t {{.ImageTag}} {{.BuildArg}} .

if [[ "$?" != "0" ]];then
	echo "build failed !!"
else
	echo "build success ~~"
fi

docker run -itd --restart=always --name {{.ContainerName}} -p {{.ExportPort}}:8000 {{.ImageTag}}
docker logs -f {{.ContainerName}}
`

const GenSwagger = `swag init --parseDependency --parseInternal --parseDepth 3`

type DockerFileParam struct {
	BinName string
}

const DockerFile = `#syntax=docker/dockerfile:latest
FROM golang:latest as builder

WORKDIR /app

ADD . .

RUN --mount=type=cache,id=go_mod,target=/go/pkg/mod \
    --mount=type=cache,id=odp_go_cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o bin/{{.BinName}} ./main.go

FROM alpine:latest

RUN echo "http://mirrors.aliyun.com/alpine/latest-stable/main/" > /etc/apk/repositories \
    && apk update \
    && apk add tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

WORKDIR /home/works/program

ARG config=config

EXPOSE 8000

ENV GIN_MODE=release
COPY --from=builder /app/bin ./
COPY ./config/${config}.yaml ./config/config.yaml
CMD ./{{.BinName}}
`

type ConfigYamlParam struct {
	Db     string
	DbName string
}

const ConfigYaml = `database:
  type: "{{.Db}}"
  user: "admin"
  pass: "admin"
  host: "localhost"
  port: 9527
  db: "{{.DbName}}"
`
