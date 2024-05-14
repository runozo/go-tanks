package game

var tileOptions = map[string][]int{
	"tileGrass1.png":                     {0, 0, 0, 0}, // 0 grass
	"tileGrass2.png":                     {0, 0, 0, 0},
	"tileGrass_roadCornerLL.png":         {0, 0, 1, 1}, // 1 road with grass
	"tileGrass_roadCornerLR.png":         {0, 1, 1, 0},
	"tileGrass_roadCornerUL.png":         {1, 0, 0, 1},
	"tileGrass_roadCornerUR.png":         {1, 1, 0, 0},
	"tileGrass_roadCrossing.png":         {1, 1, 1, 1},
	"tileGrass_roadCrossingRound.png":    {1, 1, 1, 1},
	"tileGrass_roadEast.png":             {0, 1, 0, 1},
	"tileGrass_roadNorth.png":            {1, 0, 1, 0},
	"tileGrass_roadSplitE.png":           {1, 1, 1, 0},
	"tileGrass_roadSplitN.png":           {1, 1, 0, 1},
	"tileGrass_roadSplitS.png":           {0, 1, 1, 1},
	"tileGrass_roadSplitW.png":           {1, 0, 1, 1},
	"tileGrass_roadTransitionE.png":      {4, 3, 4, 1}, //
	"tileGrass_roadTransitionE_dirt.png": {4, 3, 4, 1},
	"tileGrass_roadTransitionN.png":      {3, 6, 1, 6},
	"tileGrass_roadTransitionN_dirt.png": {3, 6, 1, 6},
	"tileGrass_roadTransitionS.png":      {1, 8, 3, 8},
	"tileGrass_roadTransitionS_dirt.png": {1, 8, 3, 8},
	"tileGrass_roadTransitionW.png":      {5, 1, 5, 3},
	"tileGrass_roadTransitionW_dirt.png": {5, 1, 5, 3},
	"tileGrass_transitionE.png":          {4, 2, 4, 0},
	"tileGrass_transitionN.png":          {2, 6, 0, 6},
	"tileGrass_transitionS.png":          {0, 8, 2, 8},
	"tileGrass_transitionW.png":          {5, 0, 5, 2},
	"tileSand1.png":                      {2, 2, 2, 2},
	"tileSand2.png":                      {2, 2, 2, 2},
	"tileSand_roadCornerLL.png":          {2, 2, 3, 3},
	"tileSand_roadCornerLR.png":          {2, 3, 3, 2},
	"tileSand_roadCornerUL.png":          {3, 2, 2, 3},
	"tileSand_roadCornerUR.png":          {3, 3, 2, 2},
	"tileSand_roadCrossing.png":          {3, 3, 3, 3},
	"tileSand_roadCrossingRound.png":     {3, 3, 3, 3},
	"tileSand_roadEast.png":              {2, 3, 2, 3},
	"tileSand_roadNorth.png":             {3, 2, 3, 2},
	"tileSand_roadSplitE.png":            {3, 3, 3, 2},
	"tileSand_roadSplitN.png":            {3, 3, 2, 3},
	"tileSand_roadSplitS.png":            {2, 3, 3, 3},
	"tileSand_roadSplitW.png":            {3, 2, 3, 3},
}

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
	return filtered
}

// stringInSlice checks if a string is present in a slice of strings.
//
// It takes a string to search for and a slice of strings to search in and returns a boolean.
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
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
func intInSlice(a int, list []int) bool {
	for _, b := range list {
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
func getMinEntropyIndexes(cells *[][]string) []int {
	minEntropy := 32767
	var minEntropyIndexes []int
	for i, cell := range *cells {
		if len(cell) > 1 && len(cell) < minEntropy {
			minEntropy = len(cell)
			minEntropyIndexes = []int{i}
		} else if len(cell) > 1 && len(cell) == minEntropy {
			minEntropyIndexes = append(minEntropyIndexes, i)
		}
	}
	return minEntropyIndexes
}
