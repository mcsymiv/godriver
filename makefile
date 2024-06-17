home:
	go test -v -count=1 test/setup_test.go test/home_test.go -run TestHome

exec:
	go test -v -count=1 test/setup_test.go test/$(t)_test.go -run $(tn)

