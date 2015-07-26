package xiuxiu

import "flag"

var (
	ESAddr   string
	EsIndice string
	EsType   string
)

func init() {
	flag.StringVar(&ESAddr, "es", "http://[fe80::fabc:12ff:fea2:64a6]:9200", "or http://testbox02.chinacloudapp.cn:9200")
	flag.StringVar(&EsIndice, "indice", "fsmedia2", "target indice")
	flag.StringVar(&EsType, "type", "media", "target indice")
}
