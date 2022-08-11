# Veracode Node Packager - Alpha Version
Please note that this is **not an official Veracode project**, not supported by Veracode in any form, and comes with no warranty whatsoever. It is simply a little pet project of mine trying to make the life of Veracode's `Node.js` customers a bit easier. Use at your own risk.

The Veracode Node Packager is a tool that packages your `Node.js` application for `Veracode Static Analysis`. The idea is to avoid common mistakes that I, in my role as a Veracode Application Security Consultant, commonly see in customer uploads.

There also is a `sample-node-project` folder that contains a "Hello World"-ish `Node.js` application with a lot of the files that we want to filter out. This test project can be used to take `Veracode Node Packager` for a spin.

Please feel free to extend the existing functionality, followed by a `Merge Request`.

## Built-in Help
Help is built-in!

- `node-packager --help` - outputs the help.

# How to Use
```text
Usage:
    node-packager [flags]

Flags:
  -source string     The path of the Node.js app you want to package (default "sample-node-project")
  -target string     The path where you want the output.zip to be stored to (default ".")
  -tests string      The path that contains your Node.js test files (relative to the source) (default "test")
```

# What does it do?
- Creates a zip of the `-source` folder and puts it into the provided `-target` directory as `upload.zip`
- `Features`: 
    - This tool creates a zip of your application ready to be uploaded to the Veracode Platform
    - It prevents common, non-required, files from being a part of the zip (such as `node_modules`, `tests`)
    - The tool also checks for "smells" that indicates something might not right with the packaging, and prints corresponding warnings if a "smell" was found
- `Files to be omitted`:
    - Omit the `node_modules` folder (usually only contains 3rd party libraries)
    - Omit the `tests` directory (that contains e.g. your unit- and integration tests)
        - Specified via `-test <path>`
    - Omit style sheets (`.css` and `.scss` files)
    - Omit images (e.g. `.jpg`, `.png`) 
    - Omit documents (e.g. `.pdf`)
    - Omit the `.git` folder
    - Omit other not-required files (e.g. `.DS_Store`)
- `Additional Checks`:
    - Check if `package-lock.json` exists (this is required for Veracode SCA)
    - Check if `public` exists (may contain resources that are not part of your source code)
    - Check if `dist` exists (may contain resources that are not part of your source code)
    - Check for `.map` files (indicates that your JS files might be minified)

# Setup
- You can simply run it from source via `go run main.go` 
- You can build the project yourself simply via `go build`

# Releases
- The `/releases` folder contains some binary releases already so that you might not have to build it yourself