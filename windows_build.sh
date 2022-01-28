# set -x

sh ./prebuild.sh $1

rm -rf bin/inno-point-manager

go build -o bin/inno-point-manager.exe main.go

cd bin
./inno-point-manager.exe -c=config.yml