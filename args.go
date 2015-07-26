package xiuxiu

import "flag"

var (
	EsAddr      string
	EsIndice    string
	EsType      string
	EsDebug     bool
	EsAppIndice string
	EsAppType   string
)

func init() {
	flag.StringVar(&EsAddr, "es", "http://[fe80::fabc:12ff:fea2:64a6]:9200", "or http://testbox02.chinacloudapp.cn:9200")
	flag.StringVar(&EsIndice, "indice", "fsmedia2", "target indice")
	flag.StringVar(&EsType, "type", "media", "data type")
	flag.BoolVar(&EsDebug, "debug", true, "diagnose mode")
	flag.StringVar(&EsAppIndice, "app.indice", "ottpomme", "app indice name in es")
	flag.StringVar(&EsAppType, "app.type", "app", "app type name in es")
}
