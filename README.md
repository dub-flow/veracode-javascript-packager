# Veracode Node Packager - Alpha Version
The Veracode Node Packager is a tool that packages your `Node.js` application for `Veracode Static Analysis`. The idea is to avoid common mistakes that I, in my role as a Veracode Application Security Consultant, commonly see in customer uploads.

Please note that this is **not an official Veracode project**, not supported by Veracode in any form, and comes with no warranty whatsoever. It is simply a little pet project of mine trying to make the life of Veracode's `Node.js` customers a bit easier. Use at your own risk.

There also is a `sample-node-project` folder that contains a "Hello World"-ish `Node.js` application with a lot of the files that we want to filter out. This test project can be used to take `Veracode Node Packager` for a spin.

## Built-in Help
Help is built-in!

- `node-packager --help` - outputs the help.

# How to Use
```text
Usage:
    node-packager [flags]

Flags:
  -source string     The path of the Node.js app you want to package (default "../test-projects/my-node-test")
  -target string     The path where you want the output.zip to be stored to (default ".")
```

# Setup
- You can simply run it from source via `go run src/main.go` 
- How to build..
- Releases?

# What does it do?
- Creates a zip of the `-source` folder and puts it into the provided `-target` directory as `upload.zip`
- `Features`:
    - Omit `node_modules`
    - Omit style sheets (`.css` and `.scss` files)
    - Omit popular image extensions (e.g. `.jpg`, `.png`) 
    - Omit popular document extensions (e.g. `.pdf`)
    - Omit `.git` folder
    - Omit other not-required files (e.g. `.DS_Store`)

# Upcoming Features (hopefully)
- Omit `test` folder
- Omit `public` folder 

