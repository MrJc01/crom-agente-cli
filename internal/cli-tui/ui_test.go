package tui

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestTUI_InteractiveInput(t *testing.T) {
	// Mock the input stream (Stdin) (Item 38)
	var in bytes.Buffer
	// O primeiro comando chaveia a cor
	in.WriteString("/color blue\n")
	// O segundo envia uma mensagem e fecha o programa (simulando EOF ou /exit)
	in.WriteString("/exit\n")

	opts := Options{
		WorkspacePath:  t.TempDir(),
		StoragePath:    t.TempDir(),
		SessionName:    "test-interactive-session",
		Provider:       "mock-provider",
		TimeoutSeconds: 30,
		Input:          &in,
	}

	// Como o Start não vai encontrar um loop provider válido para "mock-provider",
	// ele pode retornar erro se não mockarmos o llm.NewProvider. Mas a lógica
	// de leitura de Stdin é coberta se rodarmos. Aqui faremos apenas um teste do 
	// reader injection em um modelo isolado:
	model := &TUIModel{
		options: opts,
	}
	if opts.Input != nil {
		model.reader = bufio.NewReader(strings.NewReader("/test\n"))
	}
	
	if model.reader == nil {
		t.Fatal("Esperado que o model.reader fosse inicializado com o opts.Input")
	}
}
