# go-fasttld

[![Go Reference](https://img.shields.io/badge/go-reference-blue?logo=go&logoColor=white&style=for-the-badge)](https://pkg.go.dev/github.com/elliotwutingfeng/go-fasttld)
[![Go Report Card](https://goreportcard.com/badge/github.com/elliotwutingfeng/go-fasttld?style=for-the-badge)](https://goreportcard.com/report/github.com/elliotwutingfeng/go-fasttld)
[![Codecov Coverage](https://img.shields.io/codecov/c/github/elliotwutingfeng/go-fasttld?color=bright-green&logo=codecov&style=for-the-badge&token=GB00MYK51E)](https://codecov.io/gh/elliotwutingfeng/go-fasttld)
[![Mentioned in Awesome Go](https://img.shields.io/static/v1?logo=awesomelists&label=&labelColor=CCA6C4&logoColor=261120&message=Mentioned%20in%20awesome&color=494368&style=for-the-badge)](https://github.com/avelino/awesome-go)

[![GitHub license](https://img.shields.io/badge/LICENSE-BSD--3--CLAUSE-GREEN?style=for-the-badge)](LICENSE)

## Summary

**go-fasttld** is a high performance [top level domains (TLD)](https://en.wikipedia.org/wiki/Top-level_domain) extraction module that extracts subcomponents from [URLs](https://en.wikipedia.org/wiki/URL).

URLs can either contain hostnames, IPv4 addresses, or IPv6 addresses. TLD extraction is based on the [Mozilla Public Suffix List](http://www.publicsuffix.org). Private domains listed in the [Mozilla Public Suffix List](http://www.publicsuffix.org) like 'blogspot.co.uk' and 'sinaapp.com' are also supported.

![Demo](demo.gif)

Spot any bugs? Report them [here](https://github.com/elliotwutingfeng/go-fasttld/issues)

## Installation

```sh
go get github.com/elliotwutingfeng/go-fasttld
```

## Try the CLI

First, build the CLI application.

```sh
# `git clone` and `cd` to the go-fasttld repository folder first
make build_cli
```

Afterwards, try extracting subcomponents from a URL.

```sh
# `git clone` and `cd` to the go-fasttld repository folder first
./dist/fasttld extract https://user@a.subdomain.example.ac.uk:5000/a/b\?id\=42
```

## Try the example code

All of the following examples can be found at `examples/demo.go`. To play the demo, run the following command:

```sh
# `git clone` and `cd` to the go-fasttld repository folder first
make demo
```

### Hostname

```go
// Initialise fasttld extractor
extractor, _ := fasttld.New(fasttld.SuffixListParams{})

// Extract URL subcomponents
url := "https://user@a.subdomain.example.ac.uk:5000/a/b?id=42"
res, _ := extractor.Extract(fasttld.URLParams{URL: url})

// Display results
fasttld.PrintRes(url, res) // Pretty-prints res.Scheme, res.UserInfo, res.SubDomain etc.
```

| Scheme   | UserInfo | SubDomain   | Domain  | Suffix | RegisteredDomain | Port | Path       | HostType     |
|----------|----------|-------------|---------|--------|------------------|------|------------|--------------|
| https:// | user     | a.subdomain | example | ac.uk  | example.ac.uk    | 5000 | /a/b?id=42 | hostname     |

### IPv4 Address

```go
extractor, _ := fasttld.New(fasttld.SuffixListParams{})
url := "https://127.0.0.1:5000"
res, _ := extractor.Extract(fasttld.URLParams{URL: url})
```

| Scheme   | UserInfo | SubDomain | Domain    | Suffix | RegisteredDomain | Port | Path | HostType     |
|----------|----------|-----------|-----------|--------|------------------|------|------|--------------|
| https:// |          |           | 127.0.0.1 |        | 127.0.0.1        | 5000 |      | ipv4 address |

### IPv6 Address

```go
extractor, _ := fasttld.New(fasttld.SuffixListParams{})
url := "https://[aBcD:ef01:2345:6789:aBcD:ef01:2345:6789]:5000"
res, _ := extractor.Extract(fasttld.URLParams{URL: url})
```

| Scheme   | UserInfo | SubDomain | Domain                                  | Suffix | RegisteredDomain                        | Port | Path | HostType     |
|----------|----------|-----------|-----------------------------------------|--------|-----------------------------------------|------|------|--------------|
| https:// |          |           | aBcD:ef01:2345:6789:aBcD:ef01:2345:6789 |        | aBcD:ef01:2345:6789:aBcD:ef01:2345:6789 | 5000 |      | ipv6 address |

### Internationalised label separators

**go-fasttld** supports the following internationalised label separators (IETF RFC 3490)

| Full Stop  | Ideographic Full Stop | Fullwidth Full Stop | Halfwidth Ideographic Full Stop |
|------------|-----------------------|---------------------|---------------------------------|
| U+002E `.` | U+3002 `。`           | U+FF0E `．`         | U+FF61 `｡`                      |

```go
extractor, _ := fasttld.New(fasttld.SuffixListParams{})
url := "https://brb\u002ei\u3002am\uff0egoing\uff61to\uff0ebe\u3002a\uff61fk"
res, _ := extractor.Extract(fasttld.URLParams{URL: url})
```

| Scheme   | UserInfo | SubDomain                             | Domain | Suffix    | RegisteredDomain  | Port | Path | HostType     |
|----------|----------|---------------------------------------|--------|-----------|-------------------|------|------|--------------|
| https:// |          | brb\u002ei\u3002am\uff0egoing\uff61to | be     | a\uff61fk | be\u3002a\uff61fk |      |      | hostname     |

## Public Suffix List options

### Specify custom public suffix list file

You can use a custom public suffix list file by setting `CacheFilePath` in `fasttld.SuffixListParams{}` to its absolute path.

```go
cacheFilePath := "/absolute/path/to/file.dat"
extractor, _ := fasttld.New(fasttld.SuffixListParams{CacheFilePath: cacheFilePath})
```

### Updating the default Public Suffix List cache

Whenever `fasttld.New` is called without specifying `CacheFilePath` in `fasttld.SuffixListParams{}`, the local cache of the default Public Suffix List is updated automatically if it is more than 3 days old. You can also manually update the cache by using `Update()`.

```go
// Automatic update performed if `CacheFilePath` is not specified
// and local cache is more than 3 days old
extractor, _ := fasttld.New(fasttld.SuffixListParams{})

// Manually update local cache
if err := extractor.Update(); err != nil {
    log.Println(err)
}
```

### Private domains

According to the [Mozilla.org wiki](https://wiki.mozilla.org/Public_Suffix_List/Uses), the Mozilla Public Suffix List contains private domains like `blogspot.com` and `sinaapp.com`.

By default, these private domains are excluded (i.e. `IncludePrivateSuffix = false`)

```go
extractor, _ := fasttld.New(fasttld.SuffixListParams{})
url := "https://google.blogspot.com"
res, _ := extractor.Extract(fasttld.URLParams{URL: url})
```

| Scheme   | UserInfo | SubDomain | Domain   | Suffix | RegisteredDomain | Port | Path | HostType     |
|----------|----------|-----------|----------|--------|------------------|------|------|--------------|
| https:// |          | google    | blogspot | com    | blogspot.com     |      |      | hostname     |

You can _include_ private domains by setting `IncludePrivateSuffix = true`

```go
extractor, _ := fasttld.New(fasttld.SuffixListParams{IncludePrivateSuffix: true})
url := "https://google.blogspot.com"
res, _ := extractor.Extract(fasttld.URLParams{URL: url})
```

| Scheme   | UserInfo | SubDomain | Domain | Suffix       | RegisteredDomain    | Port | Path | HostType     |
|----------|----------|-----------|--------|--------------|---------------------|------|------|--------------|
| https:// |          |           | google | blogspot.com | google.blogspot.com |      |      | hostname     |

## Extraction options

### Ignore Subdomains

You can ignore subdomains by setting `IgnoreSubDomains = true`. By default, subdomains are extracted.

```go
extractor, _ := fasttld.New(fasttld.SuffixListParams{})
url := "https://maps.google.com"
res, _ := extractor.Extract(fasttld.URLParams{URL: url, IgnoreSubDomains: true})
```

| Scheme   | UserInfo | SubDomain | Domain | Suffix | RegisteredDomain | Port | Path | HostType     |
|----------|----------|-----------|--------|--------|------------------|------|------|--------------|
| https:// |          |           | google | com    | google.com       |      |      | hostname     |

### Punycode

By default, internationalised URLs are not converted to punycode before extraction.

```go
extractor, _ := fasttld.New(fasttld.SuffixListParams{})
url := "https://hello.世界.com"
res, _ := extractor.Extract(fasttld.URLParams{URL: url})
```

| Scheme   | UserInfo | SubDomain | Domain | Suffix | RegisteredDomain | Port | Path | HostType     |
|----------|----------|-----------|--------|--------|------------------|------|------|--------------|
| https:// |          | hello     | 世界   | com    | 世界.com         |      |      | hostname     |

You can convert internationalised URLs to [punycode](https://en.wikipedia.org/wiki/Punycode) before extraction by setting `ConvertURLToPunyCode = true`.

```go
extractor, _ := fasttld.New(fasttld.SuffixListParams{})
url := "https://hello.世界.com"
res, _ := extractor.Extract(fasttld.URLParams{URL: url, ConvertURLToPunyCode: true})
```

| Scheme   | UserInfo | SubDomain | Domain      | Suffix | RegisteredDomain | Port | Path | HostType     |
|----------|----------|-----------|-------------|--------|------------------|------|------|--------------|
| https:// |          | hello     | xn--rhqv96g | com    | xn--rhqv96g.com  |      |      | hostname     |

## Parsing errors

If the URL is invalid, the second variable returned by `Extract()`, **error**, will be non-nil. Partially extracted subcomponents can still be retrieved from the first variable returned, ***ExtractResult**.

```go
extractor, _ := fasttld.New(fasttld.SuffixListParams{})
url := "https://example!.com" // invalid characters in hostname
color.New().Println("The following line should be an error message")
if res, err := extractor.Extract(fasttld.URLParams{URL: url}); err != nil {
    color.New(color.FgHiRed, color.Bold).Print("Error: ")
    color.New(color.FgHiWhite).Println(err)
}
fasttld.PrintRes(url, res) // Partially extracted subcomponents can still be retrieved
```

| Scheme   | UserInfo | SubDomain | Domain | Suffix | RegisteredDomain | Port | Path | HostType |
|----------|----------|-----------|--------|--------|------------------|------|------|----------|
| https:// |          |           |        |        |                  |      |      |          |

## Testing

```sh
# `git clone` and `cd` to the go-fasttld repository folder first
make tests
```

## Benchmarks

```sh
# `git clone` and `cd` to the go-fasttld repository folder first
make bench
```

### Modules used

| Benchmark Name       | Source                           |
|----------------------|----------------------------------|
| GoFastTld            | go-fasttld (this module)         |
| JPilloraGoTld        | github.com/jpillora/go-tld       |
| JoeGuoTldExtract     | github.com/joeguo/tldextract     |
| Mjd2021USATldExtract | github.com/mjd2021usa/tldextract |
| M507Tlde             | github.com/M507/tlde             |

### Results

Benchmarks performed on AMD Ryzen 7 5800X, Manjaro Linux.

**go-fasttld** performs especially well on longer URLs.

---

#### #1

<code>https://news.google.com</code>

| Benchmark Name       | Iterations | ns/op       | B/op      | allocs/op   | Fastest            |
|----------------------|------------|-------------|-----------|-------------|--------------------|
| GoFastTld            | 5335618    | 216.5 ns/op | 144 B/op  | 1 allocs/op | :heavy_check_mark: |
| JPilloraGoTld        | 2408444    | 477.9 ns/op | 224 B/op  | 2 allocs/op |                    |
| JoeGuoTldExtract     | 2506902    | 474.2 ns/op | 160 B/op  | 5 allocs/op |                    |
| Mjd2021USATldExtract | 1519467    | 794.0 ns/op | 208 B/op  | 7 allocs/op |                    |
| M507Tlde             | 2536875    | 494.0 ns/op | 160 B/op  | 5 allocs/op |                    |

---

#### #2

<code>https://iupac.org/iupac-announces-the-2021-top-ten-emerging-technologies-in-chemistry/</code>

| Benchmark Name       | Iterations | ns/op       | B/op      | allocs/op   | Fastest            |
|----------------------|------------|-------------|-----------|-------------|--------------------|
| GoFastTld            | 5740478    | 196.9 ns/op | 144 B/op  | 1 allocs/op | :heavy_check_mark: |
| JPilloraGoTld        | 1735533    | 683.9 ns/op | 224 B/op  | 2 allocs/op |                    |
| JoeGuoTldExtract     | 2343976    | 511.0 ns/op | 272 B/op  | 5 allocs/op |                    |
| Mjd2021USATldExtract | 1592064    | 741.0 ns/op | 288 B/op  | 6 allocs/op |                    |
| M507Tlde             | 2295045    | 517.7 ns/op | 272 B/op  | 5 allocs/op |                    |

---

#### #3

<code>https://www.google.com/maps/dir/Parliament+Place,+Parliament+House+Of+Singapore,+Singapore/Parliament+St,+London,+UK/@25.2440033,33.6721455,4z/data=!3m1!4b1!4m14!4m13!1m5!1m1!1s0x31da19a0abd4d71d:0xeda26636dc4ea1dc!2m2!1d103.8504863!2d1.2891543!1m5!1m1!1s0x487604c5aaa7da5b:0xf13a2197d7e7dd26!2m2!1d-0.1260826!2d51.5017061!3e4</code>

| Benchmark Name       | Iterations | ns/op       | B/op      | allocs/op   | Fastest            |
|----------------------|------------|-------------|-----------|-------------|--------------------|
| GoFastTld            | 4744010    | 229.7 ns/op | 144 B/op  | 1 allocs/op | :heavy_check_mark: |
| JPilloraGoTld        | 448917     | 2826 ns/op  | 928 B/op  | 4 allocs/op |                    |
| JoeGuoTldExtract     | 909306     | 1283 ns/op  | 1120 B/op | 6 allocs/op |                    |
| Mjd2021USATldExtract | 1000000    | 1297 ns/op  | 1120 B/op | 6 allocs/op |                    |
| M507Tlde             | 1000000    | 1257 ns/op  | 1120 B/op | 6 allocs/op |                    |

---

## Implementation details

### Why not split on "." and take the last element instead?

Splitting on "." and taking the last element only works for simple TLDs like `com`, but not more complex ones like `oseto.nagasaki.jp`.

### TLD tries

![Trie](Trie_example.svg)

**go-fasttld** stores TLDs in [compressed tries](https://en.wikipedia.org/wiki/Trie).

Valid TLDs from the [Mozilla Public Suffix List](http://www.publicsuffix.org) are appended to the compressed trie in reverse-order.

```sh
Given the following TLDs
au
nsw.edu.au
com.ac
edu.ac
gov.ac

and the example URL host `example.nsw.edu.au`

The compressed trie will be structured as follows:

START
 ╠═ au 🚩 ✅
 ║  ╚═ edu ✅
 ║     ╚═ nsw 🚩 ✅
 ╚═ ac
    ╠═ com 🚩
    ╠═ edu 🚩
    ╚═ gov 🚩

=== Symbol meanings ===
🚩 : path to this node is a valid TLD
✅ : path to this node found in example URL host `example.nsw.edu.au`
```

The URL host subcomponents are parsed from right-to-left until no more matching nodes can be found. In this example, the path of matching nodes are `au -> edu -> nsw`. Reversing the nodes gives the extracted TLD `nsw.edu.au`.

## Acknowledgements

This module is a port of the Python [fasttld](https://github.com/jophy/fasttld) module, with additional modifications to support extraction of subcomponents from full URLs, IPv4 addresses, and IPv6 addresses.

- [fasttld (Python)](https://github.com/jophy/fasttld)
- [tldextract (Python)](https://github.com/john-kurkowski/tldextract)
- [ICANN IDN Character Validation Guidance](https://www.icann.org/resources/pages/idna-protocol-2012-02-25-en)
- [IETF RFC 2396](https://www.ietf.org/rfc/rfc2396.txt)
- [IETF RFC 3490](https://www.ietf.org/rfc/rfc3490.txt)
- [IETF RFC 3986](https://www.ietf.org/rfc/rfc3986.txt)
