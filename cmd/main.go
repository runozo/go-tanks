package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/runozo/go-tanks/game"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var serveraddress = flag.String("addr", "", "http service address")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	}

	g := game.NewGame(*serveraddress)
	ebiten.SetFullscreen(true)
	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
