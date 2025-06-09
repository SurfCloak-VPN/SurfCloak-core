# SurfCloak-core documentation
The core itself will be a golang package that you can install as a dependency in your project.

### Install
```
go get https://github.com/SurfCloak-VPN/SurfCloak-core
```
### Update
```
go get -u https://github.com/SurfCloak-VPN/SurfCloak-core
```

### Instructions for use
SurfCloak-core is a golang package that has the necessary methods for generating wireguard configurations. But before generating, you must specify the parameters with which the configuration will be generated. After creating a wg configuration, you can send the client configuration to the user so that he can connect to the server.
