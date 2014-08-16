package main

import (
	"fmt"
	"github.com/hiroeorz/mewtocol-go/mewtocol"
	"log"
)

func main() {
	setupCommand := `rdb_set confv250.enable 0 && 
                         /usr/sbin/sys -r rs232 &&
                         stty -F /dev/ttyAPP4 9600 parenb parodd cs8 -cstopb cread -crtscts`

	// Above setup command is my device special.
	// Normaly use 
	// ```setupCommand = "stty -F /dev/ttyAPP4 9600 parenb parodd cs8 -cstopb cread -crtscts"```
	f, err := mewtocol.OpenPLC("/dev/ttyAPP4", setupCommand)

	if err != nil {
		log.Fatal(err)
	}

	var dst uint = 1

	// 外部出力を 0:オン 1:オフ 2:オン に設定
	mewtocol.WriteIOSingle(f, dst, "Y", 0, true)
	mewtocol.WriteIOSingle(f, dst, "Y", 1, false)
	mewtocol.WriteIOSingle(f, dst, "Y", 2, true)

	// 外部出力のステータスをひとつづつ取得する
	fmt.Println("------------")
	res10, _ := mewtocol.ReadIOSingle(f, dst, "Y", 0)
	fmt.Println("input0->", res10)

	res11, _ := mewtocol.ReadIOSingle(f, dst, "Y", 1)
	fmt.Println("input1->", res11)

	res12, _ := mewtocol.ReadIOSingle(f, dst, "Y", 2)
	fmt.Println("input2->", res12)
	fmt.Println("------------")

	// 外部出力のステータスをワード単位で取得する
	res20, _ := mewtocol.ReadIOWord(f, dst, "Y", 0, 2)
	fmt.Println("input list->", res20)

	fmt.Println("=================================================")

	// データエリアに値を書き込む
	areaVals := []uint32{0, 1, 2, 3, 4, 55, 665, 776, 906, 9999}
	_, err = mewtocol.WriteDataArea(f, dst, "D", 0, areaVals)
	if err != nil {
		log.Fatal(err)
	}

	// データエリアから値を読み込む
	resDataArea, _ := mewtocol.ReadDataArea(f, dst, "D", 0, 9)
	fmt.Println("dataArea->", resDataArea)
	fmt.Println("------------")
}
