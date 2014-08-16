package main

import (
	"fmt"
	"github.com/hiroeorz/mewtocol-go/mewtocol"
	"log"
	"os"
)

/* 出力例
$ ./mewtocol-go

---------------------------------------------
output0-> true
output1-> false
output2-> true
---------------------------------------------
input0-> false
input1-> true
input2-> false
---------------------------------------------
output list-> [5 0 0]
---------------------------------------------
input list-> [2 0 0]
---------------------------------------------
dataArea-> [0 1 2 3 4 55 665 776 906 9999]
---------------------------------------------
*/

func main() {
	setupCommand := `rdb_set confv250.enable 0 && 
                         /usr/sbin/sys -r rs232 &&
                         stty -F /dev/ttyAPP4 9600 parenb parodd cs8 -cstopb cread -crtscts`

	// Above setup command is my device special.
	// Normally use
	// ```setupCommand = "stty -F /dev/ttyAPP4 9600 parenb parodd cs8 -cstopb cread -crtscts"```
	f, err := mewtocol.OpenPLC("/dev/ttyAPP4", setupCommand)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var dst uint = 1 //宛先アドレス。通常は1固定

	// 外部出力を 0:オン 1:オフ 2:オン に設定
	_, err = mewtocol.WriteIOSingle(f, dst, "Y", 0, true)
	_, err = mewtocol.WriteIOSingle(f, dst, "Y", 1, false)
	_, err = mewtocol.WriteIOSingle(f, dst, "Y", 2, true)

	// 外部出力のステータスをひとつづつ取得する
	// 出力接点のステータスはON:true, OFF:false として取得します。
	putLine()
	var state bool

	state, err = mewtocol.ReadIOSingle(f, dst, "Y", 0)
	fmt.Println("output0->", state)

	state, err = mewtocol.ReadIOSingle(f, dst, "Y", 1)
	fmt.Println("output1->", state)

	state, err = mewtocol.ReadIOSingle(f, dst, "Y", 2)
	fmt.Println("output2->", state)

	putLine()
	// 外部入力のステータスをひとつづつ取得する
	// 入力接点のステータスはON:true, OFF:false として取得します。
	state, err = mewtocol.ReadIOSingle(f, dst, "X", 0)
	fmt.Println("input0->", state)

	state, err = mewtocol.ReadIOSingle(f, dst, "X", 1)
	fmt.Println("input1->", state)

	state, err = mewtocol.ReadIOSingle(f, dst, "X", 2)
	fmt.Println("input2->", state)
	putLine()

	// 外部出力のステータスをワード単位で取得する。
	// 外部出力を 0:オン 1:オフ 2:オン、残り全てオフの場合、0番目の値は2進数で```00000101```となり、
	// 結果として```[5, 0, 0]```を得ます。
	var stateWords []uint32

	stateWords, err = mewtocol.ReadIOWord(f, dst, "Y", 0, 2)
	fmt.Println("output list->", stateWords)

	putLine()
	// 外部入力のステータスをワード単位で取得する。
	// 外部入力を 0:オフ 1:オン 2:オフ、残り全てオフの場合、0番目の値は2進数で```00000010```となり、
	// 結果として```[2, 0, 0]```を得ます。
	stateWords, err = mewtocol.ReadIOWord(f, dst, "X", 0, 2)
	fmt.Println("input list->", stateWords)

	putLine()

	// データエリアに値を書き込む
	areaDataList := []uint32{0, 1, 2, 3, 4, 55, 665, 776, 906, 9999}
	_, err = mewtocol.WriteDataArea(f, dst, "D", 0, areaDataList)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// データエリアから値を読み込む
	areaDataListRead, err := mewtocol.ReadDataArea(f, dst, "D", 0, 9)
	fmt.Println("dataArea->", areaDataListRead)

	putLine()
}

func putLine() {
	fmt.Println("---------------------------------------------")
}
