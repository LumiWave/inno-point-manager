# set -x

rm -rf bin/*

mkdir -p bin

cp ./etc/conf/config.stage.yml ./bin/config.yml
cp ./etc/conf/external_api.yml ./bin
cp ./etc/conf/internal_api.yml ./bin
cp ./etc/onbuffcerti.crt ./bin
cp ./etc/onbuffcerti.key ./bin

mkdir -p bin/docs/ext

cp ./etc/swagger/ext/*.* ./bin/docs/ext

rm -rf bin/inno-point-manager

go build -o bin/inno-point-manager.exe main.go

cd bin
./inno-point-manager.exe -c=config.yml