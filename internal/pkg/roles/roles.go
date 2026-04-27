package roles

type Role int

const (
	Unknown Role = iota
	Blocked
	User
	Operator
	Admin
)

var toString = map[Role]string{
	0: "unknown",
	1: "blocked",
	2: "user",
	3: "operator",
	4: "admin",
}

func (r Role) String() string {
	return toString[r]
}

var toRole = map[string]Role{
	"unknown":  Unknown,
	"blocked":  Blocked,
	"user":     User,
	"operator": Operator,
	"admin":    Admin,
}

func ToRole(r string) Role {
	if c, ok := toRole[r]; ok {
		return c
	}

	return Unknown
}

func IsOperator(r Role) bool {
	return r >= Operator
}

func IsAdmin(r Role) bool {
	return r == Admin
}

func IsBlocked(r Role) bool {
	return r == Blocked
}

func IsUser(r Role) bool {
	return r >= User
}

func IsRoleALowerThanB(a Role, b Role) bool {
	return a < b
}