package prkg

import (
	_ "embed"
)

//go:embed words_bip39_eng.txt
var wordListEnglish string
var DictionaryEnglish = NewDictionary(wordListEnglish)
