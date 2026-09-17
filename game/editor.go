package game

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/runozo/go-tanks/internal/maps"
)

// Editor tools.
const (
	toolTile = iota
	toolObstacle
	toolPlayerSpawn
	toolEnemySpawn
	toolErase
)

const (
	editorCellSize = 44
	editorOriginX  = 8
	editorOriginY  = 8
	paletteCell    = 44
	paletteCols    = 4
)

// Grass tile used as the "erase"/base cell.
const baseGrass = "tileGrass1"

// paletteEntry is one clickable item in the editor palette.
type paletteEntry struct {
	label string // display char for tool buttons, or "" for sprites
	tile  string // tile/obstacle name, or tool key ("player","enemy","erase")
	kind  int
}

type Editor struct {
	game   *Game
	name   string
	width  int
	height int
	tiles  [][]string

	obstacles   []ObstacleSpawn
	playerSpawn *SpawnPoint
	enemySpawns []SpawnPoint

	tool    int
	selTile string
	selObst string

	palette []paletteEntry
	status  string
}

func NewEditor(game *Game) *Editor {
	e := &Editor{
		game:    game,
		selTile: baseGrass,
		selObst: "crateWood",
		status:  "Map editor",
	}
	e.buildPalette()
	e.New(31, 17, "custom")
	return e
}

// buildPalette collects the terrain tiles, furniture and special tools.
func (e *Editor) buildPalette() {
	e.palette = e.palette[:0]

	// special tools
	e.palette = append(e.palette,
		paletteEntry{label: "P", tile: "player", kind: toolPlayerSpawn},
		paletteEntry{label: "E", tile: "enemy", kind: toolEnemySpawn},
		paletteEntry{label: "X", tile: "erase", kind: toolErase},
	)

	// terrain tiles (ground types), sorted for a stable layout
	var ground []string
	var furniture []string
	for name, t := range e.game.assets.TileEntries {
		if t.Type == "ground" {
			ground = append(ground, name)
		} else if isFurniture(name) {
			furniture = append(furniture, name)
		}
	}
	sort.Strings(ground)
	sort.Strings(furniture)

	for _, n := range ground {
		e.palette = append(e.palette, paletteEntry{tile: n, kind: toolTile})
	}
	for _, n := range furniture {
		e.palette = append(e.palette, paletteEntry{tile: n, kind: toolObstacle})
	}
}

func isFurniture(name string) bool {
	for _, prefix := range []string{
		"crate", "barrel", "tree", "sandbag", "barricade", "oilSpill", "rock",
	} {
		if len(name) >= len(prefix) && name[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

// New resets the editor to a blank grass grid.
func (e *Editor) New(width, height int, name string) {
	e.width, e.height, e.name = width, height, name
	e.tiles = make([][]string, height)
	for y := 0; y < height; y++ {
		e.tiles[y] = make([]string, width)
		for x := 0; x < width; x++ {
			e.tiles[y][x] = baseGrass
		}
	}
	e.obstacles = nil
	e.playerSpawn = nil
	e.enemySpawns = nil
}

// LoadFromMap fills the editor from a parsed map.
func (e *Editor) LoadFromMap(m *maps.Map) {
	e.name = m.Name
	e.width, e.height = m.Width, m.Height
	e.tiles = make([][]string, e.height)
	for y := 0; y < e.height; y++ {
		e.tiles[y] = append([]string(nil), m.Tiles[y]...)
	}
	e.obstacles = append([]ObstacleSpawn(nil), m.Obstacles...)
	if m.PlayerSpawn != nil {
		sp := *m.PlayerSpawn
		e.playerSpawn = &sp
	} else {
		e.playerSpawn = nil
	}
	e.enemySpawns = append([]SpawnPoint(nil), m.EnemySpawns...)
}

// ToMap builds a serializable Map, generating a legend for the used tiles.
func (e *Editor) ToMap(name string) *maps.Map {
	// map unique tile name -> legend char
	used := map[string]bool{}
	for y := 0; y < e.height; y++ {
		for x := 0; x < e.width; x++ {
			used[e.tiles[y][x]] = true
		}
	}
	chars := []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	legend := map[string]string{}
	rev := map[string]string{}
	i := 0
	for name := range used {
		ch := string(chars[i%len(chars)])
		i++
		legend[ch] = name
		rev[name] = ch
	}

	grid := make([]string, e.height)
	for y := 0; y < e.height; y++ {
		row := make([]rune, e.width)
		for x := 0; x < e.width; x++ {
			row[x] = []rune(rev[e.tiles[y][x]])[0]
		}
		grid[y] = string(row)
	}

	m := &maps.Map{
		Name:        name,
		Legend:      legend,
		Grid:        grid,
		Obstacles:   append([]maps.ObstacleSpawn(nil), e.obstacles...),
		PlayerSpawn: e.playerSpawn,
		EnemySpawns: append([]maps.SpawnPoint(nil), e.enemySpawns...),
		Width:       e.width,
		Height:      e.height,
	}
	return m
}

// Save writes the map to the runtime maps directory as <name>.json, creating
// the directory when needed.
func (e *Editor) Save() error {
	if e.game.mapsDir == "" {
		e.status = "no mapsdir configured"
		return fmt.Errorf("no mapsdir")
	}
	if err := os.MkdirAll(e.game.mapsDir, 0o755); err != nil {
		return err
	}

	m := e.ToMap(e.name)
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(e.game.mapsDir, e.name+".json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return err
	}
	e.status = fmt.Sprintf("saved %s", path)
	e.game.reloadRuntimeMaps()
	return nil
}

func (e *Editor) cellAt(mx, my int) (int, int, bool) {
	x := (mx - editorOriginX) / editorCellSize
	y := (my - editorOriginY) / editorCellSize
	if y >= 0 && y < e.height && x >= 0 && x < e.width {
		return x, y, true
	}
	return 0, 0, false
}

func (e *Editor) paletteRect(i int) (x, y int, ok bool) {
	if i < 0 || i >= len(e.palette) {
		return 0, 0, false
	}
	px := editorOriginX + e.width*editorCellSize + 24
	py := editorOriginY
	col := i % paletteCols
	row := i / paletteCols
	return px + col*paletteCell, py + row*paletteCell, true
}

func (e *Editor) paletteAt(mx, my int) int {
	px := editorOriginX + e.width*editorCellSize + 24
	py := editorOriginY
	col := (mx - px) / paletteCell
	row := (my - py) / paletteCell
	if col < 0 || col >= paletteCols || row < 0 {
		return -1
	}
	i := row*paletteCols + col
	if i >= len(e.palette) {
		return -1
	}
	return i
}

func (e *Editor) Update(tps float64) {
	// tool selection via number keys
	for i, k := range []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4, ebiten.Key5} {
		if inpututil.IsKeyJustPressed(k) {
			e.tool = i
			e.status = fmt.Sprintf("tool %d", i)
		}
	}

	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyS):
		if err := e.Save(); err != nil {
			e.status = "save: " + err.Error()
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyN):
		e.New(31, 17, fmt.Sprintf("custom_%d", time.Now().Unix()))
		e.status = "new map: " + e.name
	case inpututil.IsKeyJustPressed(ebiten.KeyEscape):
		e.game.exitEditor(true)
		return
	}

	mx, my := ebiten.CursorPosition()

	// palette click
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if i := e.paletteAt(mx, my); i >= 0 {
			p := e.palette[i]
			switch p.kind {
			case toolTile:
				e.tool = toolTile
				e.selTile = p.tile
				e.status = "tile: " + p.tile
			case toolObstacle:
				e.tool = toolObstacle
				e.selObst = p.tile
				e.status = "obstacle: " + p.tile
			default:
				e.tool = p.kind
				e.status = "tool: " + p.tile
			}
			return
		}
	}

	// grid painting
	if cx, cy, ok := e.cellAt(mx, my); ok {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			e.paint(cx, cy)
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
			e.erase(cx, cy)
		}
	}
}

func (e *Editor) paint(cx, cy int) {
	switch e.tool {
	case toolTile:
		e.tiles[cy][cx] = e.selTile
	case toolObstacle:
		e.obstacles = append(e.obstacles, ObstacleSpawn{X: float64(cx * 64), Y: float64(cy * 64), Type: e.selObst})
	case toolPlayerSpawn:
		e.playerSpawn = &SpawnPoint{X: float64(cx*64 + 32), Y: float64(cy*64 + 32)}
	case toolEnemySpawn:
		e.enemySpawns = append(e.enemySpawns, SpawnPoint{X: float64(cx*64 + 32), Y: float64(cy*64 + 32)})
	case toolErase:
		e.erase(cx, cy)
	}
}

func (e *Editor) erase(cx, cy int) {
	e.tiles[cy][cx] = baseGrass

	px, py := float64(cx*64), float64(cy*64)
	keep := e.obstacles[:0]
	for _, o := range e.obstacles {
		if !(o.X >= px && o.X < px+64 && o.Y >= py && o.Y < py+64) {
			keep = append(keep, o)
		}
	}
	e.obstacles = keep

	enemies := e.enemySpawns[:0]
	for _, s := range e.enemySpawns {
		if !(s.X >= px && s.X < px+64 && s.Y >= py && s.Y < py+64) {
			enemies = append(enemies, s)
		}
	}
	e.enemySpawns = enemies

	if e.playerSpawn != nil && e.playerSpawn.X >= px && e.playerSpawn.X < px+64 && e.playerSpawn.Y >= py && e.playerSpawn.Y < py+64 {
		e.playerSpawn = nil
	}
}

func (e *Editor) drawScaled(screen *ebiten.Image, name string, x, y, size float64) {
	spr := e.game.assets.GetSprite(name)
	if spr == nil {
		return
	}
	w := float64(spr.Bounds().Dx())
	h := float64(spr.Bounds().Dy())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(size/w, size/h)
	op.GeoM.Translate(x, y)
	screen.DrawImage(spr, op)
}

func (e *Editor) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 34, 255})

	// grid tiles
	for y := 0; y < e.height; y++ {
		for x := 0; x < e.width; x++ {
			px := float64(editorOriginX + x*editorCellSize)
			py := float64(editorOriginY + y*editorCellSize)
			e.drawScaled(screen, e.tiles[y][x], px, py, editorCellSize)
		}
	}

	// obstacles
	for _, o := range e.obstacles {
		cx := int(o.X) / 64
		cy := int(o.Y) / 64
		if cx < 0 || cx >= e.width || cy < 0 || cy >= e.height {
			continue
		}
		px := float64(editorOriginX + cx*editorCellSize)
		py := float64(editorOriginY + cy*editorCellSize)
		e.drawScaled(screen, o.Type, px, py, editorCellSize)
	}

	// spawn markers
	if e.playerSpawn != nil {
		cx := int(e.playerSpawn.X) / 64
		cy := int(e.playerSpawn.Y) / 64
		px := float32(editorOriginX + cx*editorCellSize)
		py := float32(editorOriginY + cy*editorCellSize)
		vector.StrokeRect(screen, px+2, py+2, editorCellSize-4, editorCellSize-4, 3, color.RGBA{32, 255, 128, 255}, false)
	}
	for _, s := range e.enemySpawns {
		cx := int(s.X) / 64
		cy := int(s.Y) / 64
		px := float32(editorOriginX + cx*editorCellSize)
		py := float32(editorOriginY + cy*editorCellSize)
		vector.StrokeRect(screen, px+2, py+2, editorCellSize-4, editorCellSize-4, 3, color.RGBA{255, 80, 80, 255}, false)
	}

	// hover highlight
	mx, my := ebiten.CursorPosition()
	if cx, cy, ok := e.cellAt(mx, my); ok {
		px := float32(editorOriginX + cx*editorCellSize)
		py := float32(editorOriginY + cy*editorCellSize)
		vector.StrokeRect(screen, px, py, editorCellSize, editorCellSize, 1, color.RGBA{255, 255, 255, 200}, false)
	}

	// palette
	for i, p := range e.palette {
		px, py, ok := e.paletteRect(i)
		if !ok {
			continue
		}
		fx := float32(px)
		fy := float32(py)
		switch p.kind {
		case toolPlayerSpawn:
			vector.DrawFilledRect(screen, fx, fy, paletteCell, paletteCell, color.RGBA{32, 180, 96, 255}, false)
		case toolEnemySpawn:
			vector.DrawFilledRect(screen, fx, fy, paletteCell, paletteCell, color.RGBA{180, 60, 60, 255}, false)
		case toolErase:
			vector.DrawFilledRect(screen, fx, fy, paletteCell, paletteCell, color.RGBA{90, 90, 96, 255}, false)
		default:
			e.drawScaled(screen, p.tile, float64(px), float64(py), paletteCell)
		}
		if p.label != "" {
			drawCenteredText(screen, p.label, fx+paletteCell/2, fy+paletteCell/2, e.game.fontSmall)
		}
		if p.tile == e.selTile && p.kind == toolTile || p.tile == e.selObst && p.kind == toolObstacle {
			vector.StrokeRect(screen, fx, fy, paletteCell, paletteCell, 2, color.RGBA{255, 255, 0, 255}, false)
		}
		if p.kind == e.tool && p.kind > toolObstacle {
			vector.StrokeRect(screen, fx, fy, paletteCell, paletteCell, 2, color.RGBA{255, 255, 0, 255}, false)
		}
	}

	// help text
	help := "1-5 tools  |  LMB paint  RMB erase  |  S save  N new  Esc play  |  " + e.status
	drawText(screen, help, 8, screenHeight-40, e.game.fontSmall, color.White)
}

func drawText(screen *ebiten.Image, str string, x, y int, f *text.GoTextFace, c color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(c)
	text.Draw(screen, str, f, op)
}

func drawCenteredText(screen *ebiten.Image, str string, x, y float32, f *text.GoTextFace) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, str, f, op)
}
