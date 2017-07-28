# Valve Data Format
Package for working with [Valve Data Format](https://developer.valvesoftware.com/wiki/KeyValues) writen in Go. Lexer based on Rob Pike's [talk](https://youtu.be/HxaD_trXwRE).

###### I found some already existed parsers written in Go, but I did not like them and I decided to write my own.

## Start using it
1. Download and install it:
```go
go get github.com/blukai/vdf
```
2. Import it in your code:
```go
import "github.com/blukai/vdf"
```
3. Parse VDF file:
```go
dat, err := ioutil.ReadFile("./vdfFile.txt")
if err != nil {
	log.Fatal(err)
}
m := vdf.Parse(string(dat))
```
...and u can easily convert it's to json by adding:
```go
jsonized, err := json.MarshalIndent(m, "", "  ")
if err != nil {
	log.Fatal(err)
}

err = ioutil.WriteFile("./jsonFile.json", jsonized, 0644)
if err != nil {
	log.Fatal(err)
}
```

---

If you found a bug, pr's are welcome, otherwise open an issue.
