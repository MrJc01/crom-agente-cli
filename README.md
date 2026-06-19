# crom-agente-cli

Interface de terminal interativa (REPL/TUI) baseada em chat para co-programar com o `crom-agente`. Estilizado e inspirado em interfaces como Claude Code e Gemini CLI.

Desenvolvido em Go utilizando a suíte de TUI da Charmbracelet (`bubbletea`, `lipgloss`, `bubbles`, `glamour`).

---

## Como Compilar

> [!IMPORTANT]
> **Restrição de Pacotes `internal` do Go:**
> O CLI importa pacotes internos do projeto principal (`github.com/crom/crom-agente/internal/...`). De acordo com as regras do Go, pacotes `internal` não podem ser importados por módulos externos (como o `github.com/crom/crom-agente-cli` que possui seu próprio `go.mod`).
> 
> Por isso, **a compilação e execução de testes locais devem sempre ser realizadas de dentro do repositório principal** (`crom-agente`), onde o compilador reconhece todos os pacotes como parte do mesmo módulo.

Para compilar o binário de forma correta e atualizada:

```bash
# Navegue até o repositório principal
cd ../crom-agente

# Compile o CLI em modo headless (evitando dependências do systray)
go build -tags headless -o bin/crom-agente-cli ./cmd/crom-agente-cli

# Copie o binário atualizado para esta pasta
cp bin/crom-agente-cli ../crom-agente-cli/
```

Para rodar os testes unitários e de integração do CLI de dentro do repositório principal:
```bash
go test -v -tags headless ./internal/cli-tui/...
go test -v -tags headless ./internal/cli/...
```

---

## Como Executar

Para iniciar a interface interativa em seu terminal, basta rodar o binário:

```bash
./crom-agente-cli
```

### Flags Disponíveis:
- `-w, --workspace`: Caminho para o workspace do projeto (padrão `.`)
- `-s, --storage`: Diretório de persistência do estado do agente (padrão `.crom`)
- `--session`: Nome ou ID da sessão de chat (padrão `cli-session`)
- `--provider`: Override de provedor de LLM (`openai`, `gemini`, `anthropic`, `ollama`, `openrouter`)
- `--model`: Override do modelo de LLM
- `--permission-mode`: Modo de segurança para execução de ferramentas (`total_access`, `ask_every_time`, `scoped`)

Exemplo especificando o OpenRouter e Gemini na sessão `chat-otimizacao`:
```bash
./crom-agente-cli --session chat-otimizacao --provider openrouter --model google/gemini-2.5-flash
```

---

## Slash Commands no Chat
Durante a conversa com o agente, você pode digitar comandos especiais precedidos por barra (`/`):
- `/add <caminho>`: Anexa o conteúdo de um arquivo local diretamente no próximo prompt que você enviar ao agente.
- `/session <nome>`: Chaveia dinamicamente para outra sessão salvando/carregando o histórico.
- `/clear`: Limpa o histórico de chat exibido na tela.
- `/help`: Mostra a ajuda no chat.
- `/exit` ou `/quit`: Encerra o aplicativo.

---

## 📄 Licença

Este projeto é público e licenciado sob a **Licença Pública de Conteúdo CromIA**.

Copyright (C) 2026-presente CromIA, Todos os direitos reservados.

Esta licença concede o direito limitado de uso, modificação e redistribuição do código.
- **Uso Não Comercial:** Totalmente permitido.
- **Uso Comercial:** Permitido desde que o ganho monetário (direto ou indireto) utilizando o conteúdo seja menor que 1 milhão de BRL (consulte a licença para valores em dólar).

Para ver todas as regras, limites monetários e condições de conversão para a licença MIT, consulte o arquivo [LICENSE](LICENSE) na raiz deste repositório.

---

🌐 **Site Oficial:** [ia.crom.run](https://ia.crom.run)
