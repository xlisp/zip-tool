package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	// 默认分割大小：16MB
	DefaultCutSize = 16 * 1024 * 1024
)

// 分割ZIP文件为多个Base64编码的文本文件
func divisionZipFile(zipPath, outputDir, fileName, suffix string, cutSize int) error {
	// 打开源文件
	sourceFile, err := os.Open(zipPath)
	if err != nil {
		return fmt.Errorf("无法打开源文件: %v", err)
	}
	defer sourceFile.Close()

	// 获取文件信息
	fileInfo, err := sourceFile.Stat()
	if err != nil {
		return fmt.Errorf("无法获取文件信息: %v", err)
	}

	// 计算需要分割的文件数量
	fileSize := fileInfo.Size()
	numFiles := int(fileSize / int64(cutSize))
	if fileSize%int64(cutSize) != 0 {
		numFiles++
	}

	fmt.Printf("文件大小: %d 字节\n", fileSize)
	fmt.Printf("将分割为 %d 个文件\n", numFiles)

	// 确保输出目录存在
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("无法创建输出目录: %v", err)
	}

	// 创建缓冲读取器
	reader := bufio.NewReader(sourceFile)
	buffer := make([]byte, cutSize)

	for i := 0; i < numFiles; i++ {
		// 生成文件名，格式为：fileName + 三位数字 + suffix
		outputFileName := fmt.Sprintf("%s%03d%s", fileName, i, suffix)
		outputPath := filepath.Join(outputDir, outputFileName)

		// 读取数据块
		n, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("读取文件时出错: %v", err)
		}

		if n == 0 {
			break
		}

		// 对读取的数据进行Base64编码
		encodedData := base64.StdEncoding.EncodeToString(buffer[:n])

		// 写入到输出文件
		outputFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("无法创建输出文件 %s: %v", outputPath, err)
		}

		_, err = outputFile.WriteString(encodedData)
		outputFile.Close()

		if err != nil {
			return fmt.Errorf("写入文件 %s 时出错: %v", outputPath, err)
		}

		fmt.Printf("已创建文件: %s (大小: %d 字节)\n", outputFileName, n)
	}

	fmt.Println("文件分割完成！")
	return nil
}

// 合并Base64编码的文本文件为ZIP文件
func mergeToZipFile(inputDir, outputZipPath string) error {
	// 获取输入目录中的所有文件
	files, err := os.ReadDir(inputDir)
	if err != nil {
		return fmt.Errorf("无法读取输入目录: %v", err)
	}

	// 过滤并排序文件
	var textFiles []string
	for _, file := range files {
		if !file.IsDir() {
			textFiles = append(textFiles, file.Name())
		}
	}

	// 按文件名排序，确保正确的合并顺序
	sort.Strings(textFiles)

	if len(textFiles) == 0 {
		return fmt.Errorf("输入目录中没有找到文件")
	}

	fmt.Printf("找到 %d 个文件需要合并\n", len(textFiles))

	// 创建输出文件
	outputFile, err := os.Create(outputZipPath)
	if err != nil {
		return fmt.Errorf("无法创建输出文件: %v", err)
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	// 处理每个文件
	for _, fileName := range textFiles {
		filePath := filepath.Join(inputDir, fileName)
		fmt.Printf("正在处理文件: %s\n", fileName)

		// 读取文件内容
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("无法读取文件 %s: %v", fileName, err)
		}

		// Base64解码
		decodedData, err := base64.StdEncoding.DecodeString(string(fileContent))
		if err != nil {
			return fmt.Errorf("解码文件 %s 时出错: %v", fileName, err)
		}

		// 写入到输出文件
		_, err = writer.Write(decodedData)
		if err != nil {
			return fmt.Errorf("写入数据时出错: %v", err)
		}
	}

	fmt.Println("文件合并完成！")
	return nil
}

// 显示使用说明
func printUsage() {
	fmt.Println("ZIP文件分割与合并工具")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  分割模式: go run main.go split <zip文件路径> <输出目录> [文件名前缀] [文件扩展名] [分割大小(MB)]")
	fmt.Println("  合并模式: go run main.go merge <输入目录> <输出zip文件路径>")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  go run main.go split ./example.zip ./output test .txt 16")
	fmt.Println("  go run main.go merge ./output ./restored.zip")
	fmt.Println()
	fmt.Println("参数说明:")
	fmt.Println("  分割大小: 以MB为单位，默认为16MB")
	fmt.Println("  文件名前缀: 默认为'part'")
	fmt.Println("  文件扩展名: 默认为'.txt'")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "split":
		if len(os.Args) < 4 {
			fmt.Println("错误: 分割模式需要至少3个参数")
			printUsage()
			return
		}

		zipPath := os.Args[2]
		outputDir := os.Args[3]

		// 设置默认值
		fileName := "part"
		suffix := ".txt"
		cutSizeMB := 16

		// 解析可选参数
		if len(os.Args) > 4 {
			fileName = os.Args[4]
		}
		if len(os.Args) > 5 {
			suffix = os.Args[5]
		}
		if len(os.Args) > 6 {
			if size, err := strconv.Atoi(os.Args[6]); err == nil {
				cutSizeMB = size
			} else {
				fmt.Printf("警告: 无效的分割大小 '%s'，使用默认值16MB\n", os.Args[6])
			}
		}

		cutSize := cutSizeMB * 1024 * 1024

		fmt.Printf("开始分割文件...\n")
		fmt.Printf("源文件: %s\n", zipPath)
		fmt.Printf("输出目录: %s\n", outputDir)
		fmt.Printf("文件前缀: %s\n", fileName)
		fmt.Printf("文件扩展名: %s\n", suffix)
		fmt.Printf("分割大小: %dMB\n", cutSizeMB)

		if err := divisionZipFile(zipPath, outputDir, fileName, suffix, cutSize); err != nil {
			fmt.Printf("分割失败: %v\n", err)
			os.Exit(1)
		}

	case "merge":
		if len(os.Args) < 4 {
			fmt.Println("错误: 合并模式需要2个参数")
			printUsage()
			return
		}

		inputDir := os.Args[2]
		outputZipPath := os.Args[3]

		fmt.Printf("开始合并文件...\n")
		fmt.Printf("输入目录: %s\n", inputDir)
		fmt.Printf("输出文件: %s\n", outputZipPath)

		if err := mergeToZipFile(inputDir, outputZipPath); err != nil {
			fmt.Printf("合并失败: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("错误: 未知命令 '%s'\n", command)
		printUsage()
		os.Exit(1)
	}
}
