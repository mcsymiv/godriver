.PHONY: home driver g

home:
	go test -v -count=1 test/setup_test.go test/home_test.go -run TestHome

driver:
	go test -v -count=1 test/setup_test.go test/driver_test.go -run TestDriver

g:
	go test -v -count=1 test/setup_test.go test/g_test.go -run TestG

exec:
	go test -v -count=1 test/setup_test.go test/$(t)_test.go -run $(tn)

