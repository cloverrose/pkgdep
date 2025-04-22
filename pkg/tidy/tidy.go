package tidy

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/cloverrose/pkgdep/pkg/inspector"
	"github.com/cloverrose/pkgdep/pkg/orderedmap"
)

const keepAnnotation = "@keep"

var (
	ErrNotFound    = errors.New("not found")
	ErrInvalidKind = errors.New("invalid yaml node kind")
)

// YAMLKindString maps yaml.Kind to human readable string
var yamlKindString = map[yaml.Kind]string{
	yaml.DocumentNode: "Document",
	yaml.SequenceNode: "Sequence",
	yaml.MappingNode:  "Mapping",
	yaml.ScalarNode:   "Scalar",
	yaml.AliasNode:    "Alias",
}

// Config holds tidy configuration
type Config struct {
	ConfigFile  string // path to .pkgdep.yaml
	InspectFile string // path to used_rules.csv
	OutputFile  string // path to output file
}

// Tidy clean up configuration
type Tidy struct {
	cfg       Config
	root      *yaml.Node
	deps      *orderedmap.OrderedMap
	inspector *inspector.Inspector
}

func New(cfg Config) *Tidy {
	return &Tidy{cfg: cfg}
}

// Tidy removes unused rules from the config file
func (t *Tidy) Tidy() error {
	// Read config file
	configData, err := os.ReadFile(t.cfg.ConfigFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Create two independent readers from the byte slice
	configReader1 := bytes.NewReader(configData)
	configReader2 := bytes.NewReader(configData)

	inspectFile, err := os.Open(t.cfg.InspectFile)
	if err != nil {
		return fmt.Errorf("failed to open inspect file: %w", err)
	}
	defer inspectFile.Close()

	outFile, err := os.Create(t.cfg.OutputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	return t.tidyCore(configReader1, configReader2, inspectFile, outFile)
}

// tidyCore is core logic of tidy by using io.Reader and io.Writer
func (t *Tidy) tidyCore(configReader1, configReader2, inspectReader io.Reader, w io.Writer) error {
	if err := t.load(configReader1, configReader2, inspectReader); err != nil {
		return fmt.Errorf("failed to load: %w", err)
	}

	unusedRules := t.findUnusedRules()

	if err := t.removeRules(unusedRules); err != nil {
		return fmt.Errorf("failed to remove rules: %w", err)
	}

	if err := t.save(w); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	return nil
}

func (t *Tidy) load(configReader1, configReader2, inspectReader io.Reader) error {
	// Load YAML config
	var node yaml.Node
	if err := yaml.NewDecoder(configReader1).Decode(&node); err != nil {
		return fmt.Errorf("failed to decode config: %w", err)
	}
	t.root = &node

	// Load dependencies
	deps, err := t.loadDependencies(configReader2)
	if err != nil {
		return fmt.Errorf("failed to load dependencies: %w", err)
	}
	t.deps = deps

	// Load inspector
	insp, err := inspector.Load(inspectReader)
	if err != nil {
		return fmt.Errorf("failed to load inspector: %w", err)
	}
	t.inspector = insp

	return nil
}

// loadDependencies loads the dependencies from the config file
func (t *Tidy) loadDependencies(r io.Reader) (*orderedmap.OrderedMap, error) {
	// dto is the data transfer object for the config file.
	// We only need the dependencies field.
	var dto struct {
		Dependencies orderedmap.OrderedMap `yaml:"dependencies"`
	}
	if err := yaml.NewDecoder(r).Decode(&dto); err != nil {
		return nil, err
	}

	return &dto.Dependencies, nil
}

func (t *Tidy) findUnusedRules() map[string][]string {
	unused := make(map[string][]string)
	for from, toList := range t.deps.Iter() {
		for _, to := range toList {
			if !t.inspector.IsRecorded(from, to) {
				slog.Debug("found unused rule",
					slog.String("from", from),
					slog.String("to", to))
				unused[from] = append(unused[from], to)
			}
		}
	}

	slog.Info("found unused rules", slog.Int("count", len(unused)))
	return unused
}

func (t *Tidy) removeRules(unusedRules map[string][]string) error {
	depsNode, err := t.findNode(t.root.Content[0], "dependencies")
	if err != nil {
		return fmt.Errorf("failed to find dependencies node: %w", err)
	}

	// newContent is the content of the dependencies node after pruning
	newContent := make([]*yaml.Node, 0, len(depsNode.Content))

	for i := 0; i < len(depsNode.Content); i += 2 {
		fromNode := depsNode.Content[i]
		from := fromNode.Value
		toList, exists := unusedRules[from]
		if !exists || hasKeepAnnotation(fromNode) {
			// Used rule or annotated with @keep. Keep it.
			newContent = append(newContent, depsNode.Content[i], depsNode.Content[i+1])
			continue
		}

		values := depsNode.Content[i+1].Content
		newValues := make([]*yaml.Node, 0, len(values))
		for _, toNode := range values {
			if !slices.Contains(toList, toNode.Value) || hasKeepAnnotation(toNode) {
				// Used rule or annotated with @keep. Keep it.
				newValues = append(newValues, toNode)
			}
		}

		// Not empty. Keep it.
		if len(newValues) > 0 {
			depsNode.Content[i+1].Content = newValues
			newContent = append(newContent, depsNode.Content[i], depsNode.Content[i+1])
		} else {
			slog.Debug("removing empty dependency rule", slog.String("from", from))
		}
	}

	depsNode.Content = newContent
	return nil
}

func (t *Tidy) findNode(node *yaml.Node, key string) (*yaml.Node, error) {
	if node.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("%w: expected mapping, got %s",
			ErrInvalidKind, yamlKindString[node.Kind])
	}

	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return node.Content[i+1], nil
		}
	}

	return nil, fmt.Errorf("%w: key %q", ErrNotFound, key)
}

func (t *Tidy) save(w io.Writer) error {
	encoder := yaml.NewEncoder(w)
	encoder.SetIndent(2)
	return encoder.Encode(t.root)
}

func hasKeepAnnotation(node *yaml.Node) bool {
	return strings.Contains(node.LineComment, keepAnnotation) || strings.Contains(node.HeadComment, keepAnnotation)
}
