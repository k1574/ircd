package channel

import(
	"ircd/m/user"
)

type Channel struct {
	Users map[string]*user.User
}


var(
	// Names delimiter.
	NamDel = ","
	// Channel name prefixes.
	NamPre = []string{"#", "&"}
	// Name restricted characters.
	NamRes = []string{NamDel, string([]byte{7})}
	MaxNamLen = 200
)

