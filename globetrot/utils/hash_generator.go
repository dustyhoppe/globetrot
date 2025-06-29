package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"hash"
	"strings"
)

const (
	windowsLineEnding = "\r\n"
	unixLineEnding    = "\n"
	macLineEnding     = "\r"
)

type HashGenerator struct {
	hasher           hash.Hash
	normalizeEndings bool
}

func (h *HashGenerator) Init(normalizeEndings bool) {
	h.normalizeEndings = normalizeEndings
	h.hasher = sha256.New()
}

func (h *HashGenerator) GenerateHash(script string) string {
        if h.normalizeEndings {
                script = strings.ReplaceAll(script, windowsLineEnding, unixLineEnding)
                script = strings.ReplaceAll(script, macLineEnding, unixLineEnding)
        }

        b := []byte(script)
        h.hasher.Reset()
        h.hasher.Write(b)
        sha := base64.URLEncoding.EncodeToString(h.hasher.Sum(nil))

	return sha
}
