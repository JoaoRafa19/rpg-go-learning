# RPG Go!

Um jogo RPG simples desenvolvido em Go utilizando a biblioteca Ebiten.

## Índice

- [Introdução](#introdução)
- [Instalação](#instalação)
- [Como Jogar](#como-jogar)
- [Controles](#controles)
- [Estrutura do Projeto](#estrutura-do-projeto)
- [Dependências](#dependências)
- [Contribuição](#contribuição)
- [Licença](#licença)

## Introdução

**RPG Go!** é um jogo de RPG 2D desenvolvido em Go, usando a biblioteca Ebiten para renderização gráfica. O jogo apresenta um personagem principal que explora um mapa, interage com inimigos e coleta poções. O mapa é criado utilizando o editor de tilemaps [Tiled](https://www.mapeditor.org/).

## Instalação

### Pré-requisitos

- Go (versão 1.16 ou superior)
- Git

### Clonando o Repositório

```bash
git clone https://github.com/seu-usuario/rpg-go.git
cd rpg-go
```

### Instalando Dependências

Utilize o comando `go get` para instalar as dependências necessárias:

```bash
go get -u github.com/hajimehoshi/ebiten/v2
```

## Como Jogar

Para executar o jogo, utilize o comando:

```bash
go run main.go
```

Uma janela será aberta exibindo o jogo. Controle o personagem principal e explore o mapa!

## Controles

- **W**: Move o jogador para cima
- **A**: Move o jogador para a esquerda
- **S**: Move o jogador para baixo
- **D**: Move o jogador para a direita

## Estrutura do Projeto

- `main.go`: Arquivo principal que inicia o jogo e contém a lógica de atualização e renderização.
- `entities/`: Pacote que contém as definições das entidades do jogo como Player, Enemy e Potion.
  - `player.go`: Define o jogador e suas animações.
  - `enemy.go`: Define os inimigos e seu comportamento.
  - `potion.go`: Define as poções e suas propriedades.
- `animations/`: Contém o sistema de animações utilizado pelo jogador.
  - `animation.go`: Gerencia quadros de animação e atualização de estado.
- `spritesheet/`: Pacote para manipulação de spritesheets.
  - `spritesheet.go`: Define a estrutura e métodos para extrair sprites individuais.
- `assets/`: Diretório que contém todos os recursos do jogo.
  - `images/`: Contém as imagens utilizadas (personagens, tilesets, etc.).
  - `maps/`: Contém os arquivos de mapa gerados pelo Tiled.
- `tilemap.go`: Carrega e processa os mapas criados no Tiled.
- `camera.go`: Implementa a lógica da câmera que segue o jogador.

## Dependências

- [Ebiten](https://github.com/hajimehoshi/ebiten): Biblioteca de jogos 2D em Go.
- [Tiled](https://www.mapeditor.org/): Editor utilizado para criar os mapas do jogo.

## Contribuição

Contribuições são bem-vindas! Se você quiser melhorar o jogo ou adicionar novas funcionalidades, sinta-se à vontade para fazer um fork do projeto e enviar um pull request.

### Passos para Contribuir

1. Faça um fork do repositório.
2. Crie uma nova branch com a sua feature ou correção de bug: `git checkout -b minha-feature`.
3. Commit suas alterações: `git commit -m 'Adiciona nova funcionalidade'`.
4. Faça push para a branch: `git push origin minha-feature`.
5. Abra um pull request.

## Licença

Este projeto está licenciado sob a licença MIT. Consulte o arquivo [LICENSE](LICENSE) para mais detalhes.

---

**Nota:** Para quaisquer dúvidas ou problemas, sinta-se livre para abrir uma issue no repositório.