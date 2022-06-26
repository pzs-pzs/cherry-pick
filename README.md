# cherry-pick

A pure go tool to analyze cherry-pick commit, output yaml file

## install

```
go install github.com/pzs-pzs/cherry-pick/cmd/pick@latest
```

## use

```shell
pick analyze -r "https://github.com/pzs-pzs/galen.git" -o out.yaml
```

## feature

- [x] support analyze repository use   ```git cherry-pick -x```
- [x] support https schema
- [ ] support all cherry-pick commit analysis

## lib

```url
https://github.com/go-git/go-git
```