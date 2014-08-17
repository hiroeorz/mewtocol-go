package mewtocol

import (
	"errors"
	"fmt"
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
func ReadIOSingle(plc *Serial, dst uint, contactCode string, contactNo uint) (bool, error) {
	fmt.Println("channel!")
	sendStr := formatReadIOSingle(dst, contactCode, contactNo)
	plc.ReqCh <- sendStr
	recvStr := <-plc.ResCh

	if recvStr == "" {
		return false, errors.New("failed to communicate with PLC")
	} else {
		return parseReadIOSingle(recvStr)
	}
}

// 指定したワード単位での接点状態を読み出す
func ReadIOWord(plc *Serial, dst uint, contactCode string, startNo uint, endNo uint) ([]uint, error) {
	sendStr := formatReadIOWord(dst, contactCode, startNo, endNo)
	plc.ReqCh <- sendStr
	recvStr := <-plc.ResCh

	if recvStr == "" {
		return nil, errors.New("failed to communicate with PLC")
	} else {
		return parseReadIOWord(recvStr)
	}
}

// 指定した1点の接点出力を書き込む
// 出力をONにするときはstateにtrueを、OFFにするときはfalseを指定してください
func WriteIOSingle(plc *Serial, dst uint, contactCode string, contactNo uint, state bool) (bool, error) {
	sendStr := formatWriteIOSingle(dst, contactCode, contactNo, state)
	plc.ReqCh <- sendStr
	recvStr := <-plc.ResCh

	if recvStr == "" {
		return false, errors.New("failed to communicate with PLC")
	} else {
		return parseWriteIOSingle(recvStr)
	}
}

// 指定したワード単位でのデータエリア状態を読み出す
func ReadDataArea(plc *Serial, dst uint, dataCode string, startNo uint, endNo uint) ([]uint, error) {
	sendStr := formatReadDataArea(dst, dataCode, startNo, endNo)
	plc.ReqCh <- sendStr
	recvStr := <-plc.ResCh

	if recvStr == "" {
		return nil, errors.New("failed to communicate with PLC")
	} else {
		return parseReadDataArea(recvStr)
	}
}

// 配列の引数として渡された値をデータエリアへ書き込む。
func WriteDataArea(plc *Serial, dst uint, dataCode string, startNo uint, values []uint32) (bool, error) {
	sendStr := formatWriteDataArea(dst, dataCode, startNo, values)
	plc.ReqCh <- sendStr
	recvStr := <-plc.ResCh

	if recvStr == "" {
		return false, errors.New("failed to communicate with PLC")
	} else {
		return parseWriteDataArea(recvStr)
	}
}

// 水平パリティチェックの為のコード(2byte)を生成して返す。
// ヘッダからBCC直前までのデータにたいして1byteづつ排他的論理和をとった結果を16進数文字列にして返す。
func getBcc(str string) string {
	buff := []byte(str)
	result := buff[0]

	for i := 1; i < len(buff); i++ {
		result = result ^ (buff[i])
	}
	return strings.ToUpper(fmt.Sprintf("%02x", result))
}
