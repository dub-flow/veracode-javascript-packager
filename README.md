# Veracode Node Packager - Alpha Version
The Veracode Node Packager is a tool that packages your `Node.js` application for `Veracode Static Analysis`. The idea is to avoid common mistakes that I, in my role as a Veracode Application Security Consultant, commonly see in customer uploads.

Please note that this is **not an official Veracode project**, not supported by Veracode in any form, and comes with no warranty whatsoever. It is simply a little pet project of mine trying to make the life of Veracode's `Node.js` customers a bit easier. Use at your own risk.

## Built-in Help

Help is built-in!

- `node-packager --help` - outputs the help.

# Installation
- Via `go run src/main.go`
- How to build..
- Releases?

# How to Run
```text
Usage:
    node-packager [flags]

Flags:
  -source string     The path of the Node.js app you want to package (default "../test-projects/my-node-test")
  -target string     The path where you want the output.zip to be stored to (default ".")
```

# What does it do?
- Creates a zip of a source folder (`-source <path>`) and puts it into the provided target directory (`-target <path>`) as `upload.zip`
- `Features`:
    - Omit `node_modules`
    - Omit `.css` and `.scss` files

# Upcoming Features (hopefully)
    - Omit images
    - Omit `test` folder
    - Omit `public` folder 


# Test Projects
- To test this packager, we have used `The node.js example app` from `https://github.com/contentful/the-example-app.nodejs`
- After downloading it `npm install`