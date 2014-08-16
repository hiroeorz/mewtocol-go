package mewtocol

import (
	"errors"
	"fmt"
	"strconv"
)

type Mewtocol struct {
	Header  string
	Address uint16
	Code    string
	Body    []byte
}

// PLCからの接点単体取得応答をパースして値を返す
func parseReadIOSingle(str string) (bool, error) {
	res, err := parseHeader(str)
	if err != nil {
		return false, err
	}

	contactData := string(res.Body)

	if contactData == "1" {
		return true, nil
	} else {
		return false, nil
	}
}

// PLCからの接点ワード単位取得応答をパースして値を返す
func parseReadIOWord(str string) ([]uint32, error) {
	res, err := parseHeader(str)
	if err != nil {
		return nil, err
	}

	return parseListData(res.Body), nil
}

// PLCからの接点出力単体書き込み応答をパースして返します
func parseWriteIOSingle(str string) (bool, error) {
	_, err := parseHeader(str)
	if err != nil {
		return false, err
	}

	return true, nil
}

// PLCからのデータエリア取得応答をパースして値を返す
func parseReadDataArea(str string) ([]uint32, error) {
	res, err := parseHeader(str)
	if err != nil {
		return nil, err
	}

	return parseListData(res.Body), nil
}

// PLCからの接点出力単体書き込み応答をパースして返します
func parseWriteDataArea(str string) (bool, error) {
	_, err := parseHeader(str)
	if err != nil {
		return false, err
	}

	return true, nil
}

// byteのスライスを2バイトづつ数値に変換して、数値のスライスにして返す。
func parseListData(data []byte) []uint32 {
	count := len(data) / 4
	list := []uint32{}

	for i := 0; i < count; i++ {
		n := i * 4
		valLower := data[(n + 0):(n + 2)] // 下位を前半に受信する
		valUpper := data[(n + 2):(n + 4)] // 上位を後半で受信する
		val := []byte{}
		val = append(val, valUpper...)
		val = append(val, valLower...)
		intVal, _ := strconv.ParseUint(string(val), 16, 32)
		list = append(list, uint32(intVal))
	}

	return list
}

func parseHeader(str string) (*Mewtocol, error) {
	buff := []byte(str)

	success := string(buff[3])
	if success == "$" {
		header := string(buff[0])
		address, _ := strconv.ParseInt(string(buff[1:3]), 10, 16)
		code := string(buff[4:6])
		body := getReqBody(buff)
		return &Mewtocol{header, uint16(address), code, body}, nil
	} else if success == "!" {
		errNo, _ := strconv.ParseInt(string(buff[4:6]), 16, 16)
		return nil, errors.New(fmt.Sprintf("mewtocol error response:%x", errNo))
	} else {
		panic(fmt.Sprintf("invalid success code:", success))
	}

}

func getReqBody(buff []byte) []byte {
	headerSize := 6
	footerSize := 2
	bodySize := len(buff) - (headerSize + footerSize)
	bodyStart := headerSize
	bodyEnd := headerSize + bodySize
	return buff[bodyStart:bodyEnd]
}
