package mewtocol

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// 接点コード
// "X"  外部入力
// "Y"  外部出力
// "R"  内部リレー
// "L"  リンクリレー
// "T"  タイマ
// "C"  カウンタ

// 指定した1点の接点状態を読み出す
func ReadIOSingle(f *os.File, dst uint, contactCode string, contactNo uint) (bool, error) {
	sendStr := formatReadIOSingle(dst, contactCode, contactNo)
	recvStr, err := send(f, sendStr)
	if err != nil {
		return false, err
	} else {
		return ParseReadIOSingle(recvStr)
	}
}

// 指定したワード単位での接点状態を読み出す
func ReadIOWord(f *os.File, dst uint, contactCode string, startNo uint, endNo uint) ([]uint32, error) {
	sendStr := formatReadIOWord(dst, contactCode, startNo, endNo)
	recvStr, err := send(f, sendStr)
	if err != nil {
		return []uint32{}, err
	} else {
		return ParseReadIOWord(recvStr)
	}
}

// 指定した1点の接点出力を書き込む
// 出力をONにするときはstateにtrueを、OFFにするときはfalseを指定してください
func WriteIOSingle(f *os.File, dst uint, contactCode string, contactNo uint, state bool) (bool, error) {
	sendStr := formatWriteIOSingle(dst, contactCode, contactNo, state)
	recvStr, err := send(f, sendStr)
	if err != nil {
		return false, err
	} else {
		return ParseWriteIOSingle(recvStr)
	}
}

// 指定したワード単位でのデータエリア状態を読み出す
func ReadDataArea(f *os.File, dst uint, dataCode string, startNo uint, endNo uint) ([]uint32, error) {
	sendStr := formatReadDataArea(dst, dataCode, startNo, endNo)
	recvStr, err := send(f, sendStr)
	if err != nil {
		return []uint32{}, err
	} else {
		return ParseReadDataArea(recvStr)
	}
}

// 配列の引数として渡された値をデータエリアへ書き込む。
func WriteDataArea(f *os.File, dst uint, dataCode string, startNo uint, values []uint32) (bool, error) {
	sendStr := formatWriteDataArea(dst, dataCode, startNo, values)
	recvStr, err := send(f, sendStr)
	if err != nil {
		return false, err
	} else {
		return ParseWriteDataArea(recvStr)
	}
}

// 与えられた文字列を送信し、指定サイズのレスポンスを受信して返す。
func send(f *os.File, sendStr string) (string, error) {
	_, err := Write(f, sendStr)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	recvStr, err := Read(f)
	if err != nil {
		log.Fatal(err)
		return "", err
	} else {
		return recvStr, nil
	}
}

func getBcc(str string) string {
	buff := []byte(str)
	result := buff[0]

	for i := 1; i < len(buff); i++ {
		result = result ^ (buff[i])
	}
	return strings.ToUpper(fmt.Sprintf("%02x", result))
}
