# httpsrescheck
An experimental project that is aimed to find insecure resources on HTTPS sites. It is a console tool written in Go.

## Installation

```
go get github.com/kkomelin/httpsrescheck
```

## Usage

```
./httpsrescheck https://example.com
```

## Roadmap

- Check IFRAME and OBJECT tags
- Implement debug option and hide url log by default
- Improve output format (CSV?)
- Handle trailing slash (https://example.com and https://example.com/)
- Handle domains w/ and w/o WWW
