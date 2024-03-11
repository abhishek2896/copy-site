# Extract site information
This package extracts site data

## How to use

### In machine

- Run `go mod tidy` for installing required packages.
- Create a build using command `go build -o fetch main.go`.
- You can use `GOOS` and `GOARCH` env for creating executable for any specific operating system.
- Run command `./fetch www.example.com` for fetching details about any specific website.
- We can use same for multiple sites `./fetch www.example.com www.google.com ...`
- Use `--metadata` flag for fetching site metadata `./fetch --metadata www.example.com www.google.com ...`

### Using Docker

- Create a docker image using `docker build --tag fetch .`. You can change image name by replacing `fetch` in the command.
- Run normal fetch command using `-v` flag to sync the file created inside container and local `docker run -v <absoulte_path_where_you_want_the_files>:/app fetch https://www.google.com`
- Use `--metadata` flag for fetching site metadata `docker run fetch --metadata https://www.google.com`