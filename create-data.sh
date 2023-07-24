
#create database and table in postgres. insert sample record to postgres table
docker exec -it db psql -U park_announce -c "CREATE DATABASE IF NOT EXISTS padb;"
docker exec -it db psql -U park_announce -c "CREATE ROLE postgres;"
docker exec -it db psql -U park_announce padb -c "CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, phone VARCHAR(20), nickname VARCHAR(50), name VARCHAR(50), surname VARCHAR(50), email VARCHAR(100));"
docker exec -it db psql -U park_announce padb -c "CREATE TABLE IF NOT EXISTS geolocations (id SERIAL PRIMARY KEY, user_id INTEGER, latitude FLOAT, longitude FLOAT, FOREIGN KEY (user_id) REFERENCES users (id));"
docker exec -it db psql -U park_announce padb -c "INSERT INTO users(phone, nickname, name, surname, email) VALUES ('123456789', 'user1', 'John', 'Doe', 'john.doe@example.com'),('987654321', 'user2', 'Jane', 'Smith', 'jane.smith@example.com');"
docker exec -it db psql -U park_announce padb -c "INSERT INTO geolocations (user_id, latitude, longitude) VALUES (1, 37.7749, -122.4194),(2, 34.0522, -118.2437);"

docker exec -it db psql -U park_announce padb -c "CREATE TABLE IF NOT EXISTS foo ( geog geography );"
docker exec -it db psql -U park_announce padb -c "CREATE INDEX ON foo USING gist(geog);"
docker exec -it db psql -U park_announce padb -c "INSERT INTO foo (geog) VALUES (ST_MakePoint(590454.7399891173, 4519145.719617855));"
docker exec -it db psql -U park_announce padb -c "INSERT INTO foo (geog) VALUES (ST_MakePoint(590250.10594797, 4518558.019924332));"
docker exec -it db psql -U park_announce padb -c "INSERT INTO foo (geog) VALUES (ST_MakePoint(583304.1823994748, 4506069.654048115));"
docker exec -it db psql -U park_announce padb -c "INSERT INTO foo (geog) VALUES (ST_MakePoint(583324.4866324601, 4506805.373160211));"


