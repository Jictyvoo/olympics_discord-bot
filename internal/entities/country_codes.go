package entities

import (
	_ "embed"
)

/*
//go:embed ioc_codes.json
var iocCodesList []byte

func init() {
	iocCodes := make(map[string]string, 10)
	_ = json.Unmarshal(iocCodesList, &iocCodes)

	normalizeName := func(a string) string {
		return strings.ToLower(strings.TrimSpace(a))
	}
	var totalMatched uint64
iocLoop:
	for ioc, countryName := range iocCodes {
		countryName = normalizeName(countryName)
		for mapKey, info := range countriesData {
			if normalizeName(info.Name) == countryName {
				info.IOCCode = ioc
				countriesData[mapKey] = info
				totalMatched++
				continue iocLoop
			}
		}
		slog.Warn(fmt.Sprintf("Not found for %s - %s", ioc, countryName))
	}
	slog.Debug("Finish running IOC aggregation", slog.Uint64("totalMatched", totalMatched))
}
*/

type CountryInfo struct {
	Name       string
	CodeNum    string
	ISOCode    [2]string
	IOCCode    string
	Population uint64
	AreaKm2    float64
	GDPUSD     string
}

func GetCountryByCode(countryCode string) CountryInfo {
	found, ok := countriesData[countryCode]
	if ok {
		return found
	}

	for _, country := range countriesData {
		if country.IOCCode == countryCode || country.ISOCode[0] == countryCode ||
			country.ISOCode[1] == countryCode {
			return country
		}
	}

	return CountryInfo{IOCCode: countryCode, ISOCode: [2]string{countryCode}, Name: countryCode}
}

func GetCountryList() []CountryInfo {
	countries := make([]CountryInfo, 0, len(countriesData))
	for _, country := range countriesData {
		countries = append(countries, country)
	}

	return countries
}

var countriesData = map[string]CountryInfo{
	"ANT": {"Antigua and Barbuda", "1-268", [2]string{"AG", "ATG"}, "ANT", 86754, 443.000000, "1.22 Billion"},
	"SLO": {"Slovenia", "386", [2]string{"SI", "SVN"}, "SLO", 2007000, 20273.000000, "46.82 Billion"},
	"AFG": {"Afghanistan", "93", [2]string{"AF", "AFG"}, "AFG", 29121286, 647500.000000, "20.65 Billion"},
	"NAM": {"Namibia", "264", [2]string{"NA", "NAM"}, "NAM", 2128471, 825418.000000, "12.3 Billion"},
	"VAT": {"Vatican", "379", [2]string{"VA", "VAT"}, "VAT", 921, 0.000000, "NA"},
	"CRO": {"Croatia", "385", [2]string{"HR", "HRV"}, "CRO", 4491000, 56542.000000, "59.14 Billion"},
	"GRE": {"Greece", "30", [2]string{"GR", "GRC"}, "GRE", 11000000, 131940.000000, "243.3 Billion"},
	"JAM": {"Jamaica", "1-876", [2]string{"JM", "JAM"}, "JAM", 2847232, 10991.000000, "14.39 Billion"},
	"ALG": {"Algeria", "213", [2]string{"DZ", "DZA"}, "ALG", 34586184, 2381740.000000, "215.7 Billion"},
	"IRI": {"Iran", "98", [2]string{"IR", "IRN"}, "IRI", 76923300, 1648000.000000, "411.9 Billion"},
	"AHO": {"Netherlands Antilles", "599", [2]string{"AN", "ANT"}, "AHO", 300000, 960.000000, "NA"},
	"NCA": {"Nicaragua", "505", [2]string{"NI", "NIC"}, "NCA", 5995928, 129494.000000, "11.26 Billion"},
	"QAT": {"Qatar", "974", [2]string{"QA", "QAT"}, "QAT", 840926, 11437.000000, "213.1 Billion"},
	"AGU": {"Anguilla", "1-264", [2]string{"AI", "AIA"}, "AGU", 13254, 102.000000, "175.4 Million"},
	"ARU": {"Aruba", "297", [2]string{"AW", "ABW"}, "ARU", 71566, 193.000000, "2.516 Billion"},
	"BEL": {"Belgium", "32", [2]string{"BE", "BEL"}, "BEL", 10403000, 30510.000000, "507.4 Billion"},
	"CPV": {"Cape Verde", "238", [2]string{"CV", "CPV"}, "CPV", 508659, 4033.000000, "1.955 Billion"},
	"GUM": {"Guam", "1-671", [2]string{"GU", "GUM"}, "GUM", 159358, 549.000000, "4.6 Billion"},
	"KOS": {"Kosovo", "383", [2]string{"XK", "XKX"}, "KOS", 1800000, 10887.000000, "7.15 Billion"},
	"TUR": {"Turkey", "90", [2]string{"TR", "TUR"}, "TUR", 77804122, 780580.000000, "821.8 Billion"},
	"VAN": {"Vanuatu", "678", [2]string{"VU", "VUT"}, "VAN", 221552, 12200.000000, "828 Million"},
	"ANG": {"Angola", "244", [2]string{"AO", "AGO"}, "ANG", 13068161, 1246700.000000, "124 Billion"},
	"CUB": {"Cuba", "53", [2]string{"CU", "CUB"}, "CUB", 11423000, 110860.000000, "72.3 Billion"},
	"DMA": {"Dominica", "1-767", [2]string{"DM", "DMA"}, "DMA", 72813, 754.000000, "495 Million"},
	"SOL": {"Solomon Islands", "677", [2]string{"SB", "SLB"}, "SOL", 559198, 28450.000000, "1.099 Billion"},
	"SYR": {"Syria", "963", [2]string{"SY", "SYR"}, "SYR", 22198110, 185180.000000, "64.7 Billion"},
	"THA": {"Thailand", "66", [2]string{"TH", "THA"}, "THA", 67089500, 514000.000000, "400.9 Billion"},
	"BOL": {"Bolivia", "591", [2]string{"BO", "BOL"}, "BOL", 9947418, 1098580.000000, "30.79 Billion"},
	"COK": {"Cook Islands", "682", [2]string{"CK", "COK"}, "COK", 21388, 240.000000, "183.2 Million"},
	"GEQ": {"Equatorial Guinea", "240", [2]string{"GQ", "GNQ"}, "GEQ", 1014999, 28051.000000, "17.08 Billion"},
	"PLE": {"Palestine", "970", [2]string{"PS", "PSE"}, "PLE", 3800000, 5970.000000, "6.641 Billion"},
	"VEN": {"Venezuela", "58", [2]string{"VE", "VEN"}, "VEN", 27223228, 912050.000000, "367.5 Billion"},
	"SHP": {"Ascension Island", "", [2]string{"", ""}, "SHP", 0, 0.000000, ""},
	"BDI": {"Burundi", "257", [2]string{"BI", "BDI"}, "BDI", 9863117, 27830.000000, "2.676 Billion"},
	"FRA": {"France", "33", [2]string{"FR", "FRA"}, "FRA", 64768389, 547030.000000, "2.739 Trillion"},
	"TKL": {"Tokelau", "690", [2]string{"TK", "TKL"}, "TKL", 1466, 10.000000, "1.5 Million"},
	"UZB": {"Uzbekistan", "998", [2]string{"UZ", "UZB"}, "UZB", 27865738, 447400.000000, "55.18 Billion"},
	"YEM": {"Yemen", "967", [2]string{"YE", "YEM"}, "YEM", 23495361, 527968.000000, "43.89 Billion"},
	"NCD": {"New Caledonia", "687", [2]string{"NC", "NCL"}, "NCD", 216494, 19060.000000, "10.41 Billion"},
	"UAE": {"United Arab Emirates", "971", [2]string{"AE", "ARE"}, "UAE", 4975593, 82880.000000, "390 Billion"},
	"KIR": {"Kiribati", "686", [2]string{"KI", "KIR"}, "KIR", 92533, 811.000000, "173 Million"},
	"WAF": {"Wallis and Futuna", "681", [2]string{"WF", "WLF"}, "WAF", 16025, 274.000000, "188 Million"},
	"CKI": {"Cocos Islands", "61", [2]string{"CC", "CCK"}, "CKI", 628, 14.000000, ""},
	"KAZ": {"Kazakhstan", "7", [2]string{"KZ", "KAZ"}, "KAZ", 15340000, 2717300.000000, "224.9 Billion"},
	"USA": {"United States", "1", [2]string{"US", "USA"}, "USA", 310232863, 9629091.000000, "16.72 Trillion"},
	"BRA": {"Brazil", "55", [2]string{"BR", "BRA"}, "BRA", 201103330, 8511965.000000, "2.19 Trillion"},
	"SSD": {"South Sudan", "211", [2]string{"SS", "SSD"}, "SSD", 8260490, 644329.000000, "11.77 Billion"},
	"ALB": {"Albania", "355", [2]string{"AL", "ALB"}, "ALB", 2986952, 28748.000000, "12.8 Billion"},
	"AND": {"Andorra", "376", [2]string{"AD", "AND"}, "AND", 84000, 468.000000, "4.8 Billion"},
	"BAN": {"Bangladesh", "880", [2]string{"BD", "BGD"}, "BAN", 156118464, 144000.000000, "140.2 Billion"},
	"CZE": {"Czech Republic", "420", [2]string{"CZ", "CZE"}, "CZE", 10476000, 78866.000000, "194.8 Billion"},
	"CIV": {"Ivory Coast", "225", [2]string{"CI", "CIV"}, "CIV", 21058798, 322460.000000, "28.28 Billion"},
	"NIG": {"Niger", "227", [2]string{"NE", "NER"}, "NIG", 15878271, 1267000.000000, "7.304 Billion"},
	"PAR": {"Paraguay", "595", [2]string{"PY", "PRY"}, "PAR", 6375830, 406750.000000, "30.56 Billion"},
	"CGO": {"Republic of the Congo", "242", [2]string{"CG", "COG"}, "CGO", 3039126, 342000.000000, "14.25 Billion"},
	"RNN": {"La Reunion", "262", [2]string{"RE", "REU"}, "RNN", 776948, 2517.000000, "NA"},
	"GBR": {"United Kingdom", "44", [2]string{"GB", "GBR"}, "GBR", 62348447, 244820.000000, "2.435 Trillion"},
	"FIN": {"Finland", "358", [2]string{"FI", "FIN"}, "FIN", 5244000, 337030.000000, "259.6 Billion"},
	"STP": {"Sao Tome and Principe", "239", [2]string{"ST", "STP"}, "STP", 175808, 1001.000000, "311 Million"},
	"ESP": {"Spain", "34", [2]string{"ES", "ESP"}, "ESP", 46505963, 504782.000000, "1.356 Trillion"},
	"CRC": {"Costa Rica", "506", [2]string{"CR", "CRI"}, "CRC", 4516220, 51100.000000, "48.51 Billion"},
	"DOM": {"Dominican Republic", "1-809, 1-829, 1-849", [2]string{"DO", "DOM"}, "DOM", 9823821, 48730.000000, "59.27 Billion"},
	"GRL": {"Greenland", "299", [2]string{"GL", "GRL"}, "GRL", 56375, 2166086.000000, "2.16 Billion"},
	"RUS": {"Russia", "7", [2]string{"RU", "RUS"}, "RUS", 140702000, 17100000.000000, "2.113 Trillion"},
	"FLI": {"Falkland Islands", "500", [2]string{"FK", "FLK"}, "FLI", 2638, 12173.000000, "164.5 Million"},
	"MYA": {"Myanmar", "95", [2]string{"MM", "MMR"}, "MYA", 53414374, 678500.000000, "59.04 Billion"},
	"ARG": {"Argentina", "54", [2]string{"AR", "ARG"}, "ARG", 41343201, 2766890.000000, "484.6 Billion"},
	"BIH": {"Bosnia and Herzegovina", "387", [2]string{"BA", "BIH"}, "BIH", 4590000, 51129.000000, "18.87 Billion"},
	"GUI": {"Guinea", "224", [2]string{"GN", "GIN"}, "GUI", 10324025, 245857.000000, "6.544 Billion"},
	"SWZ": {"Swaziland", "268", [2]string{"SZ", "SWZ"}, "SWZ", 1354051, 17363.000000, "3.807 Billion"},
	"TAN": {"Tanzania", "255", [2]string{"TZ", "TZA"}, "TAN", 41892895, 945087.000000, "31.94 Billion"},
	"TRI": {"Trinidad and Tobago", "1-868", [2]string{"TT", "TTO"}, "TRI", 1228691, 5128.000000, "43.69 Billion"},
	"UKR": {"Ukraine", "380", [2]string{"UA", "UKR"}, "UKR", 45415596, 603700.000000, "175.5 Billion"},
	"IND": {"India", "91", [2]string{"IN", "IND"}, "IND", 1173108018, 3287590.000000, "1.67 Trillion"},
	"TPE": {"Taiwan", "886", [2]string{"TW", "TWN"}, "TPE", 22894384, 35980.000000, "489.2 Billion"},
	"BER": {"Bermuda", "1-441", [2]string{"BM", "BMU"}, "BER", 65365, 53.000000, "5.6 Billion"},
	"COM": {"Comoros", "269", [2]string{"KM", "COM"}, "COM", 773407, 2170.000000, "658 Million"},
	"INA": {"Indonesia", "62", [2]string{"ID", "IDN"}, "INA", 242968342, 1919440.000000, "867.5 Billion"},
	"IOM": {"Isle of Man", "44-1624", [2]string{"IM", "IMN"}, "IOM", 75049, 572.000000, "4.076 Billion"},
	"ISR": {"Israel", "972", [2]string{"IL", "ISR"}, "ISR", 7353985, 20770.000000, "272.7 Billion"},
	"JPN": {"Japan", "81", [2]string{"JP", "JPN"}, "JPN", 127288000, 377835.000000, "5.007 Trillion"},
	"LBR": {"Liberia", "231", [2]string{"LR", "LBR"}, "LBR", 3685076, 111369.000000, "1.977 Billion"},
	"MDV": {"Maldives", "960", [2]string{"MV", "MDV"}, "MDV", 395650, 300.000000, "2.27 Billion"},
	"FIJ": {"Fiji", "679", [2]string{"FJ", "FJI"}, "FIJ", 875983, 18274.000000, "4.218 Billion"},
	"KGZ": {"Kyrgyzstan", "996", [2]string{"KG", "KGZ"}, "KGZ", 5776500, 199900.000000, "7.234 Billion"},
	"MDA": {"Moldova", "373", [2]string{"MD", "MDA"}, "MDA", 4324000, 33843.000000, "7.932 Billion"},
	"MON": {"Monaco", "377", [2]string{"MC", "MCO"}, "MON", 32965, 2.000000, "5.748 Billion"},
	"PAN": {"Panama", "507", [2]string{"PA", "PAN"}, "PAN", 3410676, 78200.000000, "40.62 Billion"},
	"CAM": {"Cambodia", "855", [2]string{"KH", "KHM"}, "CAM", 14453680, 181040.000000, "15.64 Billion"},
	"EST": {"Estonia", "372", [2]string{"EE", "EST"}, "EST", 1291170, 45227.000000, "24.28 Billion"},
	"ITA": {"Italy", "39", [2]string{"IT", "ITA"}, "ITA", 60340328, 301230.000000, "2.068 Trillion"},
	"GAB": {"Gabon", "241", [2]string{"GA", "GAB"}, "GAB", 1545255, 267668.000000, "19.97 Billion"},
	"KOR": {"South Korea", "82", [2]string{"KR", "KOR"}, "KOR", 48422644, 98480.000000, "1.198 Trillion"},
	"ETH": {"Ethiopia", "251", [2]string{"ET", "ETH"}, "ETH", 88013491, 1127127.000000, "47.34 Billion"},
	"ZAM": {"Zambia", "260", [2]string{"ZM", "ZMB"}, "ZAM", 13460305, 752614.000000, "22.24 Billion"},
	"COL": {"Colombia", "57", [2]string{"CO", "COL"}, "COL", 47790000, 1138910.000000, "369.2 Billion"},
	"GIC": {"Gibraltar", "350", [2]string{"GI", "GIB"}, "GIC", 27884, 6.000000, "1.106 Billion"},
	"TUN": {"Tunisia", "216", [2]string{"TN", "TUN"}, "TUN", 10589025, 163610.000000, "48.38 Billion"},
	"BAR": {"Barbados", "1-246", [2]string{"BB", "BRB"}, "BAR", 285653, 431.000000, "4.262 Billion"},
	"GAM": {"Gambia", "220", [2]string{"GM", "GMB"}, "GAM", 1593256, 11295.000000, "896 Million"},
	"SEN": {"Senegal", "221", [2]string{"SN", "SEN"}, "SEN", 12323252, 196722.000000, "15.36 Billion"},
	"SVB": {"Svalbard", "47", [2]string{"SJ", "SJM"}, "SVB", 2550, 62049.000000, "NA"},
	"SUI": {"Switzerland", "41", [2]string{"CH", "CHE"}, "SUI", 7581000, 41290.000000, "646.2 Billion"},
	"IRQ": {"Iraq", "964", [2]string{"IQ", "IRQ"}, "IRQ", 29671605, 437072.000000, "221.8 Billion"},
	"CAN": {"Canada", "1", [2]string{"CA", "CAN"}, "CAN", 33679000, 9984670.000000, "1.825 Trillion"},
	"LIB": {"Lebanon", "961", [2]string{"LB", "LBN"}, "LIB", 4125247, 10400.000000, "43.49 Billion"},
	"SVK": {"Slovakia", "421", [2]string{"SK", "SVK"}, "SVK", 5456362, 48845.000000, "97.82 Billion"},
	"CMR": {"Cameroon", "237", [2]string{"CM", "CMR"}, "CMR", 19294149, 475440.000000, "27.88 Billion"},
	"CHA": {"Chad", "235", [2]string{"TD", "TCD"}, "CHA", 10543464, 1284000.000000, "13.59 Billion"},
	"GBS": {"Guinea-Bissau", "245", [2]string{"GW", "GNB"}, "GBS", 1565126, 36125.000000, "880 Million"},
	"SMT": {"Sint Maarten", "1-721", [2]string{"SX", "SXM"}, "SMT", 37429, 34.000000, "794.7 Million"},
	"BHU": {"Bhutan", "975", [2]string{"BT", "BTN"}, "BHU", 699847, 47000.000000, "2.133 Billion"},
	"BUL": {"Bulgaria", "359", [2]string{"BG", "BGR"}, "BUL", 7148785, 110910.000000, "53.7 Billion"},
	"CAF": {"Central African Republic", "236", [2]string{"CF", "CAF"}, "CAF", 4844927, 622984.000000, "2.05 Billion"},
	"SRB": {"Serbia", "381", [2]string{"RS", "SRB"}, "SRB", 7344847, 88361.000000, "43.68 Billion"},
	"UGA": {"Uganda", "256", [2]string{"UG", "UGA"}, "UGA", 33398682, 236040.000000, "22.6 Billion"},
	"MAC": {"Macau", "853", [2]string{"MO", "MAC"}, "MAC", 449198, 254.000000, "51.68 Billion"},
	"PRK": {"North Korea", "850", [2]string{"KP", "PRK"}, "PRK", 22912177, 120540.000000, "28 Billion"},
	"HEL": {"Saint Helena", "290", [2]string{"SH", "SHN"}, "HEL", 7460, 410.000000, "18 Million"},
	"GRN": {"Grenada", "1-473", [2]string{"GD", "GRD"}, "GRN", 107818, 344.000000, "811 Million"},
	"LUX": {"Luxembourg", "352", [2]string{"LU", "LUX"}, "LUX", 497538, 2586.000000, "60.54 Billion"},
	"POL": {"Poland", "48", [2]string{"PL", "POL"}, "POL", 38500000, 312685.000000, "513.9 Billion"},
	"SRI": {"Sri Lanka", "94", [2]string{"LK", "LKA"}, "SRI", 21513990, 65610.000000, "65.12 Billion"},
	"TOG": {"Togo", "228", [2]string{"TG", "TGO"}, "TOG", 6587239, 56785.000000, "4.299 Billion"},
	"AUS": {"Australia", "61", [2]string{"AU", "AUS"}, "AUS", 21515754, 7686850.000000, "1.488 Trillion"},
	"CUR": {"Curaçao", "599", [2]string{"CW", "CUW"}, "CUR", 141766, 444.000000, "5.6 Billion"},
	"MRI": {"Mauritius", "230", [2]string{"MU", "MUS"}, "MRI", 1294104, 2040.000000, "11.9 Billion"},
	"RWA": {"Rwanda", "250", [2]string{"RW", "RWA"}, "RWA", 11055976, 26338.000000, "7.7 Billion"},
	"TKM": {"Turkmenistan", "993", [2]string{"TM", "TKM"}, "TKM", 4940916, 488100.000000, "40.56 Billion"},
	"CYP": {"Cyprus", "357", [2]string{"CY", "CYP"}, "CYP", 1102677, 9250.000000, "21.78 Billion"},
	"ASA": {"American Samoa", "1-684", [2]string{"AS", "ASM"}, "ASA", 57881, 199.000000, "462.2 Million"},
	"MHL": {"Marshall Islands", "692", [2]string{"MH", "MHL"}, "MHL", 65859, 181.000000, "193 Million"},
	"AZE": {"Azerbaijan", "994", [2]string{"AZ", "AZE"}, "AZE", 8303512, 86600.000000, "76.01 Billion"},
	"BRN": {"Bahrain", "973", [2]string{"BH", "BHR"}, "BRN", 738004, 665.000000, "28.36 Billion"},
	"BRU": {"Brunei", "673", [2]string{"BN", "BRN"}, "BRU", 395027, 5770.000000, "16.56 Billion"},
	"MGL": {"Mongolia", "976", [2]string{"MN", "MNG"}, "MGL", 3086918, 1565000.000000, "11.14 Billion"},
	"CAY": {"Cayman Islands", "1-345", [2]string{"KY", "CYM"}, "CAY", 44270, 262.000000, "2.25 Billion"},
	"GER": {"Germany", "49", [2]string{"DE", "DEU"}, "GER", 81802257, 357021.000000, "3.593 Trillion"},
	"JCI": {"Jersey", "44-1534", [2]string{"JE", "JEY"}, "JCI", 90812, 116.000000, "5.1 Billion"},
	"LAO": {"Laos", "856", [2]string{"LA", "LAO"}, "LAO", 6368162, 236800.000000, "10.1 Billion"},
	"FSM": {"Micronesia", "691", [2]string{"FM", "FSM"}, "FSM", 107708, 702.000000, "339 Million"},
	"MTS": {"Montserrat", "1-664", [2]string{"MS", "MSR"}, "MTS", 9341, 102.000000, "54.72 Million"},
	"OMA": {"Oman", "968", [2]string{"OM", "OMN"}, "OMA", 2967717, 212460.000000, "81.95 Billion"},
	"SOM": {"Somalia", "252", [2]string{"SO", "SOM"}, "SOM", 10112453, 637657.000000, "NA"},
	"ISV": {"British Virgin Islands", "1-284", [2]string{"VG", "VGB"}, "ISV", 21730, 153.000000, "1.095 Billion"},
	"FPN": {"French Polynesia", "689", [2]string{"PF", "PYF"}, "FPN", 270485, 4167.000000, "5.65 Billion"},
	"HAI": {"Haiti", "509", [2]string{"HT", "HTI"}, "HAI", 9648924, 27750.000000, "8.287 Billion"},
	"SKN": {"Saint Kitts and Nevis", "1-869", [2]string{"KN", "KNA"}, "SKN", 51134, 261.000000, "767 Million"},
	"SEY": {"Seychelles", "248", [2]string{"SC", "SYC"}, "SEY", 88340, 455.000000, "1.271 Billion"},
	"XMI": {"Christmas Island", "61", [2]string{"CX", "CXR"}, "XMI", 1500, 135.000000, ""},
	"GUY": {"Guyana", "592", [2]string{"GY", "GUY"}, "GUY", 748486, 214970.000000, "3.02 Billion"},
	"NRU": {"Nauru", "674", [2]string{"NR", "NRU"}, "NRU", 10065, 21.000000, "NA"},
	"NGR": {"Nigeria", "234", [2]string{"NG", "NGA"}, "NGR", 154000000, 923768.000000, "502 Billion"},
	"PAK": {"Pakistan", "92", [2]string{"PK", "PAK"}, "PAK", 184404791, 803940.000000, "236.5 Billion"},
	"MAS": {"Malaysia", "60", [2]string{"MY", "MYS"}, "MAS", 28274729, 329750.000000, "312.4 Billion"},
	"PER": {"Peru", "51", [2]string{"PE", "PER"}, "PER", 29907003, 1285220.000000, "210.3 Billion"},
	"ZIM": {"Zimbabwe", "263", [2]string{"ZW", "ZWE"}, "ZIM", 11651858, 390580.000000, "10.48 Billion"},
	"MOZ": {"Mozambique", "258", [2]string{"MZ", "MOZ"}, "MOZ", 22061451, 801590.000000, "14.67 Billion"},
	"RSA": {"South Africa", "27", [2]string{"ZA", "ZAF"}, "RSA", 49000000, 1219912.000000, "353.9 Billion"},
	"BAH": {"Bahamas", "1-242", [2]string{"BS", "BHS"}, "BAH", 301790, 13940.000000, "8.373 Billion"},
	"IRL": {"Ireland", "353", [2]string{"IE", "IRL"}, "IRL", 4622917, 70280.000000, "220.9 Billion"},
	"MLT": {"Malta", "356", [2]string{"MT", "MLT"}, "MLT", 403000, 316.000000, "9.541 Billion"},
	"POR": {"Portugal", "351", [2]string{"PT", "PRT"}, "POR", 10676000, 92391.000000, "220.3 Billion"},
	"TJK": {"Tajikistan", "992", [2]string{"TJ", "TJK"}, "TJK", 7487489, 143100.000000, "8.513 Billion"},
	"DJI": {"Djibouti", "253", [2]string{"DJ", "DJI"}, "DJI", 740528, 23000.000000, "1.459 Billion"},
	"JOR": {"Jordan", "962", [2]string{"JO", "JOR"}, "JOR", 6407085, 92300.000000, "34.08 Billion"},
	"BRT": {"Saint Barthelemy", "590", [2]string{"BL", "BLM"}, "BRT", 8450, 21.000000, "255 Million"},
	"LCA": {"Saint Lucia", "1-758", [2]string{"LC", "LCA"}, "LCA", 160922, 616.000000, "1.377 Billion"},
	"GHA": {"Ghana", "233", [2]string{"GH", "GHA"}, "GHA", 24339838, 239460.000000, "45.55 Billion"},
	"MRT": {"Saint Martin", "590", [2]string{"MF", "MAF"}, "MRT", 35925, 53.000000, "561.5 Million"},
	"SAM": {"Samoa", "685", [2]string{"WS", "WSM"}, "SAM", 192001, 2944.000000, "705 Million"},
	"BEN": {"Benin", "229", [2]string{"BJ", "BEN"}, "BEN", 9056010, 112620.000000, "8.359 Billion"},
	"MAR": {"Morocco", "212", [2]string{"MA", "MAR"}, "MAR", 31627428, 446550.000000, "104.8 Billion"},
	"ROU": {"Romania", "40", [2]string{"RO", "ROU"}, "ROU", 21959278, 237500.000000, "188.9 Billion"},
	"ARM": {"Armenia", "374", [2]string{"AM", "ARM"}, "ARM", 2968000, 29800.000000, "10.44 Billion"},
	"LIE": {"Liechtenstein", "423", [2]string{"LI", "LIE"}, "LIE", 35000, 160.000000, "5.113 Billion"},
	"NZL": {"New Zealand", "64", [2]string{"NZ", "NZL"}, "NZL", 4252277, 268680.000000, "181.1 Billion"},
	"PUR": {"Puerto Rico", "1-787, 1-939", [2]string{"PR", "PRI"}, "PUR", 3916632, 9104.000000, "93.52 Billion"},
	"TUV": {"Tuvalu", "688", [2]string{"TV", "TUV"}, "TUV", 10472, 26.000000, "38 Million"},
	"BOT": {"Botswana", "267", [2]string{"BW", "BWA"}, "BOT", 2029307, 600370.000000, "15.53 Billion"},
	"ECU": {"Ecuador", "593", [2]string{"EC", "ECU"}, "ECU", 14790608, 283560.000000, "91.41 Billion"},
	"EGY": {"Egypt", "20", [2]string{"EG", "EGY"}, "EGY", 80471869, 1001450.000000, "262 Billion"},
	"LES": {"Lesotho", "266", [2]string{"LS", "LSO"}, "LES", 1919552, 30355.000000, "2.457 Billion"},
	"VIN": {"Saint Vincent and the Grenadines", "1-784", [2]string{"VC", "VCT"}, "VIN", 104217, 389.000000, "742 Million"},
	"SWE": {"Sweden", "46", [2]string{"SE", "SWE"}, "SWE", 9828655, 449964.000000, "552 Billion"},
	"KEN": {"Kenya", "254", [2]string{"KE", "KEN"}, "KEN", 40046566, 582650.000000, "45.31 Billion"},
	"PHI": {"Philippines", "63", [2]string{"PH", "PHL"}, "PHI", 99900177, 300000.000000, "272.2 Billion"},
	"BIZ": {"Belize", "501", [2]string{"BZ", "BLZ"}, "BIZ", 314522, 22966.000000, "1.637 Billion"},
	"MYT": {"Mayotte", "262", [2]string{"YT", "MYT"}, "MYT", 159042, 374.000000, "2.254 Billion"},
	"NEP": {"Nepal", "977", [2]string{"NP", "NPL"}, "NEP", 28951852, 140800.000000, "19.64 Billion"},
	"NIU": {"Niue", "683", [2]string{"NU", "NIU"}, "NIU", 2166, 260.000000, "10.01 Million"},
	"TGA": {"Tonga", "676", [2]string{"TO", "TON"}, "TGA", 122580, 748.000000, "477 Million"},
	"ESA": {"El Salvador", "503", [2]string{"SV", "SLV"}, "ESA", 6052064, 21041.000000, "24.67 Billion"},
	"HUN": {"Hungary", "36", [2]string{"HU", "HUN"}, "HUN", 9982000, 93030.000000, "130.6 Billion"},
	"LTU": {"Lithuania", "370", [2]string{"LT", "LTU"}, "LTU", 2944459, 65200.000000, "46.71 Billion"},
	"KSA": {"Saudi Arabia", "966", [2]string{"SA", "SAU"}, "KSA", 25731776, 1960582.000000, "718.5 Billion"},
	"MEX": {"Mexico", "52", [2]string{"MX", "MEX"}, "MEX", 112468855, 1972550.000000, "1.327 Trillion"},
	"SMR": {"San Marino", "378", [2]string{"SM", "SMR"}, "SMR", 31477, 61.000000, "1.866 Billion"},
	"ERI": {"Eritrea", "291", [2]string{"ER", "ERI"}, "ERI", 5792984, 121320.000000, "3.438 Billion"},
	"LAT": {"Latvia", "371", [2]string{"LV", "LVA"}, "LAT", 2217969, 64589.000000, "30.38 Billion"},
	"SPM": {"Saint Pierre and Miquelon", "508", [2]string{"PM", "SPM"}, "SPM", 7012, 242.000000, "215.3 Million"},
	"BUR": {"Burkina Faso", "226", [2]string{"BF", "BFA"}, "BUR", 16241811, 274200.000000, "12.13 Billion"},
	"COD": {"Democratic Republic of the Congo", "243", [2]string{"CD", "COD"}, "COD", 70916439, 2345410.000000, "18.56 Billion"},
	"KUW": {"Kuwait", "965", [2]string{"KW", "KWT"}, "KUW", 2789132, 17820.000000, "165.3 Billion"},
	"AUT": {"Austria", "43", [2]string{"AT", "AUT"}, "AUT", 8205000, 83858.000000, "417.9 Billion"},
	"TLS": {"East Timor", "670", [2]string{"TL", "TLS"}, "TLS", 1154625, 15007.000000, "6.129 Billion"},
	"HON": {"Honduras", "504", [2]string{"HN", "HND"}, "HON", 7989415, 112090.000000, "18.88 Billion"},
	"NOR": {"Norway", "47", [2]string{"NO", "NOR"}, "NOR", 5009150, 324220.000000, "515.8 Billion"},
	"PNG": {"Papua New Guinea", "675", [2]string{"PG", "PNG"}, "PNG", 6064515, 462840.000000, "16.1 Billion"},
	"SUR": {"Suriname", "597", [2]string{"SR", "SUR"}, "SUR", 492829, 163820.000000, "5.009 Billion"},
	"CHI": {"Chile", "56", [2]string{"CL", "CHL"}, "CHI", 16746491, 756950.000000, "281.7 Billion"},
	"CHN": {"China", "86", [2]string{"CN", "CHN"}, "CHN", 1330044000, 9596960.000000, "9.33 Trillion"},
	"GUA": {"Guatemala", "502", [2]string{"GT", "GTM"}, "GUA", 13550440, 108890.000000, "53.9 Billion"},
	"GCI": {"Guernsey", "44-1481", [2]string{"GG", "GGY"}, "GCI", 65228, 78.000000, "2.742 Billion"},
	"HKG": {"Hong Kong", "852", [2]string{"HK", "HKG"}, "HKG", 6898686, 1092.000000, "272.1 Billion"},
	"MNE": {"Montenegro", "382", [2]string{"ME", "MNE"}, "MNE", 666730, 14026.000000, "4.518 Billion"},
	"NED": {"Netherlands", "31", [2]string{"NL", "NLD"}, "NED", 16645000, 41526.000000, "722.3 Billion"},
	"PLW": {"Palau", "680", [2]string{"PW", "PLW"}, "PLW", 19907, 459.000000, "221 Million"},
	"SIN": {"Singapore", "65", [2]string{"SG", "SGP"}, "SIN", 4701069, 692.000000, "295.7 Billion"},
	"ISL": {"Iceland", "354", [2]string{"IS", "ISL"}, "ISL", 308910, 103000.000000, "14 Billion"},
	"VIE": {"Vietnam", "84", [2]string{"VN", "VNM"}, "VIE", 89571130, 329560.000000, "170 Billion"},
	"BLR": {"Belarus", "375", [2]string{"BY", "BLR"}, "BLR", 9685000, 207600.000000, "69.24 Billion"},
	"DEN": {"Denmark", "45", [2]string{"DK", "DNK"}, "DEN", 5484000, 43094.000000, "324.3 Billion"},
	"LBA": {"Libya", "218", [2]string{"LY", "LBY"}, "LBA", 6461454, 1759540.000000, "70.92 Billion"},
	"SLE": {"Sierra Leone", "232", [2]string{"SL", "SLE"}, "SLE", 5245695, 71740.000000, "4.607 Billion"},
	"SUD": {"Sudan", "249", [2]string{"SD", "SDN"}, "SUD", 35000000, 1861484.000000, "52.5 Billion"},
	"FAI": {"Faroe Islands", "298", [2]string{"FO", "FRO"}, "FAI", 48228, 1399.000000, "2.32 Billion"},
	"GEO": {"Georgia", "995", [2]string{"GE", "GEO"}, "GEO", 4630000, 69700.000000, "15.95 Billion"},
	"MAD": {"Madagascar", "261", [2]string{"MG", "MDG"}, "MAD", 21281844, 587040.000000, "10.53 Billion"},
	"MAW": {"Malawi", "265", [2]string{"MW", "MWI"}, "MAW", 15447500, 118480.000000, "3.683 Billion"},
	"MKD": {"North Macedonia", "389", [2]string{"MK", "MKD"}, "MKD", 2062294, 25333.000000, "10.65 Billion"},
	"NMI": {"Northern Mariana Islands", "1-670", [2]string{"MP", "MNP"}, "NMI", 53883, 477.000000, "733 Million"},
	"URU": {"Uruguay", "598", [2]string{"UY", "URY"}, "URU", 3477000, 176220.000000, "57.11 Billion"},
	"MLI": {"Mali", "223", [2]string{"ML", "MLI"}, "MLI", 13796354, 1240192.000000, "11.37 Billion"},
	"MTN": {"Mauritania", "222", [2]string{"MR", "MRT"}, "MTN", 3205060, 1030700.000000, "4.183 Billion"},
}
