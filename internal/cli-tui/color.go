package tui

import (
	"os"
	"strings"
)

// IsColorSupported verifica se o terminal atual suporta cores ANSI (Item 39)
func IsColorSupported() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	term := os.Getenv("TERM")
	if term == "dumb" || term == "" {
		return false
	}
	return true
}

// Colorize aplica o código ANSI se houver suporte a cores, caso contrário retorna o texto puro
func Colorize(ansiCode, text string) string {
	if !IsColorSupported() {
		return text
	}
	return ansiCode + text + "\033[0m"
}

// ColorMap exporta os códigos comuns para facilitar
var ColorMap = map[string]string{
	"red":    "\033[1;31m",
	"green":  "\033[1;32m",
	"yellow": "\033[1;33m",
	"blue":   "\033[1;34m",
	"purple": "\033[1;35m",
	"cyan":   "\033[1;36m",
	"orange": "\033[38;5;208m",
	"pink":   "\033[38;5;205m",
}
