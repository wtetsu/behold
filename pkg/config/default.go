/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package config

// Default returns the default configuration
func Default() string {
	return `# Gaze configuration(priority: default < ~/.gaze.yml < ~/.config/gaze/gaze.yml < -f option)
commands:
- ext: .bash
  cmd: bash "{{file}}"
- ext: .cpp
  cmd: |
    gcc "{{file}}" -o"{{base0}}.out"
    ./"{{base0}}.out"
- ext: .d
  cmd: dmd -run "{{file}}"
- ext: .go
  cmd: go run "{{file}}"
- ext: .groovy
  cmd: groovy "{{file}}"
- ext: .java
  cmd: java "{{file}}"
- ext: .js
  cmd: node "{{file}}"
- ext: .kts
  cmd: kotlinc -script "{{file}}"
- ext: .php
  cmd: php "{{file}}"
- ext: .py
  cmd: python "{{file}}"
- ext: .rb
  cmd: ruby "{{file}}"
- ext: .rs
  cmd: |
    rustc "{{file}}" -o"{{base0}}.out"
    ./"{{base0}}.out"
- ext: .sh
  cmd: sh "{{file}}"
- ext: .ts
  cmd: |
    tsc "{{file}}" --out "{{base0}}.out"
    node ./"{{base0}}.out"
- re: ^Dockerfile$
  cmd: docker build -f "{{file}}" .
`
}
