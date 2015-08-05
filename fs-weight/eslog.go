package main

func es_log_url(path, params string) string {
	return "http://" + es_front + path + "?" + params
}
