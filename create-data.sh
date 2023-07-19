
#creare database and table in postgres. insert sample record to postgres table
docker exec -it db psql -U park-announce -c "create database padb;"
docker exec -it db psql -U park-announce padb -c "CREATE TABLE users (id SERIAL PRIMARY KEY, phone VARCHAR(20), nickname VARCHAR(50), name VARCHAR(50), surname VARCHAR(50), email VARCHAR(100));"
docker exec -it db psql -U park-announce padb -c "CREATE TABLE geolocations (id SERIAL PRIMARY KEY, user_id INTEGER, latitude FLOAT, longitude FLOAT, FOREIGN KEY (user_id) REFERENCES users (id));"
docker exec -it db psql -U park-announce padb -c "INSERT INTO users(phone, nickname, name, surname, email) VALUES ('123456789', 'user1', 'John', 'Doe', 'john.doe@example.com'),('987654321', 'user2', 'Jane', 'Smith', 'jane.smith@example.com');"
docker exec -it db psql -U park-announce padb -c "INSERT INTO geolocations (user_id, latitude, longitude) VALUES (1, 37.7749, -122.4194),(2, 34.0522, -118.2437);"