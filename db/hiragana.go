package db

import (
	"math/rand"
	"strings"
	"time"
)

var HiraganaTable = `
	あ a	い i	う u	え e	お o
	か ka	き ki	く ku	け ke	こ ko
	が ga	ぎ gi	ぐ gu	げ ge	ご go
	さ sa	し shi	す su	せ se	そ so
	ざ za	じ ji	ず zu	ぜ ze	ぞ zo
	た ta	ち chi	つ tsu	て te	と to
	だ da	ぢ ji	づ zu	で de	ど do
	な na	に ni	ぬ nu	ね ne	の no
	は ha	ひ hi	ふ fu	へ he	ほ ho
	ば ba	び bi	ぶ bu	べ be	ぼ bo
	ぱ pa	ぴ pi	ぷ pu	ぺ pe	ぽ po
	ま ma	み mi	む mu	め me	も mo
	や ya		ゆ yu		よ yo
	ら ra	り ri	る ru	れ re	ろ ro
	わ wa	を wo	ん n/m	
`

func GetCharactersInHiragana() []Word {
	var words []Word
	rows := strings.Split(HiraganaTable, "\n") 
	for _, row := range rows[1:] {
		phrases := strings.Split(row, "\t")[1:]
		for _, phrase := range phrases {
			word := strings.Split(phrase, " ")
			if len(word) > 1 {
				words = append(words, Word{word[0],word[1]})
			}
		}
	}
	return words
}

func GetWordInHiragana(words []Word) Word {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	max := len(words)
	index := r.Intn(max)
	return words[index]
}
