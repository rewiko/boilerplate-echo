PID_FILE = /tmp/boilerplate-echo.pid
GO_FILES = main.go
APP      = ./boilerplate-echo
serve: restart
	@fswatch -or --event=Updated -e ".*" -i ".\.go$$" . | xargs -n1 -I{}  make restart || make kill

kill:
	@echo "### Killing server ###"
	@-if [ -f "$(PID_FILE)" ]; then \
	  kill `pstree -p \`cat $(PID_FILE)\` | tr "\n" " " |sed "s/[^0-9]/ /g" |sed "s/\s\s*/ /g"` > /dev/null 2>&1; rm $(PID_FILE) || true ; \
	fi

before:
	@echo "### Starting server ###" && printf '%*s\n' "40" '' | tr ' ' -
# Start task performs "go run main.go" command and writes it's process id to PID_FILE.
start:
	go run $(GO_FILES) & echo $$! > $(PID_FILE)
restart: kill before start

.PHONY: serve restart kill before
