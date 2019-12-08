#!/bin/bash

cd backend
go build . 
./healthchain &

cd ../frontend
npm install 
npm run dev &

echo " $$$$$$$ WELCOME TO HEALTHCHAIN $$$$$$$\n"
