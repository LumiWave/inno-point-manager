set -x

sh ./prebuild.sh

go build -o bin/inno_point_manager rest_server/main.go

cd bin
./inno_point_manager -c=config.yml