package pricing

/*************************************************************************************************************/
/* Domestic Line Haul Prices Types

Used for:

2) Domestic Price Tabs
        2a) Domestic Linehaul Prices
	    2b) Domestic Service Area Prices
	    2c) Other Domestic Prices
*/
/*************************************************************************************************************/

const dlhWeightBandNumCellsExpected int = 10 //cells per band verify against dlhWeightBandNumCells
const dlhWeightBandCountExpected int = 3     //expected number of weight bands verify against weightBandCount

type dlhWeightBand struct {
	lowerLbs int
	upperLbs int
}

var dlhWeightBands = []dlhWeightBand{
	{
		lowerLbs: 500,
		upperLbs: 4999,
	},
	{
		lowerLbs: 5000,
		upperLbs: 9999,
	},
	{
		lowerLbs: 10000,
		upperLbs: 999999,
	},
}

type dlhMilesRange struct {
	lower int
	upper int
}

var dlhMilesRanges = []dlhMilesRange{
	{
		lower: 0,
		upper: 250,
	},
	{
		lower: 251,
		upper: 500,
	},
	{
		lower: 501,
		upper: 1000,
	},
	{
		lower: 1001,
		upper: 1500,
	},
	{
		lower: 1501,
		upper: 2000,
	},
	{
		lower: 2001,
		upper: 2500,
	},
	{
		lower: 2501,
		upper: 3000,
	},
	{
		lower: 3001,
		upper: 3500,
	},
	{
		lower: 3501,
		upper: 4000,
	},
	{
		lower: 4001,
		upper: 999999,
	},
}

var dlhWeightBandNumCells = len(dlhMilesRanges)
