package main

import "math/rand"

var PredefinedRoutes = [][][]float64{
	{
		{49.019780, 12.097520}, // Regensburg Hbf area
		{49.019140, 12.099950}, // towards Albertstraße
		{49.018150, 12.102350}, // towards D.-Martin-Luther-Str.
		{49.016990, 12.097980}, // near Bismarckplatz
	},
	{
		{49.020250, 12.096540}, // Bismarckplatz
		{49.019650, 12.097780}, // towards Arnulfsplatz
		{49.019120, 12.098920}, // Haidplatz area
		{49.018750, 12.099760}, // Rathaus / Altstadt
		{49.018240, 12.100980}, // Domplatz (Cathedral vicinity)
		{49.017930, 12.102120}, // towards Steinerne Brücke
		{49.018630, 12.103320}, // Steinerne Brücke (south end vicinity)
		{49.019610, 12.104320}, // Stadtamhof side (north end vicinity)
		{49.020250, 12.103210}, // Stadtamhof center-ish
		{49.020020, 12.101770}, // back towards bridge approach
		{49.019050, 12.101050}, // Dom / old town edge
	},
	{
		{49.017320, 12.095980}, // Stadtpark / Dörnbergpark edge
		{49.017980, 12.097260}, // towards Arnulfsplatz
		{49.018720, 12.098550}, // Haidplatz direction
		{49.018980, 12.099620}, // Rathaus area
		{49.019380, 12.100760}, // Dom area
		{49.018560, 12.101960}, // towards Danube / bridge approach
		{49.017750, 12.103020}, // near bridge ramp / riverside
		{49.016920, 12.102240}, // back into old town lanes
	},
	{
		{49.016920, 12.102240}, // reverse-ish loop
		{49.017750, 12.103020},
		{49.018560, 12.101960},
		{49.019380, 12.100760},
		{49.018980, 12.099620},
		{49.018720, 12.098550},
		{49.017980, 12.097260},
		{49.018240, 12.100980},
	},
}

func GenerateRandomPlate() string {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	plate := ""
	for i := 0; i < 3; i++ {
		plate += string(letters[rand.Intn(len(letters))])
	}

	return plate
}
