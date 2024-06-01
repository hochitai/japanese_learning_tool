package database

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
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

var katakanaTable = `
	ア a	イ i	ウ u	エ e	オ o
	カ ka	キ ki	ク ku	ケ ke	コ ko
	ガ ga	ギ gi	グ gu	ゲ ge	ゴ go
	サ sa	シ shi	ス su	セ se	ソ so
	ザ za	ジ ji	ズ zu	ゼ ze	ゾ zo
	タ ta	チ chi	ツ tsu	テ te	ト to
	ダ da	ヂ ji	ヅ zu	デ de	ド do
	ナ na	ニ ni	ヌ nu	ネ ne	ノ no
	ハ ha	ヒ hi	フ fu	ヘ he	ホ ho
	バ ba	ビ bi	ブ bu	ベ be	ボ bo
	パ pa	ピ pi	プ pu	ペ pe	ポ po
	マ ma	ミ mi	ム mu	メ me	モ mo
	ヤ ya		ユ yu		ヨ yo
	ラ ra	リ ri	ル ru	レ re	ロ ro
	ワ wa	ヲ wo	ン n/m	
`

type Word struct {
	Id 				int    `json:"id"`
	Characters 		string `json:"characters"`
	Pronunciation   string `json:"pronunciation"`
	Meaning			string `json:"meaning"`
}

func (w *Word) SetId(id int) {
	w.Id = id
}

// var defaultVocabulary = []WordModel{
// 	{1,"ありがとう","arigatoo","Cảm ơn"},
// }

func GetHiraganaCharacters() []Word {
	return getCharacters(HiraganaTable)
}

func GetKatakanaCharacters() []Word {
	return getCharacters(katakanaTable)
}

func GetRandomCharacter(words []Word) Word {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	max := len(words)
	index := r.Intn(max)
	return words[index]
}

func getCharacters(alphabet string) []Word {
	var words []Word
	rows := strings.Split(alphabet, "\n") 
	for _, row := range rows[1:] {
		phrases := strings.Split(row, "\t")[1:]
		for _, phrase := range phrases {
			word := strings.Split(phrase, " ")
			if len(word) > 1 {
				words = append(words, Word{
					Characters: word[0],
					Pronunciation: word[1],
				})
			}
		}
	}
	return words
}

func GetVocabularies(db *pg.DB) ([]Word, error) {
	var words []Word
	err := db.Model(&words).Select()
	if err != nil {
		return nil, err
	}
	return words, nil
}

func GetWords(db *pg.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		words, err := GetVocabularies(db)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusNotFound, struct{}{})
			return
		}
		c.JSON(http.StatusOK, words)
	}

}

func AddWord(db *pg.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		decoder := json.NewDecoder(c.Request.Body)
		var word Word
		err := decoder.Decode(&word)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, struct{}{})
			return
		}

		_, err = db.Model(&word).Returning("*").Insert()
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, struct{}{})
			return
		}
		c.JSON(http.StatusCreated, word)
	}
}

func UpdateWord(db *pg.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Error(err)		
			c.JSON(http.StatusBadRequest, struct{}{})
			return
		}
		decoder := json.NewDecoder(c.Request.Body)
		var word Word
		err = decoder.Decode(&word)
		if err != nil {
			c.Error(err)	
			c.JSON(http.StatusBadRequest, struct{}{})
			return
		}
		word.SetId(id)

		_, err = db.Model(&word).WherePK().Returning("*").Update()
		if err != nil {
			c.Error(err)	
			c.JSON(http.StatusBadRequest, struct{}{})
			return
		}
		c.JSON(http.StatusOK, word)
	}
}

func DeleteWord(db *pg.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Error(err)	
			c.JSON(http.StatusBadRequest, struct{}{})
			return
		}
		decoder := json.NewDecoder(c.Request.Body)
		var word Word
		err = decoder.Decode(&word)
		if err != nil {
			c.Error(err)	
			c.JSON(http.StatusBadRequest, struct{}{})
			return
		}
		word.SetId(id)

		_, err = db.Model(&word).WherePK().Returning("*").Delete()
		if err != nil {
			c.Error(err)	
			c.JSON(http.StatusBadRequest, struct{}{})
			return
		}
		c.JSON(http.StatusOK, word)
	}
}