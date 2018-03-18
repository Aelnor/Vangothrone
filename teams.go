package main

type Team struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	FunName string `json:"funName"`
}

var Teams = []Team{
	Team{
		Id:      1,
		Name:    "London Spitfire",
		Code:    "LDN",
		FunName: "Landan Spit",
	},
	Team{
		Id:      2,
		Name:    "Boston Uprising",
		Code:    "BOS",
		FunName: "Boston Downfalling",
	},
	Team{
		Id:      3,
		Name:    "Seoul Dynasty",
		Code:    "SEO",
		FunName: "Seoul Diveasty",
	},
	Team{
		Id:      4,
		Name:    "Houston Outlaws",
		Code:    "HOU",
		FunName: "J LUL KE",
	},
	Team{
		Id:      5,
		Name:    "New York Excelsior",
		Code:    "NYE",
		FunName: "New York Dablords",
	},
	Team{
		Id:      6,
		Name:    "Los Angeles Gladiators",
		Code:    "GLA",
		FunName: "Los Angeles Fissure",
	},
	Team{
		Id:      7,
		Name:    "Los Angeles Valiant",
		Code:    "LAV",
		FunName: "Los Angeles Valiants",
	},
	Team{
		Id:      8,
		Name:    "Shanghai Dragons",
		Code:    "SHD",
		FunName: "Shanghai Kappa",
	},
	Team{
		Id:      9,
		Name:    "Dallas Fuel",
		Code:    "DAL",
		FunName: "Dallas Cocksuckers",
	},
	Team{
		Id:      10,
		Name:    "San Francisco Shock",
		Code:    "SFS",
		FunName: "San Francisco Cock",
	},
	Team{
		Id:      11,
		Name:    "Florida Mayhem",
		Code:    "FLA",
		FunName: "Phlorida Memehem",
	},
	Team{
		Id:      12,
		Name:    "Philadelphia Fusion",
		Code:    "PHI",
		FunName: "Philadelphia Phusion",
	},
}
