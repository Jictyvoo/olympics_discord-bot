package utils

const gymnasticsIcon = ":person_doing_cartwheel:"

var disciplineIconPerCode = map[string]string{
	"ATH": ":athletic_shoe:",
	"BDM": ":badminton:", "BKB": ":basketball:", "BK3": ":basketball:", "BOX": ":boxing_glove:",
	"BKG": ":dancer:", "CSL": ":canoe:", "CSP": ":canoe:", "BMF": ":bike:", "BMX": ":person_biking:",
	"CRD": ":woman_mountain_biking:", "CTR": ":man_mountain_biking:", "MTB": ":person_mountain_biking:", "CLB": ":person_climbing:", "FEN": ":person_fencing:",
	"FBL": ":soccer:", "GAR": gymnasticsIcon, "GTR": gymnasticsIcon, "GRY": gymnasticsIcon, "GLF": ":golf:",
	"HBL": ":person_playing_handball:", "EQU": ":horse_racing:", "HOC": ":field_hockey:", "JUD": ":martial_arts_uniform:", "WLF": ":person_lifting_weights:",
	"WRE": ":men_wrestling:", "OWS": ":one_piece_swimsuit:", "SWA": ":woman_swimming: ", "SWM": ":person_swimming:", "MPN": ":person_fencing:",
	"WPO": ":person_playing_water_polo:", "ROW": ":person_rowing_boat:", "RU7": ":rugby_football:", "DIV": ":person_juggling:", "SKB": ":skateboard:",
	"SRF": ":person_surfing:", "TKW": ":martial_arts_uniform:", "TEN": ":tennis:", "TTE": ":ping_pong:", "ARC": ":bow_and_arrow:",
	"SHO": ":bow_and_arrow:", "TRI": ":triangular_ruler:", "SAL": ":sailboat:", "VVO": ":volleyball:", "VBV": ":volleyball: :beach:",
}

func DisciplineIcon(code string) string {
	return disciplineIconPerCode[code]
}
