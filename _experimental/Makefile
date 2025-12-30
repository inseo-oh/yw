DOWNLOAD_DIR := downloads
HTML_ENTITIES_JSON := $(DOWNLOAD_DIR)/entities.json
HTML_ENTITIES_AUTOGEN_I := yw/yw_html_entities_autogen.i

all: prepare yw yw-all-platforms

$(HTML_ENTITIES_JSON):
	wget https://html.spec.whatwg.org/entities.json -O $@

prepare:
	mkdir -p $(DOWNLOAD_DIR)

yw: prepare autogen-entities
	$(MAKE) -C yw all DEBUG=$(DEBUG) PLATFORM=$(PLATFORM)

yw-all-platforms: prepare autogen-entities
	$(MAKE) -C yw all-platforms DEBUG=$(DEBUG) PLATFORM=$(PLATFORM)

autogen-entities: prepare $(HTML_ENTITIES_JSON)
	$(MAKE) -C html_entities_gen DEBUG=$(DEBUG)
	./out_linux/html_entities_gen $(HTML_ENTITIES_JSON) $(HTML_ENTITIES_AUTOGEN_I)

clean:
	$(MAKE) -C html_entities_gen clean
	$(MAKE) -C yw clean

clean-all-platforms:
	$(MAKE) -C html_entities_gen clean
	$(MAKE) -C yw clean

.PHONY: all prepare yw yw-all-platforms autogen-entities clean clean-all-platforms
