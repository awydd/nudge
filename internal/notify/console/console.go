package console

import (
	"fmt"
	"os"
	"text/tabwriter"
)

type ConsoleNotifier struct {
	Title       string
	Description string
	Date        string
	Years       int
}

func (c *ConsoleNotifier) Send() error {
	fmt.Println("纪念日提醒")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "标题:\t%s (%d周年)\n", c.Title, c.Years)
	fmt.Fprintf(w, "日期:\t%s\n", c.Date)
	fmt.Fprintf(w, "描述:\t%s\n", c.Description)
	w.Flush()

	fmt.Println()
	return nil
}
