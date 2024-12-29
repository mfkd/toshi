# Toshi - A Simple Book Search Tool

Toshi is a command-line tool designed for searching books.

"Toshi" (とし / 知 or 智) is a Japanese word that translates to "wisdom" or
"knowledge."

## Installation

Ensure you have Go installed, then run:

```sh
go install github.com/mfkd/toshi
```

## Configuration

Configure `toshi` by setting either:

### Runtime Environment Variable

```sh
export DOMAIN="example.com"
```

The DOMAIN value can be changed anytime without rebuilding.

### Embedded domains.txt

File location: `./internal/embed/domains/domains.txt`. Only a single domain is supported at the moment.

```text
example.com
```

These domains are baked into the program when you build it. You'll need to rebuild to change them.

### Notes

Environment variable takes priority if both are set. Domain must be valid (e.g. "example.com" not "https://example.com").

## Usage

Search for the book *The Iliad* by Homer.

```sh
toshi The Iliad Homer
```

## Disclaimer

This software is provided for educational and research purposes only. The
authors do not condone or support the use of this tool to access or download
copyrighted materials without proper authorization.

By using this software, you agree to comply with all applicable laws in your
jurisdiction. The authors are not responsible for any misuse of the tool,
including but not limited to violations of copyright law.

It is the user's responsibility to ensure that their use of the software
complies with all legal requirements.
