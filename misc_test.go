package xiuxiu

import (
	"fmt"
	"testing"
)

func TestEmCleanName(t *testing.T) {
	directors := EmCleanName("约翰·C·瑞利")

	fmt.Println(directors)
}

//EmSplitNameSep
func TestEmSplitNameSep(t *testing.T) {
	directors := EmSplitName("约翰·Ca·瑞利")

	fmt.Println(directors)
}
func TestEmSplitCnEng(t *testing.T) {
	directors := EmSplitCnEng("约翰·Ca·瑞利")

	fmt.Println(directors)
}
func TestEmSplitSep(t *testing.T) {
	directors := EmSplitNameSep("约翰·Ca·瑞利")

	fmt.Println(directors)
}
