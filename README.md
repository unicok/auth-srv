# auth
[![Build Status](https://travis-ci.org/unicok/auth-srv.svg?branch=master)](https://travis-ci.org/unicok/auth-srv)

## Getting started

1. Install Consul

	Consul is the default registry/discovery for go-micro apps. It's however pluggable.
	[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

2. Run Consul
	```
	$ consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul
	```

4. Download and start the service

	```shell
	go get github.com/unicok/auth-srv
	./auth-srv --mongodb_url="mongodb://127.0.0.1:27017/account"
	```


## 设计理念
用户中心，支持各种第三方登陆，暂未实现。