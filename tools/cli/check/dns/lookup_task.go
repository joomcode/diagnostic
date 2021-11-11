package dns

import (
	"context"
	"fmt"
	"net"

	"github.com/joomcode/diagnostic/tools/cli/logger"
	"github.com/joomcode/diagnostic/tools/cli/tasks"
)

type LookupCallback func(ips []net.IP)

func NewSystemLookupHostTask(host string) tasks.Task {
	return tasks.NewGenericTask(fmt.Sprintf("Lookup host by system resolver: %q", host), func(ctx context.Context, log logger.Logger) error {
		_, err := SystemLookupHost(ctx, log, host)
		return err
	})
}
