# Veracode JavaScript Packager - Beta Version ⚡

Please note that this is **not an official Veracode project**, not supported by Veracode in any form, and comes with no warranty whatsoever. It is simply a little pet project of mine trying to make the life of Veracode's `JavaScript` customers a bit easier. Use at your own risk.

The `Veracode JavaScript Packager` is a tool that packages your `JavaScript` applications (i.e., `Node.js`, `Angular`, `React`, or `Vue`) for `Veracode Static Analysis`. The idea is to avoid common mistakes that I, in my role as a Veracode Application Security Consultant, commonly see in customer uploads.

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
  -s, --source string   The path to the JavaScript app you want to package
  -t, --target string   The path where you want the vc-output.zip to be stored to (default ".")
  --tests string        The path that contains your test files (relative to the source). Uses a heuristic to identify tests automatically in case no path is provided (default "")
  --debug bool          Sets the log level to Debug if set to '-debug=true' (default "false")

Examples:
    ./veracode-js-packager -s my-js-app -t . 
    ./veracode-js-packager --source my-js-app --target . --tests tests --debug=false
```

# What does it do? 🔎 

- Creates a zip of the `--source` folder and puts it into the provided `--target` directory as `vc-output.zip`
- `Features`: 
    - This tool creates a zip of your application ready to be uploaded to the Veracode Platform
    - It prevents common, non-required, files from being a part of the zip (such as `node_modules`, `tests`)
    - The tool also checks for "smells" that indicate something might not be right with the packaging, and prints corresponding warnings/errors if a "smell" was found
- `Omitted Files/Folders`:
    - Omit the `node_modules` folder (usually only contains 3rd party libraries)
    - Omit the `tests` directory (that contains e.g. your unit- and integration tests)
        - Specified via `--tests <path>`
    - Omit style sheets (e.g. `.css` and `.scss` files)
    - Omit images (e.g. `.jpg`, `.png`) and videos (e.g. `.mp4`)
    - Omit documents (e.g. `.pdf`, `.docx`)
    - Omit the `.git` folder
    - Omit fonts
    - ...

# Setup ✅

- You can simply run this tool from source via `go run .` 
- You can build the tool yourself via `go build`

# Run Tests 🧪

- To run the tests, run `go test` or `go test -v` (for more details)

# Releases 🔑 

- The `Releases` section contains some already compiled binaries for you so that you might not have to build the tool yourself
- For the `Mac M1 release`, your Mac will throw a warning (`"cannot be opened because it is from an unidentified developer"`)
    - To avoid this warning in the first place, you could simply build the app yourself (see Setup)
    - Alternatively, you may - at your own risk - bypass this warning following the guidance here: https://support.apple.com/guide/mac-help/apple-cant-check-app-for-malicious-software-mchleab3a043/mac
    - Afterwards, you can simply run the binary from the command line and provide the required flags

