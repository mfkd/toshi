# Toshi - A Simple Book Search Tool

Toshi is a command-line tool designed for searching books.

"Toshi" (とし / 知 or 智) is a Japanese word that translates to "wisdom" or
"knowledge."

## Installation

Install the latest tagged release with Go:

```sh
go install github.com/mfkd/toshi@latest
```

To build and update `toshi` from a local checkout, run `go install .` from the
repository root.

## Configuration

The default catalog is `libgen.bz`. Set `DOMAIN` to override it:

```sh
export DOMAIN="libgen.bz,libgen.vg"
```

Use hostnames without a scheme or path. Toshi tries the comma-separated
sources in order until it finds EPUB results. Toshi uses `/index.php` for
`libgen.bz`, `libgen.gl`, `libgen.la`, `libgen.li`, and `libgen.vg`.
Other hostnames use the legacy `/search.php` catalog.

## Usage

Search for the book *The Iliad* by Homer.

```sh
toshi The Iliad Homer
```

Select an EPUB by number. Toshi saves the download in `output/` under the
current working directory. If the filename already exists, Toshi adds a
numbered suffix. Use `-v` for diagnostics and `--help` for usage.

## Disclaimer

This software is provided for educational and research purposes only. The
authors do not condone or support the use of this tool to access or download
copyrighted materials without proper authorization.

By using this software, you agree to comply with all applicable laws in your
jurisdiction. The authors are not responsible for any misuse of the tool,
including but not limited to violations of copyright law.

It is the user's responsibility to ensure that their use of the software
complies with all legal requirements.
