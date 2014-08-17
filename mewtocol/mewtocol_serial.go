package mewtocol

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

const (
	MAX_FLAME_SIZE int = 118
)

/*
Test by console

First, you have to setting serial device (ex:/dev/ttyAPP4)

```
$ stty -F /dev/ttyAPP4 9600 parenb parodd cs8 -cstopb cread -crtscts
```

console1

```
$ cat /dev/ttyAPP4


```

and console2

```
$ echo -e "%01#RCSX00001D\r" > /dev/ttyAPP4
```

Now, you caught response in console1

```
%01$RC021
```

*/

type Serial struct {
	ReqCh chan string
	ResCh chan string
}

// Mewtocol通信用にシリアルポートをオープンします
// setupCommand: exp: "stty -F /dev/ttyAPP4 9600 parenb parodd cs8 -cstopb cread -crtscts"
func OpenPLC(name string, setupCommand string) (*Serial, error) {
	execCommand(setupCommand)
	f, err := os.OpenFile(name, syscall.O_RDWR|syscall.O_NOCTTY, 0666)
	if err != nil {
		return nil, err
	}

	reqCh := make(chan string, 0) // for request
	resCh := make(chan string, 0) // for response
	plc := &Serial{ReqCh: reqCh, ResCh: resCh}
	go deviceLoop(f, plc)

	return plc, nil
}

// ファイルに直接書き込み/読み込みするgoroutine
func deviceLoop(f *os.File, plc *Serial) {
	defer f.Close()

	for {
		sendStr := <-plc.ReqCh
		recvStr, err := send(f, sendStr)

		if err != nil {
			plc.ResCh <- ""
		} else {
			plc.ResCh <- recvStr
		}
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

// シリアルポートにデータを書き込みます。
func Write(f *os.File, str string) (int, error) {
	sendData := []byte(str + "\r")
	n, err := f.Write(sendData)
	if err != nil {
		log.Fatal(fmt.Sprintf("serial write error: %s", str))
	}

	return n, err
}

// シリアルポートからデータを読み込みます。
func Read(f *os.File) (string, error) {
	buff := make([]byte, 0, MAX_FLAME_SIZE)

	for {
		byte := make([]byte, 1)
		_, err := f.Read(byte)
		if err != nil {
			log.Fatal("serial read error")
			return "", err
		}

		if string(byte) == "\r" {
			break
		}
		buff = append(buff, byte...)
	}

	if isValidBCC(buff) {
		return string(buff), nil
	} else {
		return "", errors.New(fmt.Sprintf("invalid BCC:", string(buff)))
	}

}

// OSコマンドを実行する
func execCommand(cmdStr string) {
	cmd := exec.Command("sh", "-c", cmdStr)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		fmt.Println(stderr.String())
	} else {
		fmt.Println(stdout.String())
	}
}

// データに付属していたBCCと計算したBCCが一致すればtrueをかえす。
func isValidBCC(buff []byte) bool {
	lengthBeforeBCC := len(buff) - 2
	command := string(buff[:lengthBeforeBCC])
	bcc := string(buff[lengthBeforeBCC:])
	okBcc := getBcc(command)
	return (bcc == okBcc)
}
