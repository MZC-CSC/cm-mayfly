
## `cm-mayfly`의 `Docker Compose`를 이용한 Cloud-Migrator 설치 및 실행 가이드

이 가이드에서는 `cm-mayfly`의 `docker` 서브 커맨드를 이용하여 `Docker Compose` 기반으로 Cloud-Migrator 시스템을 구축 및 실행하는 방법에 대해 소개합니다. 


## 순서
1. 개발환경 준비
1. 필요사항 설치
   1. Docker
   1. Docker Compose
1. cm-mayfly 소스코드 다운로드
1. 환경설정 확인 및 변경
1. cm-mayfly 소스코드 빌드
1. cm-mayfly 이용하여 Cloud-Migrator 실행
1. Cloud-Migrator 실행상태 확인
1. [참고] 프레임워크별 컨테이너 구성 및 API Endpoint


## 개발환경 준비

[권장사항]
- Ubuntu 20.04

## 필요사항 설치

### Docker 설치
- https://docs.docker.com/engine/install/ubuntu/ 에서 설명하는 방법대로 설치합니다.

<details>
  <summary>[클릭하여 예시 보기]</summary>
  
```bash
# 기존에 Docker 가 설치되어 있었다면 삭제
sudo apt remove docker docker-engine docker.io containerd runc

# Docker 설치를 위한 APT repo 추가
sudo apt update

sudo apt install \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# x86_64 / amd64
echo \
  "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt update

sudo apt install docker-ce docker-ce-cli containerd.io
```
</details>

### Docker Compose 설치
- APT 패키지 매니저를 이용하여 설치합니다.
```bash
sudo apt install docker-compose
```

## cm-mayfly 소스코드 다운로드
```bash
git clone https://github.com/cm-mayfly/cm-mayfly.git
```

## 도커 서비스 정의 확인
Cloud-Migrator 시스템 구성에 필요한 서비스 정보는 `conf/docker` 폴더 하위에 정의되어 있으니 `cm-mayfly/conf/docker/docker-compose.yaml` 파일 및 `conf/docker/conf` 폴더의 내용들을 살펴 보고 필요한 경우 수정합니다.

## cm-mayfly 실행파일 또는 소스코드 빌드
bin 폴더에 최신 실행 파일이 있으니 사용하기 바라며, 소스 빌드가 필요한 경우에는 README에 설명된 go 설치 방법과 build 명령 또는 make 명령의 빌드 방법을 참고합니다.
```bash
cd cm-mayfly/src
go build -o mayfly main.go
```

## cm-mayfly 이용하여 Cloud-Migrator 실행
```bash
./bin/mayfly docker run
```

## Cloud-Migrator 실행상태 확인
```bash
./bin/mayfly docker info
```


## Cloud-Migrator 중지
```bash
./bin/mayfly docker stop
```

## [참고] 프레임워크별 컨테이너 구성 및 API Endpoint
| Framework별 Container Name | REST-API Endpoint |
|---|---|
| cb-spider | http://{{host}}:1024/spider |
| --- |
| cb-tumblebug | http://{{host}}:1323/tumblebug |
| --- |

