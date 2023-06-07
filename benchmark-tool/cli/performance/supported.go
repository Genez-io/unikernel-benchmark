package performance

type SupportedUnikernel struct {
	UnikernelName string
	SupportedVMMs []string
}

var SupportedUnikernels = map[string]SupportedUnikernel{
	"cloudius-systems/osv": {
		UnikernelName: "osv",
		SupportedVMMs: []string{"qemu", "firecracker"},
	},
	"unikraft/unikraft": {
		UnikernelName: "unikraft",
		SupportedVMMs: []string{"qemu"},
	},
}
