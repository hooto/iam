PROTOC_CMD = protoc
PROTOC_ARGS = --proto_path=./iamapi/ --go_opt=paths=source_relative --go_out=./iamapi/ --go-grpc_out=./iamapi/ ./iamapi/types.proto

HTOML_TAG_FIX_CMD = htoml-tag-fix
HTOML_TAG_FIX_ARGS = iamapi/types.pb.go

BINDATA_CMD = httpsrv-bindata
BINDATA_ARGS_UI = -src webui/ -dst bindata/iam_ws_webui/ -inc js,css,png,svg,ico,tpl,woff2
BINDATA_ARGS_VIEW = -src websrv/views/ -dst bindata/iam_ws_views/ -inc tpl
BINDATA_ARGS_I18N = -src i18n/ -dst bindata/iam_i18n -inc json


all: bindata_build
	@echo ""
	@echo "build complete"
	@echo ""

clean: bindata_clean
	@echo ""
	@echo "clean complete"
	@echo ""

proto_build:
	$(PROTOC_CMD) $(PROTOC_ARGS)
	$(HTOML_TAG_FIX_CMD) $(HTOML_TAG_FIX_ARGS)

bindata_build:
	$(BINDATA_CMD) $(BINDATA_ARGS_UI)
	$(BINDATA_CMD) $(BINDATA_ARGS_VIEW)
	$(BINDATA_CMD) $(BINDATA_ARGS_I18N)

bindata_clean:
	rm -f bindata/iam_ws_webui/statik.go
	rm -f bindata/iam_ws_views/statik.go
	rm -f bindata/iam_ws_i18n/statik.go

