create_file:
	mkdir cache
	touch access.log
	touch bups.log
	cp ../config/dev/config.toml ./config.toml
	touch errors.log
	cp ./config.toml test_read.toml
	touch test_write.toml

delete_file:
	rm -rf cache
	rm -rf access.log
	rm -rf bups.log
	rm -rf error.log
	rm -rf config.toml
	rm -rf errors.log
	rm -rf test_read.toml
	rm -rf test_write.toml

test:create_file
test:delete_file
	go test -race -v .
	make delete_file