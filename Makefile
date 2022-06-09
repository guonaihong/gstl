all: fmt test

# lvim没有很好适配泛型语, 不能自动格式化. 这里先手动执行下
bug:
	#go test -test.run=Test_Btree_RangePrev ./...

fmt:
	go fmt ./...
test:
	go test ./...

