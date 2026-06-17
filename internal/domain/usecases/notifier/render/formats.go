package render

import (
	"fmt"
	"strings"
	"time"

	"github.com/jictyvoo/olhojogo/pkg/strutil"
)

const (
	gymnasticsIcon = ":person_doing_cartwheel:"
	athleticsIcon  = ":athletic_shoe:"
	soccerIcon     = ":soccer:"
)

var disciplineIconPerCode = map[string]string{
	"ATH": athleticsIcon,
	"BDM": ":badminton:",
	"BKB": ":basketball:",
	"BK3": ":basketball:",
	"BOX": ":boxing_glove:",
	"BKG": ":dancer:",
	"CSL": ":canoe:",
	"CSP": ":canoe:",
	"BMF": ":bike:",
	"BMX": ":person_biking:",
	"CRD": ":woman_mountain_biking:",
	"CTR": ":man_mountain_biking:",
	"MTB": ":person_mountain_biking:",
	"CLB": ":person_climbing:",
	"FEN": ":person_fencing:",
	"FBL": soccerIcon,
	"GAR": gymnasticsIcon,
	"GTR": gymnasticsIcon,
	"GRY": gymnasticsIcon,
	"GLF": ":golf:",
	"HBL": ":person_playing_handball:",
	"EQU": ":horse_racing:",
	"HOC": ":field_hockey:",
	"JUD": ":martial_arts_uniform:",
	"WLF": ":person_lifting_weights:",
	"WRE": ":men_wrestling:",
	"OWS": ":one_piece_swimsuit:",
	"SWA": ":woman_swimming: ",
	"SWM": ":person_swimming:",
	"MPN": ":person_fencing:",
	"WPO": ":person_playing_water_polo:",
	"ROW": ":person_rowing_boat:",
	"RU7": ":rugby_football:",
	"DIV": ":person_juggling:",
	"SKB": ":skateboard:",
	"SRF": ":person_surfing:",
	"TKW": ":martial_arts_uniform:",
	"TEN": ":tennis:",
	"TTE": ":ping_pong:",
	"ARC": ":bow_and_arrow:",
	"SHO": ":bow_and_arrow:",
	"TRI": ":triangular_ruler:",
	"SAL": ":sailboat:",
	"VVO": ":volleyball:",
	"VBV": ":volleyball: :beach:",
}

// disciplineIconPerName maps a diacritic-folded discipline name to its icon, for
// providers whose competition code is not an Olympics discipline code.
var disciplineIconPerName = map[string]string{
	"football": soccerIcon,
	"futebol":  soccerIcon,
	"futbol":   soccerIcon,
	"fußball":  soccerIcon,
	"calcio":   soccerIcon,
}

func DisciplineIcon(code string) string {
	return disciplineIconPerCode[code]
}

func DisciplineIconByName(name string) string {
	return disciplineIconPerName[strings.ToLower(strutil.FoldDiacritics(name))]
}

// DiscordTimestamp renders a time as a Discord relative-timestamp markdown token.
func DiscordTimestamp(timestamp time.Time) string {
	return fmt.Sprintf("<t:%d:R>", timestamp.Unix())
}
