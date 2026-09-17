# go-tanks

A top-down 2D tank game written in Go with [Ebitengine](https://ebitengine.org/).
Drive a tank over hand-authored maps, fight AI enemies with ballistic shells,
and join online multiplayer sessions.

## Features

- **Tank controls** — move, rotate the hull and aim the turret independently.
- **Ballistic shooting** — shells follow a gravity arc; hold to charge the
  elevation and release to fire.
- **Enemy AI** — enemies seek, orbit and back away from the closest player,
  steer around obstacles, and fire with a ballistic aim.
- **Hand-authored maps** — levels are defined as editable JSON files (roads,
  terrain, obstacles and spawn points), not randomly generated.
- **Built-in map editor** — press `E` in game to paint tiles, place obstacles
  and spawn points, then save/play your map without recompiling.
- **Obstacles & collisions** — crates, barrels, trees and sandbags block tanks
  and shells (backed by [resolv](https://github.com/SolarLune/resolv)).
- **Score & HUD** — your kills vs. the computer's kills, crosshair while
  charging, high-speed fps overlay.
- **Online multiplayer** — a WebSocket hub server relays player state and
  simulates the shared AI enemies and score.

## Requirements

- [Go](https://go.dev/dl/) 1.24+
- A desktop environment with a windowing system to run the game
  (Ebitengine renders to a native window).

## Running (single player)

```bash
go run ./cmd
```

Optionally pick a specific map:

```bash
go run ./cmd -map district_01
```

## Running (multiplayer)

Start one server, then connect clients to it:

```bash
# terminal 1 — the hub server (also simulates enemies and score)
go run ./server -addr localhost:8080

# terminal 2+ — one per player
go run ./cmd -addr localhost:8080
```

The server chooses the map and owns the session; each connected player is
relayed to the others and the AI enemies are shared (simulated by the server).

## Controls

| Key(s)                  | Action                                        |
| ----------------------- | --------------------------------------------- |
| `←` / `→`               | Rotate the tank hull                          |
| `↑` / `↓`               | Move forward / backward                       |
| `A` / `D`               | Rotate the turret                             |
| `Space` (hold & release)| Charge elevation and fire the shell           |
| `P`                     | Regenerate a random map (single player only)  |

## Maps

Levels live in `internal/maps/maps/*.json` and use a compact character grid.
Each cell is a legend key mapped to a tile sprite, plus optional
`obstacles`, `playerSpawn` and `enemySpawns`.

Available maps: `city_01`, `city_02`, `district_01`, `desert_01`.

## Map editor

Press `E` while playing (single player) to open the built-in editor:

| Input                          | Action                                  |
| ------------------------------ | --------------------------------------- |
| `1`–`5`                        | Select tool (tile / obstacle / player spawn / enemy spawn / erase) |
| Left click on palette          | Pick a tile or obstacle type            |
| Left click on grid             | Paint / place the selected tool         |
| Right click on grid            | Erase (reset to grass, remove stuff)    |
| `S`                            | Save the map to the maps directory      |
| `N`                            | Start a new blank map                    |
| `Esc`                          | Save and play the edited map            |

Edited maps are saved to a runtime directory (`maps/` by default, override with
`-mapsdir <dir>`) as `<name>.json` and become available in single player (e.g.
`go run ./cmd -map custom`) without recompiling. To author several maps, rename
the generated files.

## Flags

| Flag            | Description                                           |
| --------------- | ----------------------------------------------------- |
| `-map <name>`   | Pick the map (single player; ignored in multiplayer)  |
| `-addr <host>`  | Connect to a multiplayer server (`host:port`)         |
| `-mapsdir <dir>`| Directory for runtime/edited maps (default `maps`)     |
| `-cpuprofile`   | Write a CPU profile to a file                         |

## Tests

```bash
go test ./...
```

The suite validates every map against the sprite atlas, verifies the tank AI
geometry (heading, barrel aim, ballistic slope), the network protocol, and the
editor save/parse round trip.

## Project layout

```
cmd/                game entry point (flags, ebiten run loop)
game/               game logic: tanks, bullets, AI, playfield, network client, editor
server/             multiplayer hub server (state relay + AI simulation)
internal/protocol/  shared wire messages (JSON over WebSocket)
internal/enemyai/   pure steering/ballistic math reused by client and server
internal/maps/      hand-authored level format + embedded maps
internal/assets/    sprite-sheet atlas loader
```
