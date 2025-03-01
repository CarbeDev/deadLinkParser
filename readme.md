# Dead Link Parser

Dead Link Parser is a Go-based CLI tool that crawls websites to detect broken links (dead links). It recursively checks all internal links and their HTTP status codes, helping you maintain the integrity of your website.

## Features

- Concurrent link checking using goroutines
- Handles both internal and external links
- Real-time status reporting for each checked link
- Automatic base URL detection
- Memory-efficient link tracking to avoid duplicate checks

## Installation

To build the project, make sure you have Go installed on your system, then run:

```bash
go build
```

## Usage

Run the parser by providing the target website URL:

```bash
./deadLinkParser https://example.com/
```

The tool will output the status of each link it finds:
- ✅ for successful responses (HTTP 200-399)
- ❌ for failed responses (HTTP 400+) or connection errors

Example output:
```
Link : https://example.com/about | Status : 200 OK ✅
Link : https://example.com/broken-link | Status : 404 Not Found ❌
Link : https://example.com/contact | Status : 200 OK ✅
```

## Upcoming Features

### Report Generation
We plan to add comprehensive report generation capabilities:
- HTML report export
- Summary statistics (total links, broken links, response time)
- Detailed error categorization
- Path visualization for broken link discovery
- Export options (JSON, CSV, PDF)

### Advanced Configuration
- Custom timeout settings
- Rate limiting for external requests
- Allow/deny list for domains
- Custom HTTP headers
- Proxy support
- Authentication handling for protected pages
- Custom User-Agent configuration

## Contributing

Contributions are welcome! Feel free to submit issues and pull requests.
