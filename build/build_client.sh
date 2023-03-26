#!/bin/sh
echo create folder for build package ...
if [ ! -e package ]; then
  mkdir package
fi

if [ ! -e package/web ]; then
  mkdir package/web
fi

echo build the code ...
cd ../client
npm install
npm run build
cd ../build

echo remove last package if exist
if [ -e package/web/client ]; then
  rm -rf package/web/client
fi

mv ../client/build ./package/web/client

echo client build over.
