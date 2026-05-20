.PHONY: build run migrate clean vet

BINARY = tmf915-api.exe
CONFIG = config.yaml

build:
	go build -o $(BINARY) ./cmd/server

run: build
	./$(BINARY) $(CONFIG)

migrate:
	@echo "Running migrations against DESKTOP-V2M1TS5..."
	sqlcmd -S DESKTOP-V2M1TS5 -E -i migrations\001_create_database.sql
	sqlcmd -S DESKTOP-V2M1TS5 -d TMF915 -E -i migrations\002_create_aimodelspecification.sql
	sqlcmd -S DESKTOP-V2M1TS5 -d TMF915 -E -i migrations\003_create_aimodel.sql
	sqlcmd -S DESKTOP-V2M1TS5 -d TMF915 -E -i migrations\004_create_aicontractspecification.sql
	sqlcmd -S DESKTOP-V2M1TS5 -d TMF915 -E -i migrations\005_create_aicontract.sql
	sqlcmd -S DESKTOP-V2M1TS5 -d TMF915 -E -i migrations\006_create_aicontractviolation.sql
	sqlcmd -S DESKTOP-V2M1TS5 -d TMF915 -E -i migrations\007_create_alarm.sql
	sqlcmd -S DESKTOP-V2M1TS5 -d TMF915 -E -i migrations\008_create_rule.sql
	sqlcmd -S DESKTOP-V2M1TS5 -d TMF915 -E -i migrations\009_create_monitor.sql
	sqlcmd -S DESKTOP-V2M1TS5 -d TMF915 -E -i migrations\010_create_topic.sql
	sqlcmd -S DESKTOP-V2M1TS5 -d TMF915 -E -i migrations\011_create_event.sql
	sqlcmd -S DESKTOP-V2M1TS5 -d TMF915 -E -i migrations\012_create_hub_listener.sql
	@echo "Migrations complete."

vet:
	go vet ./...

clean:
	del /f $(BINARY) 2>nul || true
