# Extract site information
This package extracts site data

## How to use
- Run `go mod tidy` for installing required packages.
- Create a build using command `go build -o fetch main.go`.
- You can use `GOOS` and `GOARCH` env for creating executable for any specific operating system.
- Run command `./fetch www.example.com` for fetching details about any specific website.
- We can use same for multiple sites `./fetch www.example.com www.google.com ...`
- Use `--metadata` flag for fetching site metadata `./fetch --metadata www.example.com www.google.com ...`