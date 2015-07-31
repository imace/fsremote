package main

type Fuzzy int

func (*Fuzzy) Guess(txt string, suggests *[]int) error {
	*suggests = fuzzy_suggest(txt)
	return nil
}
