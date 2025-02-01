#!/bin/sh

# Run swag init to generate Swagger docs
swag init
swag fmt
export POSTGRES_URL=postgres://avnadmin:AVNS_GtJH0K64gPBRu0Oi1KL@pg-3d7cb495-popcart-3578.k.aivencloud.com:22974/defaultdb?sslmode=require
export SECRET=javainuse-secret-key
export SCHEMES=http
export HOST=localhost:8080/api/v1
export AES_SECRET_KEY=pSQXuTQTDyuZHweq9UQZ0RFXSLqr4O4j
export AES_IV=8O0phPKXUQF9RMHb
air

