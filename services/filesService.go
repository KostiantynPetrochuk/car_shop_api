package filesService

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/chai2010/webp"
)

func ConvertToWebp(inputPath, outputPath string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var buf bytes.Buffer
	options := &webp.Options{Lossless: true}
	if err := webp.Encode(&buf, img, options); err != nil {
		fmt.Println(err)
		return err
	}

	if err := os.WriteFile(outputPath, buf.Bytes(), 0666); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func ConvertHeicToWebp(inputPath, outputPath string) error {
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		fmt.Println("error creating directory: ", err)
		return err
	}
	absInputPath, err := filepath.Abs(inputPath)
	if err != nil {
		fmt.Println("error getting absolute path of input file: ", err)
		return err
	}
	absOutputPath, err := filepath.Abs(outputPath)
	if err != nil {
		fmt.Println("error getting absolute path of output file: ", err)
		return err
	}
	cmd := exec.Command("convert", absInputPath, absOutputPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		return err
	}
	if _, err := os.Stat(absOutputPath); os.IsNotExist(err) {
		fmt.Println("error: output file not created")
		return err
	}
	return nil
}
