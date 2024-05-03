rm -rf bin/*

mkdir -p bin

cp ./etc/conf/config.$1.yml ./bin/config.yml
cp ./etc/conf/external_api.yml ./bin
cp ./etc/conf/internal_api.yml ./bin
cp ./etc/lumiwavecerti.crt ./bin
cp ./etc/lumiwavecerti.key ./bin

mkdir -p bin/docs/ext

cp ./etc/swagger/ext/*.* ./bin/docs/ext