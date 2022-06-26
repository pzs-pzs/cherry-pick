# cherry-pick

A pure go tool to analyze cherry-pick commit, output yaml file

## install

```go install github.com/pzs-pzs/cherry-pick@latest```

## use

```shell
pick -r "https://github.com/pzs-pzs/galen" -o rst.yaml
```

## feature

- [x] support analyze repository use ```git cherry-pick -x```
- [ ] support all cherry-pick commit analysis

## lib

```url
https://github.com/go-git/go-git
```