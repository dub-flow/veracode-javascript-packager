![Go Version](https://img.shields.io/github/go-mod/go-version/dub-flow/veracode-javascript-packager)
![GitHub Downloads](https://img.shields.io/github/downloads/dub-flow/veracode-javascript-packager/total)
![Docker Image Size](https://img.shields.io/docker/image-size/fw10/veracode-js-packager/latest)
![Docker Pulls](https://img.shields.io/docker/pulls/fw10/veracode-js-packager)


# Veracode JavaScript Packager ⚡

Please note that this is **not an official Veracode project**, not supported by Veracode in any form, and comes with no warranty whatsoever. It is simply a little pet project of mine trying to make the life of Veracode's `JavaScript`/`TypeScript` customers a bit easier. Use at your own risk.

The `Veracode JavaScript Packager` is a tool that packages your `JavaScript`/`TypeScript` applications (i.e., `Node.js`, `Angular`, `React`, or `Vue`) for `Veracode Static Analysis`. The output of the tool is a `zip` that you then need to upload to Veracode. The idea is to avoid common mistakes that I, in my role as a Veracode Application Security Consultant, commonly see in customer uploads.

There also is a set of sample applications (in `./sample-projects`) that can be used to test to take this tool for a spin.

Please feel free to extend the existing functionality, followed by a `Merge Request` ❤️.

# Built-in Help 🆘

Help is built-in!

- `veracode-js-packager --help` - outputs the help.

# How to Use ⚙

```text
Usage:
    veracode-js-packager [flags]

Flags:
  -source string     The path to the JavaScript app you want to package (required)
  -target string     The path where you want the vc-output.zip to be stored to (default ".")
  -tests string      The path that contains your test files (relative to the source). (default: Uses a heuristic to identify tests automatically in case no path is provided)

Examples:
    ./veracode-js-packager -source my-js-app -target . 
    ./veracode-js-packager -source my-js-app -target . -tests tests
```

# What does it do? 🔎 

- Creates a zip of the `-source` folder and puts it into the provided `-target` directory as `vc-output.zip`
- `Features`: 
    - This tool creates a zip of your application ready to be uploaded to the Veracode Platform
    - It prevents common, non-required, files from being a part of the zip (such as `node_modules`, `tests`)
    - The tool also checks for "smells" that indicate something might not be right with the packaging, and prints corresponding warnings/errors if a "smell" was found
- `Omitted Files/Folders`:
    - Omit the `node_modules` folder (usually only contains 3rd party libraries)
    - Omit the `tests` directory (that contains e.g. your unit- and integration tests)
        - Specified via `-tests <path>`
    - Omit style sheets (e.g. `.css` and `.scss` files)
    - Omit images (e.g. `.jpg`, `.png`) and videos (e.g. `.mp4`)
    - Omit documents (e.g. `.pdf`, `.docx`)
    - Omit the `.git` folder
    - Omit fonts
    - ...

# Setup ✅

- You can simply run this tool from source via `go run .` 
- You can build the tool yourself via `go build`
- You can build the `docker` image yourself via `docker build . -t fw10/veracode-js-packager`
    - `Note`: This would only work on Unix at the moment

# Run Tests 🧪

- To run the tests, run `go test` or `go test -v` (for more details)

# Run via Docker 🐳

1. Traverse **into** the directory of the `JavaScript app` that you want to package
2. From within there, run `docker run -it --rm -v "$(pwd):/app/js-app" --name packager fw10/veracode-js-packager`

# Run from within Azure DevOps ☁️

- To run the tool from within `Azure DevOps` with a `Linux` command line, you can copy the below task into your pipeline script (note that you need to change `-source` and `-target`)

```text
- task: CmdLine@2
  displayName: 'Veracode JavaScript Packager'
  inputs:
    script: |
      wget https://github.com/fw10/veracode-javascript-packager/releases/latest/download/veracode-js-packager-linux-amd64
      chmod +x veracode-js-packager-linux-amd64
      ./veracode-js-packager-linux-amd64 -source <path-to-js-app> -target <path-of-output-zip>
```

# Releases 🔑 

- The `Releases` section contains some already compiled binaries for you so that you might not have to build the tool yourself
- For the `Mac releases`, your Mac will throw a warning (`"cannot be opened because it is from an unidentified developer"`)
    - To avoid this warning in the first place, you could simply build the app yourself (see Setup)
    - Alternatively, you may - at your own risk - bypass this warning following the guidance here: https://support.apple.com/guide/mac-help/apple-cant-check-app-for-malicious-software-mchleab3a043/mac
    - Afterwards, you can simply run the binary from the command line and provide the required flags

# Kubernetes ☸ 

- You can now also run the app via `Kubernetes` (I don't really see a reason for doing so at this point though)
- To do this, you have to change the `hostPath.path` in the `k8s-manifest.yml` to the **absolute** path where the JavaScript app resides that you want to be package
- Afterwards, run `kubectl apply -f k8s-manifest.yml`
    - This will output the `zip` (that you can upload to `Veracode Static Analysis`) into the provided `hostPath.path`

# Bug Reports 🐞

If you find a bug, please file an Issue right here in GitHub, and I will try to resolve it in a timely manner.
