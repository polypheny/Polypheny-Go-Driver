goplay:
	rm -rf build
	mkdir build
	cp -r play protos protoclient polypheny ./build
	cd build/protos && ./generate.sh && go mod init polypheny.com/protos && go mod tidy
	cd ./build/protoclient && go mod init polypheny.com/protoclient && go mod edit -replace polypheny.com/protos=../protos && go mod tidy
	cd ./build/polypheny && go mod init polypheny.com/play && go mod edit -replace polypheny.com/protos=../protos && go mod edit -replace polypheny.com/protoclient=../protoclient && go mod tidy
	cd ./build/play && go mod init polypheny.com/play && go mod edit -replace polypheny.com/protos=../protos && go mod edit -replace polypheny.com/protoclient=../protoclient && go mod edit -replace polypheny.com/polypheny=../polypheny && go mod tidy && go run .

prepush:
	rm -rf build
