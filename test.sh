#!/bin/bash

cd backend
go test -v

cd ../frontend
npm test

echo "$$$$ TESTS COMPLETE $$$$\n"
