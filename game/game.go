package game

import (
	"bytes"
	"embed"
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/runozo/go-tanks/internal/assets"
	"github.com/runozo/go-tanks/internal/maps"
	"github.com/solarlune/resolv"
)

const (
	screenWidth     = 1960
	screenHeight    = 1088
	fontSizeMedium  = 42
	fontSizeSmall   = 22
	numberOfEnemies = 5 // *3
	cellWidth       = 64
	cellHeight      = 64
)

//go:embed assets/*
var assetsFS embed.FS

var (
	TagPlayer    = resolv.NewTag("Player")
	TagBarrel    = resolv.NewTag("Barrel")
	TagEnemy     = resolv.NewTag("Enemy")
	TagBullet    = resolv.NewTag("Bullet")
	TagExplosion = resolv.NewTag("Explosion")
	TagObstacle  = resolv.NewTag("Obstacle")
)

func FlipVertical(source *ebiten.Image) *ebiten.Image {
	flipped := ebiten.NewImage(source.Bounds().Dx(), source.Bounds().Dy())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(1, -1)
	op.GeoM.Translate(0, float64(source.Bounds().Dy()))
	flipped.DrawImage(source, op)
	return flipped
}

type Game struct {
	width         int
	height        int
	assets        *assets.Assets
	Tanks         []*Tank
	playfield     *Playfield
	scoreLine     *ScoreLine
	editor        *Editor
	obstacles     []*Obstacle
	fontMedium    *text.GoTextFace
	fontSmall     *text.GoTextFace
	serverAddress string
	mapsDir       string
	netClient     *NetClient
	networkTanks  map[string]*Tank
	netEnemies    map[string]*Tank // server-simulated enemies (multiplayer)
	localPlayer   *Tank
	myClientID    string
	netTick       uint64
	explosions    []*Explosion // standalone explosions (multiplayer fallback)
	decals        []TrackDecal // tank-track imprints on soft terrain
	space         *resolv.Space
	maps          []*Map
	mapName       string
	state         int
	debug         bool
}

// game state
const (
	PLAYING = iota
	PAUSED
	RENDERINGPLAYFIELD
	GAMEOVER
	HELP
	OPTIONS
	INTRO
	EDITING
)

func NewGame(serverAddress, mapName, mapsDir string) *Game {
	// Load sprite sheet
	spriteSheetData, err := assetsFS.ReadFile("assets/allSprites_default.png")
	if err != nil {
		log.Fatal(err)
	}

	// Load tile map
	jsonData, err := assetsFS.ReadFile("assets/mapped_tiles.json")
	if err != nil {
		log.Fatal(err)
	}
	assets := assets.NewAssets(
		spriteSheetData,
		jsonData,
	)

	mapsData, err := loadMaps(assets, mapsDir)
	if err != nil {
		log.Fatal(err)
	}

	// Load fonts
	fontData, err := assetsFS.ReadFile("assets/gomarice_no_continue.ttf")
	if err != nil {
		log.Fatal(err)
	}

	textFS, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	g := &Game{
		assets:    assets,
		width:     screenWidth,
		height:    screenHeight,
		playfield: nil,
		obstacles: []*Obstacle{},
		fontMedium: &text.GoTextFace{
			Source:    textFS,
			Direction: text.DirectionLeftToRight,
			Size:      fontSizeMedium,
		},
		fontSmall: &text.GoTextFace{
			Source:    textFS,
			Direction: text.DirectionLeftToRight,
			Size:      fontSizeSmall,
		},
		serverAddress: serverAddress,
		mapsDir:       mapsDir,
		networkTanks:  make(map[string]*Tank),
		netEnemies:    make(map[string]*Tank),
		space:         resolv.NewSpace(screenWidth, screenHeight, cellWidth, cellHeight),
		maps:          mapsData,
		mapName:       mapName,
		state:         RENDERINGPLAYFIELD,
		debug:         false,
	}

	if serverAddress != "" {
		// Multiplayer: the session (map, enemies, players) comes from the
		// server snapshot; the playfield is created when it arrives.
		g.myClientID = uuid.NewString()
		g.netClient = NewNetClient(serverAddress, g.myClientID)
	} else {
		// Single player: build the playfield immediately.
		g.playfield = NewPlayfield(g)
	}

	return g
}

// loadMaps reads the embedded maps plus any runtime maps in dir, validating
// them against the available sprite map.
func loadMaps(sprites *assets.Assets, dir string) ([]*Map, error) {
	all, err := maps.All()
	if err != nil {
		return nil, err
	}

	if dir != "" {
		runtimeMaps, err := maps.LoadDir(dir)
		if err != nil {
			return nil, fmt.Errorf("loading runtime maps: %w", err)
		}
		all = append(all, runtimeMaps...)
	}

	for _, m := range all {
		if err := m.Validate(func(name string) bool {
			_, ok := sprites.TileEntries[name]
			return ok
		}); err != nil {
			return nil, err
		}
	}

	return all, nil
}

// pickMap returns the map selected by g.mapName if set, otherwise a random one.
// For a named map, the last match wins so that runtime maps (appended after the
// embedded ones) can override a built-in map with the same name.
func (g *Game) pickMap() *Map {
	if g.mapName != "" {
		var found *Map
		for _, m := range g.maps {
			if m.Name == g.mapName {
				found = m
			}
		}
		if found != nil {
			return found
		}
		log.Fatalf("map %q not found", g.mapName)
	}
	return g.maps[rand.Intn(len(g.maps))]
}

func (g *Game) Update() error {
	tps := float64(ebiten.TPS())

	// Map editor toggle (single player).
	if inpututil.IsKeyJustPressed(ebiten.KeyE) && g.netClient == nil && g.state != EDITING {
		g.enterEditor()
		return nil
	}
	if g.state == EDITING {
		g.editor.Update(tps)
		return nil
	}

	if g.netClient != nil {
		g.processNetMessages()
	}

	if g.state == RENDERINGPLAYFIELD {
		if g.playfield != nil {
			g.playfield.Update(tps)
		}

		// prepare the game once the map is ready (single player: immediately,
		// multiplayer: after the server snapshot builds the playfield).
		if g.playfield != nil && g.playfield.ready && g.state == RENDERINGPLAYFIELD {
			g.setupSinglePlayer()
		}
	}

	// recreate playfield (single player only; the network session is owned
	// by the server)
	if g.netClient == nil && inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.state = RENDERINGPLAYFIELD
		g.space.RemoveAll()
		g.obstacles = []*Obstacle{}
		g.Tanks = nil
		g.decals = nil
		g.playfield = NewPlayfield(g)
	}

	if g.state == PLAYING {
		if g.netClient != nil {
			g.sendLocalTransform()
		}

		// update tanks (local player input, remote tank interpolation)
		for _, t := range g.Tanks {
			t.Update(tps)
		}

		// update server-driven enemies (visuals and replicated bullets only)
		if g.netClient != nil {
			for _, e := range g.netEnemies {
				e.updateWeapons(tps)
			}
		}

		g.updateDecals(tps)

		// standalone explosions (multiplayer fallback)
		active := g.explosions[:0]
		for _, e := range g.explosions {
			e.Update(tps)
			if !e.bullet.exploded {
				active = append(active, e)
			}
		}
		g.explosions = active
	}

	return nil
}

// setupSinglePlayer builds the single-player match once the playfield is ready.
func (g *Game) setupSinglePlayer() {
	g.decals = nil

	m := g.playfield.mapData

	// add obstacles (authored in the map)
	for _, o := range m.Obstacles {
		g.obstacles = append(g.obstacles, NewObstacleAt(g, o))
	}

	// add player at the authored spawn (or a random position as fallback)
	if m.PlayerSpawn != nil {
		g.Tanks = []*Tank{NewTankAt(g, resolv.Vector{X: m.PlayerSpawn.X, Y: m.PlayerSpawn.Y}, m.PlayerSpawn.Rotation, false)}
	} else {
		g.Tanks = []*Tank{NewRandomTank(g, 0, false)}
	}

	// add enemies at the authored spawns (or random as fallback)
	if len(m.EnemySpawns) > 0 {
		for _, s := range m.EnemySpawns {
			g.Tanks = append(g.Tanks, NewTankAt(g, resolv.Vector{X: s.X, Y: s.Y}, s.Rotation, true))
		}
	} else {
		for i := 0; i < numberOfEnemies; i++ {
			g.Tanks = append(g.Tanks, NewRandomTank(g, 0, true))
		}
	}

	g.scoreLine = NewScoreLine(g)

	g.state = PLAYING
}

// ---------------------------------------------------------------------------
// Map editor
// ---------------------------------------------------------------------------

func (g *Game) enterEditor() {
	if g.editor == nil {
		g.editor = NewEditor(g)
	}
	if g.playfield != nil {
		g.editor.LoadFromMap(g.playfield.mapData)
	}
	g.state = EDITING
}

// exitEditor leaves the editor; when play is true, it persists the edited map
// and rebuilds the playfield so it can be played immediately.
func (g *Game) exitEditor(play bool) {
	m := g.editor.ToMap(g.editor.name)
	g.upsertMap(m)
	g.mapName = m.Name

	if !play {
		g.state = PLAYING
		return
	}

	if g.mapsDir != "" {
		if err := g.editor.Save(); err != nil {
			log.Printf("save map: %v", err)
		}
	}

	g.space.RemoveAll()
	g.obstacles = []*Obstacle{}
	g.Tanks = nil
	g.decals = nil
	g.playfield = NewPlayfield(g)
	g.state = RENDERINGPLAYFIELD
}

// reloadRuntimeMaps re-reads embedded + runtime maps (used after a save).
func (g *Game) reloadRuntimeMaps() {
	all, err := loadMaps(g.assets, g.mapsDir)
	if err != nil {
		log.Printf("reload maps: %v", err)
		return
	}
	g.maps = all
}

// upsertMap inserts the map into the in-memory list, replacing by name.
func (g *Game) upsertMap(m *Map) {
	for i, existing := range g.maps {
		if existing.Name == m.Name {
			g.maps[i] = m
			return
		}
	}
	g.maps = append(g.maps, m)
}

// ---------------------------------------------------------------------------
// Networking (multiplayer)
// ---------------------------------------------------------------------------

// processNetMessages drains the incoming network queue.
func (g *Game) processNetMessages() {
	for {
		select {
		case msg := <-g.netClient.Incoming:
			g.handleNetMessage(msg)
		default:
			return
		}
	}
}

func (g *Game) handleNetMessage(msg NetMessage) {
	switch msg.Type {
	case msgSnapshot:
		g.applySnapshot(msg)

	case msgTransform:
		g.ensureNetTank(msg.ClientID, msg.X, msg.Y, msg.R)

	case msgFire:
		g.applyRemoteFire(msg)

	case msgEnemyState:
		g.applyEnemyState(msg)

	case msgEnemyFire:
		g.applyEnemyFire(msg)

	case msgEnemyHitPlayer:
		g.applyEnemyHitPlayer(msg)

	case msgScore:
		if g.scoreLine != nil {
			g.scoreLine.Set(msg.Score.Player, msg.Score.Comp)
		}

	case msgLeave:
		if t, ok := g.networkTanks[msg.ClientID]; ok {
			t.Destroy()
			delete(g.networkTanks, msg.ClientID)
			for i, tt := range g.Tanks {
				if tt == t {
					g.Tanks = append(g.Tanks[:i], g.Tanks[i+1:]...)
					break
				}
			}
		}
	}
}

// applySnapshot builds the whole multiplayer session from the server state.
func (g *Game) applySnapshot(msg NetMessage) {
	if g.playfield != nil {
		return // session already started
	}

	// The server owns the map selection.
	g.mapName = msg.MapName
	g.playfield = NewPlayfield(g)
	g.decals = nil

	// obstacles (shared, authored in the map file)
	for _, o := range g.playfield.mapData.Obstacles {
		g.obstacles = append(g.obstacles, NewObstacleAt(g, o))
	}

	// enemy visuals, driven by the server
	for _, es := range msg.Enemies {
		t := NewTankAt(g, resolv.Vector{X: es.X, Y: es.Y}, es.R, true)
		t.ID = es.ID
		t.Object.SetRotation(es.R)
		if len(t.barrels) > 0 {
			t.barrels[0].relativeRotation = es.Aim - es.R
		}
		g.netEnemies[es.ID] = t
	}

	// remote players already in the session
	for _, p := range msg.Players {
		if p.ClientID == g.myClientID {
			continue
		}
		g.ensureNetTank(p.ClientID, p.X, p.Y, p.R)
	}

	// local player at the authored spawn
	sp := g.playfield.mapData.PlayerSpawn
	pos := resolv.Vector{X: float64(screenWidth) / 2, Y: float64(screenHeight) / 2}
	if sp != nil {
		pos = resolv.Vector{X: sp.X, Y: sp.Y}
	}
	g.localPlayer = NewTankAt(g, pos, 0, false)
	g.Tanks = append(g.Tanks, g.localPlayer)

	g.scoreLine = NewScoreLine(g)
	g.scoreLine.Set(msg.Score.Player, msg.Score.Comp)

	log.Printf("session started on map %q with %d players and %d enemies", g.mapName, len(msg.Players), len(msg.Enemies))
	g.state = PLAYING
}

// ensureNetTank returns the remote tank for clientID, creating it on first
// sight and updating its interpolation target.
func (g *Game) ensureNetTank(clientID string, x, y, r float64) *Tank {
	if t, ok := g.networkTanks[clientID]; ok {
		t.netTarget = resolv.Vector{X: x, Y: y}
		t.netTargetR = r
		return t
	}

	t := NewTankAt(g, resolv.Vector{X: x, Y: y}, r, false)
	t.IsRemote = true
	t.ID = clientID
	g.networkTanks[clientID] = t
	g.Tanks = append(g.Tanks, t)

	return t
}

// applyRemoteFire renders a projectile fired by a remote player.
func (g *Game) applyRemoteFire(msg NetMessage) {
	t := g.ensureNetTank(msg.ClientID, msg.X, msg.Y, msg.R)
	if len(t.barrels) == 0 {
		return
	}
	b := t.barrels[0]
	b.relativeRotation = msg.R - t.Object.Rotation()
	b.slope = msg.Slope
	bullet := b.Fire()
	bullet.visual = true
	t.Bullets = append(t.Bullets, bullet)
}

// applyEnemyState updates the server-driven enemy visuals.
func (g *Game) applyEnemyState(msg NetMessage) {
	for _, es := range msg.Enemies {
		e, ok := g.netEnemies[es.ID]
		if !ok {
			e = NewTankAt(g, resolv.Vector{X: es.X, Y: es.Y}, es.R, true)
			e.ID = es.ID
			g.netEnemies[es.ID] = e
		}
		e.Object.SetPositionVec(resolv.Vector{X: es.X, Y: es.Y})
		e.Object.SetRotation(es.R)
		if len(e.barrels) > 0 {
			e.barrels[0].relativeRotation = es.Aim - e.Object.Rotation()
		}
	}
}

// applyEnemyFire renders a projectile fired by a server-simulated enemy.
func (g *Game) applyEnemyFire(msg NetMessage) {
	e, ok := g.netEnemies[msg.EnemyID]
	if !ok || len(e.barrels) == 0 {
		return
	}
	b := e.barrels[0]
	b.relativeRotation = msg.R - e.Object.Rotation()
	b.slope = msg.Slope
	bullet := b.Fire()
	bullet.visual = true
	bullet.netID = msg.BulletID
	e.Bullets = append(e.Bullets, bullet)
}

// applyEnemyHitPlayer shows the impact of an enemy bullet on a player: the
// matching visual bullet is forced to land on the player, producing an
// explosion (Server-authoritative).
func (g *Game) applyEnemyHitPlayer(msg NetMessage) {
	var target *Tank
	if msg.ClientID == g.myClientID {
		target = g.localPlayer
	} else {
		target = g.networkTanks[msg.ClientID]
	}
	if target == nil {
		return
	}

	pos := target.Object.Center()

	// Force the matching enemy bullet to land right on the player.
	for _, e := range g.netEnemies {
		for _, b := range e.Bullets {
			if b.netID == msg.BulletID {
				b.solid.SetPositionVec(pos)
				b.hasHitTarget = true
				b.altitude = initialAltitude
				b.elapsedTime = 0.2
				return
			}
		}
	}

	// Fallback: if the bullet visual was not found, spawn an explosion at the
	// player's position using a dummy bullet.
	b := &Bullet{
		solid:        resolv.NewRectangle(pos.X, pos.Y, 8, 8),
		hasHitTarget: true,
		barrel:       &Barrel{tank: target},
	}
	e := NewExplosion(b)
	target.game.explosions = append(target.game.explosions, e)
}

// sendLocalTransform relays the local player's transform at 20Hz.
func (g *Game) sendLocalTransform() {
	if g.localPlayer == nil {
		return
	}
	g.netTick++
	if g.netTick%3 != 0 {
		return
	}
	center := g.localPlayer.Object.Center()
	g.netClient.SendTransform(center.X, center.Y, g.localPlayer.Object.Rotation())
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.state == EDITING {
		g.editor.Draw(screen)
		return
	}
	if g.playfield != nil {
		g.playfield.Draw(screen)
	}

	// tank tracks are drawn just above the terrain
	g.drawDecals(screen)

	if g.state == PLAYING {

		for _, f := range g.obstacles {
			f.Draw(screen)
		}

		// server-driven enemies (multiplayer) are drawn above the terrain
		if g.netClient != nil {
			for _, e := range g.netEnemies {
				e.Draw(screen)
			}
		}

		for _, e := range g.Tanks {
			e.Draw(screen)
		}

		// standalone explosions
		for _, e := range g.explosions {
			e.Draw(screen)
		}
		// Draw bullet above all
		for _, t := range g.Tanks {
			for _, bullet := range t.Bullets {
				bullet.Draw(screen)
			}
		}

		g.scoreLine.Draw(screen)
	}

	// debug shapes
	if g.debug {
		g.space.ForEachShape(func(shape resolv.IShape, index, maxCount int) bool {

			var drawColor color.Color = color.White

			tags := shape.Tags()

			if tags.Has(TagEnemy) && !tags.Has(TagBarrel) {
				drawColor = color.RGBA{255, 128, 35, 255}
			}
			if tags.Has(TagPlayer) {
				drawColor = color.RGBA{32, 255, 128, 255}
			}
			switch o := shape.(type) {
			case *resolv.Circle:
				vector.StrokeCircle(screen, float32(o.Position().X), float32(o.Position().Y), float32(o.Radius()), 2, drawColor, false)
			case *resolv.ConvexPolygon:

				for _, l := range o.Lines() {
					vector.StrokeLine(screen, float32(l.Start.X), float32(l.Start.Y), float32(l.End.X), float32(l.End.Y), 2, drawColor, false)
				}
			}

			return true

		})
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
