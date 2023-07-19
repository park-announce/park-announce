
apt update
apt-get install -y wget gcc libpq-dev libgeos-dev proj-bin libproj-dev make libxml2-dev libprotobuf-c-dev protobuf-c-compiler g++ postgresql-server-dev-15
wget http://postgis.net/stuff/postgis-3.4.0dev.tar.gz
tar -xvzf postgis-3.4.0dev.tar.gz
cd postgis-3.4.0dev
chmod +x configure
./configure --without-raster
make
make install