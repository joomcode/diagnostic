package dns

import (
	"fmt"
	"github.com/joomcode/diagnostic/tools/cli/tasks"
	"net"
)

type LookupCallback func(ips []net.IP)

func NewSystemLookupHostTask(host string) tasks.Task {
	return tasks.NewGenericTask(fmt.Sprintf("Lookup host by system resolver: %q", host))
}
