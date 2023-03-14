package performance

var SupportedUnikernels = map[string]interface{}{
	"cloudius-systems/osv": OSvDriver,
	"unikraft/unikraft":    UnikraftDriver,
}
