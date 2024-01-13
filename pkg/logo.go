package pkg

import ( 
	"strings"
)

func Logo()string{
	logo := strings.Join([]string{
		"          _nnnn_",
		"         dGGGGMMb",
		"        @p~qp~~qMb",
		"        M|@||@) M|",
		"        @,----.JM|",
		"       JS^\\__/  qKL",
		"      dZP        qKRb",
		"     dZP          qKKb",
		"    fZP            SMMb",
		"    HZM            MMMM",
		"    FqM            MMMM",
		"  __| \".        |\\dS\"qML",
		"  |    `.       | `' \\Zq",
		" _)      \\.___.,|     .'",
		" \\____   )MMMMMP|   .'",
		"      `-'       `--'",
		"[ENTER YOUR NAME]:",
	}, "\n")

	return logo
}

