# Todo Service
>
> A simple gRPC Server that provides proto definition to manage todo.

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)

## Prerequisites

- `Go 1.22` or higher.
- Make CLI

## Installation

- Clone this repository.
- Run below command to setup local env variables.
```shell
make env-local
```
- And run below command to setup docker env variables.
```shell
make env-docker
```
- Then run this command to synchronize all the dependencies.
```shell
go mod tidy
```
- Next, setup infrastructure with below command.
```shell
make compose-up
```
- Finally, run the project.

```shell
make local
```