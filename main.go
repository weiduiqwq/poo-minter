package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/genv"
	"log"
	"os"
	"poo/poo_minter"
	"sync"
)

func main() {
	ctx := context.Background()

	// 从INIT_DATA指定的路径获取数据
	filePath := genv.GetWithCmd("INIT_DATA").String()

	// 检查文件是否存在
	fileInfo, err := os.Stat(filePath)
	if err == nil && !fileInfo.IsDir() {
		// 文件存在且不是目录，按照文件方式处理

		file, err := os.Open(filePath)
		if err != nil {
			log.Fatalf("无法打开文件：%v", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		// 使用 WaitGroup 等待所有线程完成
		var wg sync.WaitGroup

		for scanner.Scan() {
			minterData := scanner.Text()
			fmt.Println(minterData)

			// 为每个minter启动一个新线程

			wg.Add(1)
			go func(data string) {
				defer wg.Done()

				// 创建独立的minter实例
				minter := poo_minter.NewPooMinter(data)

				// 运行minter.Mint(ctx)
				err := minter.Mint(ctx)
				if err != nil {
					log.Printf("minter.Mint 错误：%+v", err)
				}
			}(minterData)
		}

		wg.Wait()

		if err := scanner.Err(); err != nil {
			log.Fatalf("读取文件错误：%v", err)
		}

	} else {
		minter := poo_minter.NewPooMinter(filePath)
		err := minter.Mint(ctx)
		if err != nil {
			log.Fatalf("%+v", err)
		}
	}
}
