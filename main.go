// Command guard checks a Go repository against the naming, interface and
// init rules of its AGENTS.md. Today's violations are recorded in an
// allowlist that may only shrink: a new violation fails, and so does an entry
// that no longer matches, so fixed names come out of the list.
//
//	go run github.com/0magnet/guard@main names
//	go run github.com/0magnet/guard@main inits
//	go run github.com/0magnet/guard@main names --write   # record today's
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/0magnet/calvin/clihelp"
)

var (
	dir   string
	list  bool
	write bool
)

var rootCmd = &cobra.Command{
	Use:   "guard",
	Short: "check a Go repository's names, interfaces and inits against its allowlists",
	Long: "Checks the Go files git tracks in the working directory against the rules of\n" +
		"AGENTS.md. Each check has an allowlist of today's violations, which may only\n" +
		"shrink. Tests, generated files, vendor, node_modules, testdata and dot\n" +
		"directories are skipped.",
	SilenceUsage:  true,
	SilenceErrors: true,
}

var checks = []struct{ name, short string }{
	{"names", "generic names (Data, Info, Manager, Get, ...), util-style packages, and is/has bools"},
	{"ifaces", "interfaces that embed another interface"},
	{"inits", "func init() outside package main"},
}

func main() {
	f := rootCmd.PersistentFlags()
	f.StringVar(&dir, "dir", "scripts/guard", "directory holding the allowlists")
	f.BoolVar(&list, "list", false, "print the current violations instead of checking")
	f.BoolVar(&write, "write", false, "record the current violations as the allowlist")
	for _, c := range checks {
		rootCmd.AddCommand(&cobra.Command{
			Use:   c.name,
			Short: "fail on new " + c.short,
			Args:  cobra.NoArgs,
			RunE: func(_ *cobra.Command, _ []string) error {
				return check(c.name, dir, list, write)
			},
		})
	}
	clihelp.Init(rootCmd, "guard")
	if err := rootCmd.Execute(); err != nil {
		if !errors.Is(err, errFindings) {
			fmt.Fprintln(os.Stderr, "guard:", err)
		}
		os.Exit(1)
	}
}
