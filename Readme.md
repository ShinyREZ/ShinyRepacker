# ShinyRepacker

ShinyRepacker focus on unpacking and repacking resources bundled by TexturePacker, which is used by THE IDOLM@STER SHINY COLORS to create sprite images.

With this tool, you can edit those packed images easily, which may be useful for localization projects.

## Installation

```bash
go get https://github.com/ShinyREZ/ShinyRepacker
```

## Usage

```
Usage of ShinyRepacker:
  -file string
        Json file name without extension
  -image string
        Image file name, used as input in unpack mode or as output in repack mode
  -mode string
        Work mode, unpack or repack (default "unpack")
  -prefix string
        Prefix path of exported files (default "unpacked")
```

## Examples

```bash
# unpack
ShinyRepacker --mode unpack --file icon

# repack
ShinyRepacker --mode repack --file icon

# batch unpack
arr=(`ls -1 **/*.json`) && for i in "${arr[@]}";do ShinyRepacker --file "$i"; done
```

## License

[MIT](./LICENSE)