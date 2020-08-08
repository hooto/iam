package main

//go:generate inpack bindata -src webui/ -dst bindata/iam_ws_webui/ -inc js,css,png,svg,ico,tpl,woff2
//go:generate inpack bindata -src websrv/views/ -dst bindata/iam_ws_views/ -inc tpl
//go:generate inpack bindata -src i18n/ -dst bindata/iam_i18n -inc json

/**
go:generate statik -src webui/ -ns iam_ws_webui -dest bindata/ -p iam_ws_webui -f -include *.js,*.css,*.png,*.svg,*.ico,*.tpl,*.woff2
go:generate statik -src websrv/views/ -ns iam_ws_views -dest bindata/ -p iam_ws_views -f -include *.tpl
go:generate statik -src i18n/ -ns iam_i18n -dest bindata/ -p iam_i18n -f -include *.json
*/
