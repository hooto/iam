PROTOC_CMD = protoc
PROTOC_ARGS = --proto_path=./iamapi/ --go_opt=paths=source_relative --go_out=./iamapi/ --go-grpc_out=./iamapi/ ./iamapi/types.proto

HTOML_TAG_FIX_CMD = htoml-tag-fix
HTOML_TAG_FIX_ARGS = iamapi/types.pb.go

BINDATA_CMD = inpack bindata
BINDATA_ARGS_UI = -src webui/ -dst bindata/iam_ws_webui/ -inc js,css,png,svg,ico,tpl,woff2
BINDATA_ARGS_VIEW = -src websrv/views/ -dst bindata/iam_ws_views/ -inc tpl
BINDATA_ARGS_I18N = -src i18n/ -dst bindata/iam_i18n -inc json

BUILDCOLOR="\033[34;1m"
BINCOLOR="\033[37;1m"
ENDCOLOR="\033[0m"

ifndef V
	QUIET_BUILD = @printf '%b %b\n' $(BUILDCOLOR)BUILD$(ENDCOLOR) $(BINCOLOR)$@$(ENDCOLOR) 1>&2;
	QUIET_INSTALL = @printf '%b %b\n' $(BUILDCOLOR)INSTALL$(ENDCOLOR) $(BINCOLOR)$@$(ENDCOLOR) 1>&2;
endif


all: bindata_build
	@echo ""
	@echo "build complete"
	@echo ""

clean: bindata_clean
	@echo ""
	@echo "clean complete"
	@echo ""

proto_build:
	$(QUIET_BUILD)$(PROTOC_CMD) $(PROTOC_ARGS) $(CCLINK)
	$(QUIET_BUILD)$(HTOML_TAG_FIX_CMD) $(HTOML_TAG_FIX_ARGS) $(CCLINK)

bindata_build:
	sed -i 's/debug:\ true/debug:\ false/g' webui/iam/js/main.js
	$(QUIET_BUILD)$(BINDATA_CMD) $(BINDATA_ARGS_UI) $(CCLINK)
	$(QUIET_BUILD)$(BINDATA_CMD) $(BINDATA_ARGS_VIEW) $(CCLINK)
	$(QUIET_BUILD)$(BINDATA_CMD) $(BINDATA_ARGS_I18N) $(CCLINK)
	sed -i 's/debug:\ false/debug:\ true/g' webui/iam/js/main.js

bindata_clean:
	rm -f bindata/iam_ws_webui/statik.go
	rm -f bindata/iam_ws_views/statik.go
	rm -f bindata/iam_ws_i18n/statik.go

