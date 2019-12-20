package main

import (
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	//"github.com/mitchellh/go-wordwrap"

	"github.com/russross/blackfriday/v2"
	"github.com/samfoo/ansi"
)

func main() {
	fmt.Println(string(blackfriday.Run([]byte(example),
		blackfriday.WithRenderer(&ConsoleRenderer{}), blackfriday.WithExtensions(Extensions))))
}

func Render(input string) string {
	return string(blackfriday.Run([]byte(input),
		blackfriday.WithRenderer(&ConsoleRenderer{}), blackfriday.WithExtensions(Extensions)))
}

// Extensions are the blackfriday extensions used when parsing markdown
var Extensions = 0 |
	blackfriday.NoIntraEmphasis |
	blackfriday.FencedCode |
	blackfriday.Autolink |
	blackfriday.Strikethrough |
	blackfriday.SpaceHeadings |
	blackfriday.BackslashLineBreak |
	blackfriday.DefinitionLists |
	blackfriday.HeadingIDs

// ConsoleRenderer renders markdown for the console
type ConsoleRenderer struct {
	heading bool
	table   bool
	code    bool
	width   int
	item    bool
	pad     int
	sync.Once
}

func (cr *ConsoleRenderer) init(w io.Writer) {
	cr.Do(func() {
		cr.width = 80
		out, err := exec.Command("stty", "size").CombinedOutput()
		if err != nil {
			return
		}
		cl := strings.TrimSpace(strings.Split(string(out), " ")[1])
		width, err := strconv.Atoi(cl)
		if err != nil {
			return
		}
		if cr.width > 120 {
			// output doesn't look good on long consoles
			cr.width = width
		}
	})
}

func (cr *ConsoleRenderer) indent(pad int) {
	cr.pad += cr.pad
}

func (cr *ConsoleRenderer) indentCenter() {
	leftPad := (cr.width - 30) / 2
	cr.pad = leftPad
	cr.pad += cr.pad
}

func (cr *ConsoleRenderer) print(w io.Writer, val string) {
	fmt.Fprintf(w,
		fmt.Sprintf("%%-%ds", cr.pad),
		fmt.Sprintf(fmt.Sprintf("%%%ds", cr.pad), val))
}

func (cr *ConsoleRenderer) RenderNode(
	w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	//fmt.Fprintf(w, "%s\n", node.String())
	cr.init(w)
	if !entering {
		fmt.Fprintf(w, "%s", ansi.ColorCode("reset"))
		cr.heading = false
	}
	switch node.Type {
	case blackfriday.Document:
	case blackfriday.BlockQuote:
	case blackfriday.List:
	case blackfriday.Item:
		cr.item = entering
		if entering {
			fmt.Fprintf(w, "\n")
			cr.indent(4)
		} else {
			cr.indent(-4)
		}
	case blackfriday.Paragraph:
		fmt.Fprintf(w, "\n")
	case blackfriday.Heading:
		fmt.Fprintf(w, "\n")
		if entering {
			fmt.Fprintf(w, ansi.ColorCode("+b"))
			cr.indentCenter()
		} else {
			fmt.Fprintf(w, "%s", ansi.ColorCode("reset"))
			cr.pad = 0
		}
	case blackfriday.HorizontalRule:
	case blackfriday.Emph:
		if entering {
			fmt.Fprintf(w, ansi.ColorCode("+h"))
		} else {
			fmt.Fprintf(w, "%s", ansi.ColorCode("reset"))
		}
	case blackfriday.Strong:
		if entering {
			fmt.Fprintf(w, ansi.ColorCode("+bh"))
		} else {
			fmt.Fprintf(w, "%s", ansi.ColorCode("reset"))
		}
	case blackfriday.Del:
	case blackfriday.Link:
	case blackfriday.Image:
	case blackfriday.Text:
		if cr.item && node.ListData.BulletChar > 0 {
			fmt.Fprintf(w, "%s", node.ListData.Delimiter)
			cr.print(w, string([]byte{node.ListData.Delimiter}))
		}
		cr.print(w, string(node.Literal))
	case blackfriday.HTMLBlock:
	case blackfriday.CodeBlock:
		cr.print(w, "\n    "+strings.ReplaceAll(string(node.Literal), "\n", "\n    "))
	case blackfriday.Softbreak:
	case blackfriday.Hardbreak:
	case blackfriday.Code:
	case blackfriday.HTMLSpan:
	case blackfriday.Table:
	case blackfriday.TableCell:
	case blackfriday.TableHead:
	case blackfriday.TableBody:
	case blackfriday.TableRow:
	}
	return blackfriday.GoToNext
}

func (ConsoleRenderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {
	// no-op
}

func (ConsoleRenderer) RenderFooter(w io.Writer, ast *blackfriday.Node) {
	// no-op
}

const example = `
## kpt

  Git based configuration package manager.

#### Installation

    go install -v sigs.k8s.io/kustomize/kustomize/v3
    go install -v github.com/GoogleContainerTools/kpt

#### Commands

- [get](commands/get.md) -- fetch a package from git and write it to a local directory

      kpt help get # in-command help

      kpt get https://github.com/GoogleContainerTools/kpt/examples/cockroachdb@v0.1.0 my-cockroachdb
      kustomize config tree my-cockroachdb --name --replicas --image

      my-cockroachdb
      ├── [cockroachdb-statefulset.yaml]  Service cockroachdb
      ├── [cockroachdb-statefulset.yaml]  StatefulSet cockroachdb
      │   ├── spec.replicas: 3
      │   └── spec.template.spec.containers
      │       └── 0
      │           ├── name: cockroachdb
      │           └── image: cockroachdb/cockroach:v1.1.0
      ├── [cockroachdb-statefulset.yaml]  PodDisruptionBudget cockroachdb-budget
      └── [cockroachdb-statefulset.yaml]  Service cockroachdb-public

- [diff](commands/diff.md) -- display a diff between the local package copy and the upstream version

      kpt help diff # in-command help

      sed -i -e 's/replicas: 3/replicas: 5/g' my-cockroachdb/cockroachdb-statefulset.yaml
      kpt diff my-cockroachdb

      diff ...
      <   replicas: 5
      ---
      >   replicas: 3

- [update](commands/update.md) -- pull upstream package changes

      kpt help update # in-command help

      # commiting to git is required before update
      git add . && git commit -m 'updates'
      kpt update my-cockroachdb@v0.2.0

- [sync](commands/sync.md) -- declaratively manage a collection of packages

      kpt help sync # in-command help

          # dir/Kptfile
          apiVersion: kpt.dev/v1alpha1
          kind: Kptfile
          dependencies:
          - name: my-cockroachdb
            git:
              repo: "https://github.com/GoogleContainerTools/kpt"
              directory: "examples/cockroachdb"
              ref: "v0.1.0"

      kpt sync dir/

- [desc](commands/desc.md) -- show the upstream metadata for one or more packages

      kpt help desc # in-command help

      kpt desc my-cockroachdb

       PACKAGE NAME         DIR                         REMOTE                       REMOTE PATH        REMOTE REF   REMOTE COMMIT  
      my-cockroachdb   my-cockroachdb   https://github.com/kubernetes/examples   /staging/cockroachdb   master       a32bf5c        

- [man](commands/man.md) -- render the README.md from a package if possible (requires man2md README format)

      kpt help man # in-command help

      kpt man my-cockroachdb

- [init](commands/init.md) -- initialize a new package with a README.md (man2md format) and empty Kptfile
  (optional)

      mkdir my-new-package
      kpt init my-new-package/

      tree my-new-package/
      my-new-package/
      ├── Kptfile
      └── README.md

#### Design

1. **Packages are composed of Resource configuration** (rather than DSLs, templates, etc)
    * May also contain supplemental non-Resource artifacts (e.g. README.md, arbitrary other files).

2.  **Any existing git subdirectory containing Resource configuration** may be used as a package.
    * Nothing besides a git directory containing Resource configuration is required.
    * e.g. the [examples repo](https://github.com/kubernetes/examples/staging/cockroachdb) may
      be used as a package:

          # fetch the examples cockroachdb directory as a package
          kpt get https://github.com/kubernetes/examples/staging/cockroachdb my-cockroachdb

3. **Packages should use git references for versioning**.
    * Package authors should use semantic versioning when publishing packages.

          # fetch the examples cockroachdb directory as a package
          kpt get https://github.com/GoogleContainerTools/kpt/examples/cockroachdb@v0.1.0 my-cockroachdb

4. **Packages may be modified or customized in place**.
    * It is possible to directly modify the fetched package.
    * Tools may set or change fields.
    * [Kustomize functions](https://github.com/kubernetes-sigs/kustomize/blob/master/cmd/config/docs/commands/run-fns.md)
      may also be applied to the local copy of the package.

          export KUSTOMIZE_ENABLE_ALPHA_COMMANDS=true

          kpt get https://github.com/GoogleContainerTools/kpt/examples/cockroachdb my-cockroachdb
          kustomize config set my-cockroachdb/ replicas 5

5. **The same package may be fetched multiple times** to separate locations.
    * Each instance may be modified and updated independently of the others.

          export KUSTOMIZE_ENABLE_ALPHA_COMMANDS=true

          # fetch an instance of a java package
          kpt get https://github.com/GoogleContainerTools/kpt/examples/java my-java-1
          kustomize config set my-java-1/ image gcr.io/example/my-java-1:v3.0.0

          # fetch a second instance of a java package
          kpt get https://github.com/GoogleContainerTools/kpt/examples/java my-java-2
          kustomize config set my-java-2/ image gcr.io/example/my-java-2:v2.0.0

6. **Packages may pull upstream updates after they have been fetched and modified**.
    * Specify the target version to update to, and an (optional) update strategy for how to apply the
      upstream changes.

          export KUSTOMIZE_ENABLE_ALPHA_COMMANDS=true

          kpt get https://github.com/GoogleContainerTools/kpt/examples/cockroachdb my-cockroachdb
          kustomize config set my-cockroachdb/ replicas 5
          kpt update my-cockroachdb@v1.0.1 --strategy=resource-merge


#### Templates and DSLs

Note: If the use of Templates or DSLs is strongly desired, they may be fully expanded into Resource
configuration to be used as a kpt package.  These artifacts used to generated Resource configuration
may be included in the package as supplements.

#### Env Vars

  COBRA_SILENCE_USAGE
  
    Set to true to silence printing the usage on error

  COBRA_STACK_TRACE_ON_ERRORS

    Set to true to print a stack trace on an error

  KPT_NO_PAGER_HELP

    Set to true to print the help to the console directly instead of through
    a pager (e.g. less)
`
