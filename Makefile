include Makefile.conf


SPAWNER_LIST = $(WORKDIR)/spawners.list
SPAWNOPTION_LIST = $(WORKDIR)/spawnoptions.list


.PHONY: serialize-spawners
serialize-spawners:
	OUTPUT_FILE=$(SPAWNER_LIST) \
	GAMEDATA_PATH=$(GAMEDATA_PATH) \
	JWP_SERIALIZE=$(JWP_SERIALIZE) \
	JWP_CMD="$(JWP_CMD)" \
	gamedata-collector/collect-spawners.sh


.PHONY: serialize-spawnoptions 
serialize-spawnoptions:
	OUTPUT_FILE=$(SPAWNOPTION_LIST) \
	GAMEDATA_PATH=$(GAMEDATA_PATH) \
	JWP_CMD="$(JWP_CMD)" \
	gamedata-collector/collect-spawnoptions.sh


.PHONY: serialize-gamedata
serialize-gamedata: serialize-spawners serialize-spawnoptions


.PHONY: generate-spawners
generate-spawners:
	python gamedata-collector/generate.py \
	--golang \
	--output gamedata/generated_spawners.go \
	spawners list \
	$(SPAWNER_LIST)


.PHONY: generate-uncap-mod
generste-uncap-mod:
	python gamedata-collector/generate.py \
	--golang \
	--output generator/generated_uncap_mod.go \
	spawnoptions uncap \
	$(SPAWNOPTION_LIST)


.PHONY: generate-gamedata
generate-gamedata: generate-spawners generate-uncap-mod


.PHONY: cli
cli:
	go build -o build/mayham-cli cmd/mayham-cli/main.go


.PHONY: web
web:
	go build -o build/mayham-web cmd/mayham-web/main.go


.PHONY: cli-full
cli-full: generate-gamedata cli


.PHONY: all
all: serialize-gamedata generate-gamedata cli
