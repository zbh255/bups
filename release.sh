mkdir release
make build-linux
cd ./build_release ;\
  tar -czvf bups-linux.zip bups cache config.toml ;\
  mv bups-linux.zip ../release ;\
  cd ../
make build-windows
cd ./build_release;\
  tar -czvf bups-windows.zip bups cache config.toml ;\
  mv bups-windows.zip ../release ;\
  cd ../
make clean
