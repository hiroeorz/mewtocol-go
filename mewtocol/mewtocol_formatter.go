package mewtocol

import (
	"fmt"
	"log"
	"strings"
)

// 指定した1点の接点状態を読み出す電文を生成して返す
// contactCode:
//    X: 外部入力
//    Y: 外部出力
//    R: 内部リレー
//    L: リンクリレー
//    T: タイマ
//    C: カウンタ
func formatReadIOSingle(dstAddress uint, contactCode string, contactNo uint) string {
	if !isValidCode(contactCode, []string{"X", "Y", "R", "L", "T", "C"}) {
		panic(fmt.Sprintln("invalid code:", contactCode))
	}

	command := "RCS" + contactCode + fmt.Sprintf("%04d", contactNo)
	return format(dstAddress, command)
}

// 指定したワード単位での接点状態を読み出す電文を生成して返す
// contactCode:
//    X: 外部入力
//    Y: 外部出力
//    R: 内部リレー
//    L: リンクリレー
//    T: タイマ
//    C: カウンタ
func formatReadIOWord(dstAddress uint, contactCode string, startWordNo uint, endWordNo uint) string {
	if !isValidCode(contactCode, []string{"X", "Y", "R", "L", "T", "C"}) {
		panic(fmt.Sprintln("invalid code:", contactCode))
	}

	command := "RCC" + contactCode + fmt.Sprintf("%04d", startWordNo) + fmt.Sprintf("%04d", endWordNo)
	return format(dstAddress, command)
}

// 指定した1点の接点出力を書き込む電文を生成して返す。
// stateがtrueの時はON、falseの時はOFFを出力します。
// contactCode:
//    Y: 外部出力
//    R: 内部リレー
//    L: リンクリレー
func formatWriteIOSingle(dstAddress uint, contactCode string, contactNo uint, state bool) string {
	if !isValidCode(contactCode, []string{"Y", "R", "L"}) {
		panic(fmt.Sprintln("invalid code:", contactCode))
	}
	contactData := ""

	if state {
		contactData = "1"
	} else {
		contactData = "0"
	}

	command := "WCS" + contactCode + fmt.Sprintf("%04d", contactNo) + contactData
	return format(dstAddress, command)
}

// 指定したワード単位での接点状態を読み出す電文を生成して返す
// dataCode:
//    D: データレジスタ
//    L: リンクレジスタ
//    F: ファイルレジスタ
func formatReadDataArea(dstAddress uint, dataCode string, startWordNo uint, endWordNo uint) string {
	if !isValidCode(dataCode, []string{"D", "L", "F"}) {
		panic(fmt.Sprintln("invalid code:", dataCode))
	}

	command := "RD" + dataCode + fmt.Sprintf("%05d", startWordNo) + fmt.Sprintf("%05d", endWordNo)
	return format(dstAddress, command)
}

// データエリアへ、配列の引数として渡された値を書き込む
// dataCode:
//    D: データレジスタ
//    L: リンクレジスタ
//    F: ファイルレジスタ
func formatWriteDataArea(dstAddress uint, dataCode string, startWordNo uint, values []uint32) string {
	if !isValidCode(dataCode, []string{"D", "L", "F"}) {
		panic(fmt.Sprintln("invalid code:", dataCode))
	}

	valuesBin := []byte{}
	index := 0
	for i, val := range values {
		hex := []byte(fmt.Sprintf("%04x", val))
		upper := hex[:2]
		lower := hex[2:]
		valBin := append(lower, upper...)
		valuesBin = append(valuesBin, valBin...)
		index = i
	}

	endWordNo := startWordNo + uint(index)
	command := "WD" + dataCode + fmt.Sprintf("%05d", startWordNo) + fmt.Sprintf("%05d", endWordNo)
	command += strings.ToUpper(string(valuesBin))
	return format(dstAddress, command)
}

// 渡されたコマンド本体にヘッダとBCCを追加する。
// ここでは末尾のCRは付加しません。送信時に付加してください
func format(dstAddress uint, body string) string {
	command := header() + address(dstAddress) + command() + body
	sendData := command + getBcc(command)
	return sendData
}

func header() string {
	return "%"
}

// dstAddress 1 - 32, FF(255) is Global send.
func address(ad uint) string {
	if (ad < 1 || 32 < ad) || ad == 255 {
		panic(fmt.Sprintf("Invalid mewtocol address: %d", ad))
	}
	return fmt.Sprintf("%02d", ad)
}

// 指定可能なコードかどうかチェックし、OKならtrueを返す。
func isValidCode(code string, list []string) bool {
	for _, s := range list {
		if code == s {
			return true
		}
	}

	log.Fatal(fmt.Sprintf("invalid code:", code))
	return false
}

func command() string {
	return "#"
}
