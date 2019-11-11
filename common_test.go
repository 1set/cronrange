package cronrange

var (
	emptyString          = ""
	exprEveryMin         = "* * * * *"
	exprEveryXmasMorning = "0 8 25 12 *"
	exprEveryNewYear     = "0 0 1 1 *"
	timeZoneBangkok      = "Asia/Bangkok"
	timeZoneNewYork      = "America/New_York"
)

var (
	crEvery1Min, _ = New(exprEveryMin, emptyString, 1)
)
