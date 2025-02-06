# roots

`roots` is a tool for exploring multiple root directories, such as those in a monorepo project.

## Usage

```console
$ pwd
/path/to/src/github.com/fabrix-framework/fabrix
$ roots
/path/to/src/github.com/fabrix-framework/fabrix
/path/to/src/github.com/fabrix-framework/fabrix/examples/mock-todoapp-server
/path/to/src/github.com/fabrix-framework/fabrix/examples/vite-todoapp
/path/to/src/github.com/fabrix-framework/fabrix/packages/chakra-ui
/path/to/src/github.com/fabrix-framework/fabrix/packages/fabrix
/path/to/src/github.com/fabrix-framework/fabrix/packages/graphql-config
/path/to/src/github.com/fabrix-framework/fabrix/packages/unstyled
/path/to/src/github.com/fabrix-framework/fabrix/shared/eslint
/path/to/src/github.com/fabrix-framework/fabrix/shared/prettier
$ cd packages/fabrix
$ pwd
/path/to/src/github.com/fabrix-framework/fabrix/packages/fabrix
$ roots
/path/to/src/github.com/fabrix-framework/fabrix
/path/to/src/github.com/fabrix-framework/fabrix/examples/mock-todoapp-server
/path/to/src/github.com/fabrix-framework/fabrix/examples/vite-todoapp
/path/to/src/github.com/fabrix-framework/fabrix/packages/chakra-ui
/path/to/src/github.com/fabrix-framework/fabrix/packages/fabrix
/path/to/src/github.com/fabrix-framework/fabrix/packages/graphql-config
/path/to/src/github.com/fabrix-framework/fabrix/packages/unstyled
/path/to/src/github.com/fabrix-framework/fabrix/shared/eslint
/path/to/src/github.com/fabrix-framework/fabrix/shared/prettier
```

### With [ghq](https://github.com/x-motemen/ghq) and [peco](https://github.com/peco/peco)

```console
$ cd `ghq list --full-path | roots | peco`
```

## Explore Strategy

1. First, the root directory is explored. The number of root directories to explore is specified by `--parent` (default: 2).
2. Next, sub-root directories are explored from the root directory. The number of directories to explore is specified by `--depth` (default: 3).

The root directory is determined by checking whether or not the specified file (root file) exists in the target directory.

The default root files are `.git/config`, `go.mod`, `package.json`, and `Cargo.toml`.

You can specify multiple root files with `--root-file`.

## My most favorite [The Roots](https://www.theroots.com/) album

[Phrenology](https://en.wikipedia.org/wiki/Phrenology_(album))

## My most favorite [The Roots](https://www.theroots.com/) song

[The Next Movement](https://www.youtube.com/watch?v=qm7Xt2Qsjcg)

## Install

**homebrew tap:**

```console
$ brew install k1LoW/tap/roots
```

**go install:**

```console
$ go install github.com/k1LoW/roots@latest
```

**manually:**

Download binary from [releases page](https://github.com/k1LoW/roots/releases)

