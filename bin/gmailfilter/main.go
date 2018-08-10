package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/mbrt/gmailfilter/pkg/config"
	"github.com/mbrt/gmailfilter/pkg/export"
)

func readConfig(path string) (config.Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return config.Config{}, errors.Wrap(err, fmt.Sprintf("cannot read %s", path))
	}

	var res config.Config
	err = yaml.Unmarshal(b, &res)
	return res, err

}

func fatal(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	if !strings.HasSuffix(format, "\n") {
		// Add newline
		fmt.Fprintln(os.Stderr, "")
	}
	os.Exit(1)
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fatal("usage: %s <config-file>", os.Args[0])
	}
	cfg, err := readConfig(flag.Arg(0))
	if err != nil {
		fatal("error in config parse: %s", err)
	}

	rules, err := export.GenerateRules(cfg)
	if err != nil {
		fatal("error generating rules: %s", err)
	}

	err = export.DefaultXMLExporter().MarshalEntries(cfg.Author, rules, os.Stdout)
	if err != nil {
		fatal("error exporting to XML: %s", err)
	}
}