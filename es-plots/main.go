package main

import (
	"flag"
	"fmt"
	"net/rpc"
	"strconv"
	"strings"

	"github.com/hearts.zhang/xiuxiu"
	"github.com/hearts.zhang/xiuxiu/seg"
	"github.com/olivere/elastic"
)

var (
	_puncts string = "因为 同样 利用 参与 没有 那么 一栋 一生 一颗 一双 一条 不料 起来 一生 正在 对此 进行 因而 同时 因为 最后 予以 一项 非常 正在 十分 肯定 一切 一匹 一批 一般 没有 一段 对于 不已 之时 一个个 不得已：获得 之一 幸亏 不得不 在外 作为 许多 由此 于是 一番 一名 为了 一份 此外 之间 充满 叫做 年后 带出 时期 透过 解决 必须 突然 仍然 通常 之前 竟然 这是 而且 发生 之外 除了 一步 一发 这部 讲述 身为 描写 接受 到达 一声 人的 期间 伴随 投入 问起 最多 状态 率先 期待 逐渐 组建 时候 梅西 纷纷 找来 加入 由于 锁定 看到 找来  四个 出于 经常 终于 渐渐 回到 原来 虽然 不同 巨大 渺小 变化 一日 过去 这个 纵有 之内 不会 一些 可以 不可能 以致 各种 遇到 一个 开始 可是 总是 加上 她们 位于 渐渐 逐渐 之后 成为 带领 一位 一组 一场 并且 以及 有点 出来 成为 一部 容易 收到 关于 之中 及其 有关 一只 一种 这样 为什么 一直 似乎 后来 更多 下去 看见 已经 走过 加以 两个 根本 他们 可能 很多 这些 视为 知道 其中 那些 必要 还有 不准 如何 这么 拥有 处于 一起 的，在。？！、；：“” ‘’（）─…—·《》【】［］〈〉+-×÷≈＜＞%‰∞∝√∵∴∷∠⊙○π⊥∪°′〃℃{}()[].|&*/#~.,:;?!'-→．"
)

func main() {
	flag.Parse()
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(xiuxiu.EsAddr))
	panic_error(err)

	sego, err := rpc.DialHTTP("tcp", "localhost:8081")
	panic_error(err)

	xiuxiu.EsMediaScan(client, xiuxiu.EsIndice, xiuxiu.EsType, func(em xiuxiu.EsMedia) {
		when_es_media(sego, em)
	})
}
func is_stop_word(seg string) bool {
	return len([]rune(seg)) < 2 || strings.Contains(_puncts, seg)
}
func filter(terms []string) (v []string) {
	for _, t := range terms {
		if !is_stop_word(t) {
			v = append(v, t)
		}
	}
	return
}
func when_es_media(client *rpc.Client, em xiuxiu.EsMedia) {
	if em.Plots == "" {
		return
	}
	var terms []string
	client.Call("Sego.Segment", seg.SegoArg{em.Plots, true}, &terms)
	terms = xiuxiu.EsUniqSlice(filter(terms))
	if len(terms) > 0 {
		fmt.Println(em.Name, strings.Join(terms, " "))
	}
}
func print_es_media(em xiuxiu.EsMedia) {
	fmt.Println(em.Name, f2s(em.Weight), em.MediaLength, em.Day, em.Week, em.Seven, em.Month, em.Play, em.Release)
}
func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}
func f2s(input_num float64) string {
	return strconv.FormatFloat(input_num, 'f', 3, 64)
}
