package game

import "fmt"

// Wave collapse algorithm

func setup() {
	tiles := []string{
		"tileGrass_roadTransitionN_dirt.png",
		"tileGrass_roadTransitionS.png",
		"tileGrass_roadTransitionS_dirt.png",
		"tileGrass_roadTransitionW.png",
		"tileGrass_roadTransitionW_dirt.png",
		"tileGrass_transitionE.png",
		"tileGrass_transitionN.png",
		"tileGrass_transitionS.png",
		"tileGrass_transitionW.png",
		"tileSand1.png",
		"tileSand2.png",
		"tileSand_roadCornerLL.png",
		"tileSand_roadCornerLR.png",
		"tileSand_roadCornerUL.png",
		"tileSand_roadCornerUR.png",
		"tileSand_roadCrossing.png",
		"tileSand_roadCrossingRound.png",
		"tileSand_roadEast.png",
		"tileSand_roadNorth.png",
		"tileSand_roadSplitE.png",
		"tileSand_roadSplitN.png",
		"tileSand_roadSplitS.png",
		"tileSand_roadSplitW.png",
	}
	fmt.Println(tiles)
}
