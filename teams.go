package main

type Team struct {
	id      int    `json:"id"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	FunName string `json:"funName"`
}

var Teams = []Team{
	Team{
		id:      1,
		Name:    "London Spitfire",
		Code:    "LDN",
		FunName: "Landan Spit",
	},
	Team{
		id:      2,
		Name:    "Boston Uprising",
		Code:    "BOS",
		FunName: "Boston Downfalling",
	},
	Team{
		id:      3,
		Name:    "Seoul Dynasty",
		Code:    "SEO",
		FunName: "Seoul Diveasty",
	},
	Team{
		id:      4,
		Name:    "Houston Outlaws",
		Code:    "HOU",
		FunName: "J LUL KE",
	},
	Team{
		id:      5,
		Name:    "New York Excelsior",
		Code:    "NYE",
		FunName: "New York Dablords",
	},
	Team{
		id:      6,
		Name:    "Los Angeles Gladiators",
		Code:    "GLA",
		FunName: "Los Angeles Fissure",
	},
	Team{
		id:      7,
		Name:    "Los Angeles Valiant",
		Code:    "VAL",
		FunName: "Los Angeles Valiants",
	},
	Team{
		id:      8,
		Name:    "Shanghai Dragons",
		Code:    "SHD",
		FunName: "Shanghai Kappa",
	},
	Team{
		id:      9,
		Name:    "Dallas Fuel",
		Code:    "DAL",
		FunName: "Dallas Cocksuckers",
	},
	Team{
		id:      10,
		Name:    "San Francisco Shock",
		Code:    "SFS",
		FunName: "San Francisco Cock",
	},
	Team{
		id:      11,
		Name:    "Florida Mayhem",
		Code:    "FLA",
		FunName: "Phlorida Memehem",
	},
	Team{
		id:      12,
		Name:    "Philadelphia Fusion",
		Code:    "PHI",
		FunName: "Philadelphia Phusion",
	},
}
