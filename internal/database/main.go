package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	// host     = "127.0.0.1"
	// port     = 5432

	// host for run in docker
	host     = "host.docker.internal"
	// port database in docker
	port     = 5438
	user     = "postgres"
	password = "123456789"
	dbname   = "learning"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

const HiraganaTable = `
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

const KatakanaTable = `
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