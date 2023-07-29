
#create database and table in postgres. insert sample record to postgres table
docker exec -it db psql -U park_announce -c "CREATE DATABASE IF NOT EXISTS padb;"
docker exec -it db psql -U park_announce -c "CREATE ROLE postgres;"

docker exec -it db psql -U park_announce padb -c "CREATE TABLE IF NOT EXISTS pa_users (id VARCHAR(20) PRIMARY KEY, email VARCHAR(100));"

docker exec -it db psql -U park_announce padb -c "CREATE TABLE IF NOT EXISTS foo ( geog geography );"
docker exec -it db psql -U park_announce padb -c "CREATE INDEX ON foo USING gist(geog);"
docker exec -it db psql -U park_announce padb -c "INSERT INTO foo (geog) VALUES (ST_MakePoint(590454.7399891173, 4519145.719617855));"
docker exec -it db psql -U park_announce padb -c "INSERT INTO foo (geog) VALUES (ST_MakePoint(590250.10594797, 4518558.019924332));"
docker exec -it db psql -U park_announce padb -c "INSERT INTO foo (geog) VALUES (ST_MakePoint(583304.1823994748, 4506069.654048115));"
docker exec -it db psql -U park_announce padb -c "INSERT INTO foo (geog) VALUES (ST_MakePoint(583324.4866324601, 4506805.373160211));"


