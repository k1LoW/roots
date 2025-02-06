# roots

`roots` is a tool for exploring multiple root directories.

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

### With [ghq](https://github.com/x-motemen/ghq)

```console
$ ghq list --full-path | roots
```

## My most favorite [The Roots](https://www.theroots.com/) album

[Phrenology](https://en.wikipedia.org/wiki/Phrenology_(album))

## My most favorite [The Roots](https://www.theroots.com/) song

[The Next Movement](https://www.youtube.com/watch?v=qm7Xt2Qsjcg)

