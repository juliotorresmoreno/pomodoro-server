pomodoro-server
## Como usar?

## Si ya tienes todo el entorno golang configurado
go get -v github.com/juliotorresmoreno/pomodoro-server
cd $GOPATH/src/github.com/juliotorresmoreno/pomodoro-server
go get
go run main.go
Esto inicializa el servidor y lo deja pendiente para escuchar al cliente.

La aplicación usa el puerto 8080 y no requiere permisos de administrador.

## Requerimientos técnicos
1. golang 1.9
2. git

## Instalando las dependencias
### Primero Git
1. sudo apt install git git-core".
2. Configuramos git con github https://git-scm.com/book/es/v1/Empezando-Configurando-Git-por-primera-vez.

### Golang
1. wget https://redirector.gvt1.com/edgedl/go/go1.9.2.linux-amd64.tar.gz
2. sudo tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
3. echo 'export export GOPATH=$HOME/go'>>~/.bashrc
4. echo 'export PATH=$PATH:$GOPATH/bin'>>~/.bashrc
5. sudo su
6. echo 'export PATH=$PATH:/usr/local/go/bin'>>/etc/profile
7. reboot

## Clonando el proyecto

Estos comando se ejecutan sin ser root

1. mkdir $GOPATH/go
2. mkdir $GOPATH/go/src
3. mkdir $GOPATH/go/src/github.com
4. mkdir $GOPATH/go/src/github.com/juliotorresmoreno
5. cd $GOPATH/go/src/github.com/juliotorresmoreno
6. git clone git@github.com:juliotorresmoreno/pomodoro-server.git
7. go get -v

## Si da error con la importación de las dependencias de TiDB

Estos comando se ejecutan sin ser root

1. mkdir $GOPATH/go/src/github.com/pingcap
2. cd $GOPATH/go/src/github.com/pingcap
3. git clone git@github.com:pingcap/tidb.git

## Ejecución
1. cd $GOPATH/go/src/github.com/juliotorresmoreno/juliotorresmoreno
2. go run *.go

