test-e2e:
	go test \
	-race \
	./test/e2e \
	-timeout "10m" \
	-v \
	-ginkgo.v \
	-ginkgo.noColor 