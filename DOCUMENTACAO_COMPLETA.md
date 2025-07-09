# Documentação Completa do Projeto RPG Go!

## Sumário
- [Visão Geral](#visão-geral)
- [Arquitetura e Estrutura de Pastas](#arquitetura-e-estrutura-de-pastas)
- [Descrição dos Arquivos e Pacotes](#descrição-dos-arquivos-e-pacotes)
- [Fluxo de Execução](#fluxo-de-execução)
- [Detalhamento dos Componentes](#detalhamento-dos-componentes)
- [Recursos e Assets](#recursos-e-assets)
- [Dependências](#dependências)
- [Como Executar](#como-executar)
- [Como Contribuir](#como-contribuir)
- [Licença](#licença)

---

## Visão Geral

**RPG Go!** é um jogo 2D de RPG desenvolvido em Go, utilizando a biblioteca Ebiten para renderização gráfica. O jogador controla um personagem que explora mapas, enfrenta inimigos e coleta poções. Os mapas são criados com o editor Tiled.

## Arquitetura e Estrutura de Pastas

```
├── main.go                # Ponto de entrada do jogo
├── game.go                # Gerenciamento do ciclo de vida do jogo e cenas
├── assets/                # Recursos do jogo (imagens, mapas)
│   ├── images/            # Sprites, tilesets, personagens
│   └── maps/              # Mapas e tilesets do Tiled
├── animations/            # Sistema de animações
│   └── animation.go
├── camera/                # Lógica da câmera
│   └── camera.go
├── components/            # Componentes de entidades (ex: combate)
│   └── combat.go
├── constants/             # Constantes globais
│   └── constants.go
├── entities/              # Entidades do jogo (player, inimigos, poções, sprites)
│   ├── player.go
│   ├── enemy.go
│   ├── potion.go
│   └── sprite.go
├── scenes/                # Gerenciamento de cenas (jogo, pausa, início)
│   ├── gamescenes.go
│   ├── pausescene.go
│   ├── scenes.go
│   └── startscene.go
├── spritesheet/           # Manipulação de spritesheets
│   └── spritesheet.go
├── tilemap/               # Carregamento e processamento de mapas
│   └── tilemap.go
├── tileset/               # Manipulação de tilesets
│   └── tileset.go
```

## Descrição dos Arquivos e Pacotes

### Arquivos Principais
- **main.go**: Ponto de entrada. Inicializa a janela, configura o Ebiten e executa o loop principal do jogo.
- **game.go**: Gerencia o ciclo de vida do jogo, alternando entre cenas (início, jogo, pausa) e delegando atualização/desenho.

### Pacote `entities`
- **player.go**: Estrutura e lógica do jogador, incluindo animações e combate.
- **enemy.go**: Estrutura dos inimigos, incluindo IA simples e combate.
- **potion.go**: Estrutura das poções, com quantidade de cura.
- **sprite.go**: Estrutura base para entidades visuais, com posição e imagem.

### Pacote `animations`
- **animation.go**: Estrutura e lógica para animações de sprites, controle de frames e looping.

### Pacote `components`
- **combat.go**: Interfaces e implementações de combate para entidades (vida, ataque, dano, cooldowns).

### Pacote `constants`
- **constants.go**: Constantes globais, como tamanho dos tiles.

### Pacote `camera`
- **camera.go**: Lógica da câmera que segue o jogador e limita a visualização ao mapa.

### Pacote `spritesheet`
- **spritesheet.go**: Manipulação de spritesheets, cálculo de retângulos para extração de sprites.

### Pacote `tilemap`
- **tilemap.go**: Carregamento de mapas do Tiled (JSON), parsing de camadas e tilesets.

### Pacote `tileset`
- **tileset.go**: Manipulação de tilesets, extração de imagens de tiles individuais.

### Pacote `scenes`
- **gamescenes.go**: Implementação da cena principal do jogo, incluindo lógica de atualização, desenho, carregamento de entidades, colisão, etc.
- **pausescene.go**: Cena de pausa, exibe overlay e opções de continuar ou sair.
- **startscene.go**: Cena inicial, exibe tela de início e aguarda input para começar.
- **scenes.go**: Interface e enumeração de cenas, gerenciamento de transições.

## Fluxo de Execução

1. **main.go** inicializa a janela e executa o loop principal via Ebiten.
2. **game.go** gerencia a cena ativa (início, jogo, pausa) e delega as funções de update/draw.
3. Cada cena implementa a interface `Scene` e é responsável por seu próprio ciclo de vida.
4. O jogador e inimigos são instanciados a partir de `entities/`, com animações e componentes de combate.
5. O mapa é carregado de arquivos JSON do Tiled, processado por `tilemap/` e `tileset/`.
6. A câmera segue o jogador e limita a visualização ao tamanho do mapa.
7. Recursos gráficos são carregados de `assets/images` e mapas de `assets/maps`.

## Detalhamento dos Componentes

### Entidades
- **Player**: Possui animações para cada direção, componente de combate, posição e sprite.
- **Enemy**: IA simples (segue jogador), componente de combate, sprite.
- **Potion**: Item coletável, restaura vida do jogador.
- **Sprite**: Estrutura base para renderização de entidades.

### Animações
- **Animation**: Controla frames, velocidade, looping e atualização de animações de sprites.

### Combate
- **BasicCombat**: Implementa vida, ataque, estado de ataque e dano para entidades.
- **EnemyCombat**: Extende BasicCombat com cooldown de ataque.

### Câmera
- Segue o jogador, centraliza na tela e limita para não mostrar áreas fora do mapa.

### Spritesheet
- Calcula retângulos para extrair sprites individuais de uma imagem maior.

### Tilemap e Tileset
- Carregam mapas do Tiled, processam camadas, tilesets e extraem imagens dos tiles.

### Scenes
- **GameScene**: Cena principal, gerencia entidades, colisão, lógica do jogo.
- **PauseScene**: Overlay de pausa, opções de continuar ou sair.
- **StartScene**: Tela inicial, aguarda input para iniciar o jogo.

## Recursos e Assets
- **assets/images/**: Sprites de personagens, inimigos, tilesets, etc.
- **assets/maps/**: Mapas criados no Tiled, tilesets e arquivos auxiliares.

## Dependências
- [Ebiten](https://github.com/hajimehoshi/ebiten): Biblioteca para jogos 2D em Go.
- [Tiled](https://www.mapeditor.org/): Editor de mapas tile-based.

## Como Executar

1. Instale o Go (>=1.16).
2. Instale as dependências:
   ```
   go get -u github.com/hajimehoshi/ebiten/v2
   ```
3. Execute o jogo:
   ```
   go run main.go
   ```

## Como Contribuir

1. Faça um fork do repositório.
2. Crie uma branch para sua feature/correção.
3. Faça commits e push.
4. Abra um Pull Request.

## Licença

MIT. Veja o arquivo LICENSE para detalhes.

---

**Observação:** Para detalhes de cada função, structs e métodos, consulte os comentários no código-fonte de cada arquivo.
