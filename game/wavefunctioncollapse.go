package game

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/runozo/go-tanks/assets"
)

const (
	ruleUP = iota
	ruleRIGHT
	ruleDOWN
	ruleLEFT
	ruleWEIGHT
)

// filterOptions filters the original options based on the provided options slice.
//
// It takes in two parameters:
// - orig []string: the original options slice
// - options []string: the options to filter by
// Returns []string: the filtered options slice
func filterOptions(orig, options []string) []string {
	var filtered []string
	for _, o := range orig {
		if stringInSlice(o, options) {
			filtered = append(filtered, o)
		}
	}
	// slog.Info("Filtered options", "filtered", filtered)
	return filtered
}

// stringInSlice checks if a string is present in a slice of strings.
//
// It takes a string to search for and a slice of strings to search in and returns a boolean.
func stringInSlice(a string, slice []string) bool {
	for _, b := range slice {
		if b == a {
			return true
		}
	}
	return false
}

// intInSlice checks if an integer is in a slice of integers.
//
// a int - the integer to check for in the slice
// list []int - the slice of integers to search
// bool - true if the integer is found in the slice, false otherwise
func intInSlice(a int, slice []int) bool {
	for _, b := range slice {
		if b == a {
			return true
		}
	}
	return false
}

// getMinEntropyIndexes returns the indexes of cells with the minimum entropy.
//
// It takes a pointer to a 2D slice of strings, cells, as input.
// The function iterates through each cell in the cells slice and checks if the length of the cell is greater than 1 and less than the current minimum entropy.
// If so, it updates the minimum entropy and resets the minEntropyIndexes slice to contain only the current index.
// If the length of the cell is equal to the current minimum entropy, the index is appended to the minEntropyIndexes slice.
// Finally, the function returns the minEntropyIndexes slice.
//
// Parameters:
// - cells: a pointer to a 2D slice of strings representing the cells
//
// Return type:
// - []int: a slice of integers representing the indexes of cells with the minimum entropy
func getMinEntropyIndexes(tiles *[]Tile) []int {
	minEntropy := 32767
	minEntropyIndexes := make([]int, 0)
	for i, tile := range *tiles {
		if !tile.collapsed {
			cellEntropy := len(tile.options)
			// slog.Info("Entropy", "index", i, "entropy", cellEntropy)
			if cellEntropy > 1 && cellEntropy < minEntropy {
				minEntropy = cellEntropy
				minEntropyIndexes = []int{i}
			} else if cellEntropy > 1 && cellEntropy == minEntropy {
				minEntropyIndexes = append(minEntropyIndexes, i)
			}
		}
	}
	return minEntropyIndexes
}

// collapseRandomCellWithMinEntropy collapses a random cell with the minimum entropy.
//
// Parameters:
// - tiles: a pointer to a slice of Tile representing the game tiles
// - minEntropyIndexes: a pointer to a slice of integers representing the indexes of cells with the minimum entropy
//
// Return type:
// - int: the index of the collapsed cell
func collapseRandomCellWithMinEntropy(tiles *[]Tile, minEntropyIndexes *[]int) int {
	// collapse random cell with least entropy
	index := (*minEntropyIndexes)[rand.Intn(len(*minEntropyIndexes))]

	(*tiles)[index].options = []string{(*tiles)[index].options[rand.Intn(len((*tiles)[index].options))]}
	return index
}

// lookAndFilter applies a rule-based filtering to the optionsToProcess slice
//
// It takes two integer rule indexes, two slices of strings (optionsToProcess and optionsToWatch) as parameters
// Returns a slice of strings
func lookAndFilter(ruleIndexToProcess, ruleIndexToWatch int, optionsToProcess, optionsToWatch []string) []string {
	rules := make([]int, 0, 5) // random capacity
	for _, optname := range optionsToWatch {
		rule := assets.RulesGround[optname][ruleIndexToWatch]
		rules = append(rules, rule)
	}

	newoptions := make([]string, 0, 5) // random capacity
	for k, v := range assets.RulesGround {
		if intInSlice(v[ruleIndexToProcess], rules) {
			newoptions = append(newoptions, k)
		}
	}

	return filterOptions(optionsToProcess, newoptions)
}

func resetTilesOptions(tiles *[]Tile) {
	// create a slice of all the options available
	totaloptions := 0
	for _, v := range assets.RulesGround {
		totaloptions += v[ruleWEIGHT]
	}
	initialOptions := make([]string, totaloptions)
	i := 0
	for k, v := range assets.RulesGround {
		for j := 0; j < v[ruleWEIGHT]; j++ {
			initialOptions[i] = k
			i++
		}
	}

	// setup tiles with all the options enabled and a black square as image
	for i := 0; i < len(*tiles); i++ {
		(*tiles)[i].options = initialOptions
		(*tiles)[i].image = ebiten.NewImage(tileWidth, tileHeight)
		(*tiles)[i].collapsed = false
	}
}
