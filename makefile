.PHONY: home driver g g2 exec

# If the first argument is "run"...
ifeq (home,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  args := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(args):;@:)
endif

home:
	go test -v -count=1 test/setup_test.go test/home_test.go -run TestHome $(args)

driver:
	go test -v -count=1 test/setup_test.go test/driver_test.go -run TestDriver

g:
	go test -v -count=1 test/setup_test.go test/g_test.go -run TestG

g2:
	go test -v -count=1 test/setup_test.go test/g_test.go -run TestG2

exec:
	go test -v -count=1 test/setup_test.go test/$(t)_test.go -run $(tn)

