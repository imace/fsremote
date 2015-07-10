package main

import (
	"fmt"
)

func main() {
	model := NewModel()

	// For testing only, this is not advisable on production
	model.SetThreshold(1)

	// This expands the distance searched, but costs more resources (memory and time).
	// For spell checking, "2" is typically enough, for query suggestions this can be higher
	model.SetDepth(1)

	// Train multiple words simultaneously by passing an array of strings to the "Train" function
	words := []string{"刘", "刘德华", "德华", "德行", "华德育", "行天下之中", "下之中华德育", "剪刀手爱德华", "爱德华", "刀手爱"}
	//model.Train(words)
	for i, w := range words {
		model.SetCount(w, i, true)
	}
	// Train word by word (typically triggered in your application once a given word is popular enough)
	model.TrainWord("single")

	// Suggest completions
	fmt.Println("\nQUERY SUGGESTIONS")
	fmt.Println(model.Suggestions("刘德", true))
	fmt.Println(model.Suggestions("bo", false))
	fmt.Println(model.Suggestions("dyn", false))

}
