package channel

import(
	"ircd/m/user"
)
package channel

type Channel struct {
	users map[string]*user.User
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

