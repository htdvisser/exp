package envcrypto

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/exp/slices"
)

type Source interface {
	Lookup(key string) (string, bool)
	Keys() []string
}

type EnvSource struct{}

func (s EnvSource) Lookup(key string) (string, bool) {
	return os.LookupEnv(key)
}

func (s EnvSource) Keys() []string {
	environ := os.Environ()
	keys := make([]string, len(environ))
	for i, kv := range environ {
		keys[i], _, _ = strings.Cut(kv, "=")
	}
	return keys
}

type MapSource map[string]string

func (s MapSource) Lookup(key string) (string, bool) {
	v, ok := s[key]
	return v, ok
}

func (s MapSource) Keys() []string {
	keys := make([]string, len(s))
	i := 0
	for k := range s {
		keys[i] = k
		i++
	}
	return keys
}

type multiSource[X Source] []X

func (s multiSource[X]) Lookup(key string) (string, bool) {
	for _, source := range s {
		if value, ok := source.Lookup(key); ok {
			return value, ok
		}
	}
	return "", false
}

func (s multiSource[X]) Keys() []string {
	var keys []string
	for _, source := range s {
		for _, key := range source.Keys() {
			if !slices.Contains(keys, key) {
				keys = append(keys, key)
			}
		}
	}
	return keys
}

type MultiSource = multiSource[Source]

type EnvFileSource struct {
	fsys   fs.FS
	path   string
	data   []byte
	parsed MapSource
}

func NewEnvFileSource(fsys fs.FS, name string) (*EnvFileSource, error) {
	var (
		data []byte
		err  error
	)
	if fsys != nil {
		data, err = fs.ReadFile(fsys, name)
	} else {
		data, err = os.ReadFile(name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read %q: %w", name, err)
	}
	parsed, err := godotenv.UnmarshalBytes(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %q: %w", name, err)
	}
	return &EnvFileSource{
		fsys:   fsys,
		path:   name,
		data:   data,
		parsed: parsed,
	}, nil
}

func (s *EnvFileSource) Replace(f func(data []byte) ([]byte, error)) error {
	data, err := f(s.data)
	if err != nil {
		return err
	}
	if bytes.Equal(data, s.data) {
		return nil
	}
	parsed, err := godotenv.UnmarshalBytes(data)
	if err != nil {
		return fmt.Errorf("failed to parse %q: %w", s.path, err)
	}
	s.data, s.parsed = data, parsed
	return nil
}

func (s *EnvFileSource) Lookup(key string) (string, bool) {
	return s.parsed.Lookup(key)
}

func (s *EnvFileSource) Keys() []string {
	return s.parsed.Keys()
}

type EnvFilesSource struct {
	fsys  fs.FS
	files []*EnvFileSource
}

func NewEnvFilesSource(fsys fs.FS, paths ...string) (*EnvFilesSource, error) {
	files := make([]*EnvFileSource, len(paths))
	for i, path := range paths {
		file, err := NewEnvFileSource(fsys, path)
		if err != nil {
			return nil, err
		}
		files[i] = file
	}

	return &EnvFilesSource{
		fsys:  fsys,
		files: files,
	}, nil
}

func (s *EnvFilesSource) GetFile(path string) *EnvFileSource {
	for _, file := range s.files {
		if file.path == path {
			return file
		}
	}
	return nil
}

func (s *EnvFilesSource) PrependSource(file *EnvFileSource) {
	s.files = append([]*EnvFileSource{file}, s.files...)
}

func (s *EnvFilesSource) AppendSource(file *EnvFileSource) {
	s.files = append(s.files, file)
}

func (s *EnvFilesSource) Lookup(key string) (string, bool) {
	return multiSource[*EnvFileSource](s.files).Lookup(key)
}

func (s *EnvFilesSource) Keys() []string {
	return multiSource[*EnvFileSource](s.files).Keys()
}
