# Veracode JavaScript Packager - Beta Version
Please note that this is **not an official Veracode project**, not supported by Veracode in any form, and comes with no warranty whatsoever. It is simply a little pet project of mine trying to make the life of Veracode's `JavaScript` customers a bit easier. Use at your own risk.

The `Veracode JavaScript Packager` is a tool that packages your `JavaScript` applications (i.e., `Node.js`, `Angular`, `React`, or `Vue`) for `Veracode Static Analysis`. The idea is to avoid common mistakes that I, in my role as a Veracode Application Security Consultant, commonly see in customer uploads.

There also is a set of sample applications (in `./sample-projects`) that can be used to test to take this tool for a spin.

Please feel free to extend the existing functionality, followed by a `Merge Request`.

## Built-in Help
Help is built-in!

- `veracode-js-packager --help` - outputs the help.

# How to Use
```text
Usage:
    veracode-js-packager [flags]

Flags:
  -source string     The path of the JavaScript app you want to package (default "./sample-projects/sample-node-project")
  -target string     The path where you want the vc-output.zip to be stored to (default ".")
  -tests string      The path that contains your test files (relative to the source) (default "").  Uses a heuristic to identifiy tests automatically in case no path is provided
  
```

# What does it do?
- Creates a zip of the `-source` folder and puts it into the provided `-target` directory as `vc-output.zip`
- `Features`: 
    - This tool creates a zip of your application ready to be uploaded to the Veracode Platform
    - It prevents common, non-required, files from being a part of the zip (such as `node_modules`, `tests`)
    - The tool also checks for "smells" that indicate something might not be right with the packaging, and prints corresponding warnings/errors if a "smell" was found
- `Omitted Files/Folders`:
    - Omit the `node_modules` folder (usually only contains 3rd party libraries)
    - Omit the `tests` directory (that contains e.g. your unit- and integration tests)
        - Specified via `-tests <path>`
    - Omit style sheets (`.css` and `.scss` files)
    - Omit images (e.g. `.jpg`, `.png`) 
    - Omit documents (e.g. `.pdf`)
    - Omit the `.git` folder
    - Omit other non-required files (e.g. `.DS_Store`)
- `Additional Checks`:
    - Check if `package-lock.json` exists (this is required for Veracode SCA)
    - Check if `/public` exists (may contain resources that are not part of your actual 1st party source code)
    - Check if `/dist` exists (may contain minified JavaScript)
    - Check for `.map` files (indicates that your JS files might be minified)

# Setup
- You can simply run this tool from source via `go run main.go` 
- You can build the tool yourself via `go build`

# Releases
- The `Releases` section contains some binary releases already so that you might not have to build it yourself
