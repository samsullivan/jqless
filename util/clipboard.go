package util

import "golang.design/x/clipboard"

var isClipboardInitialized bool

func InitClipboard() error {
	if isClipboardInitialized {
		return nil
	}

	return clipboard.Init()
}

func WriteClipboard(data []byte) error {
	if err := InitClipboard(); err != nil {
		return err
	}

	clipboard.Write(clipboard.FmtText, data)
	return nil
}
