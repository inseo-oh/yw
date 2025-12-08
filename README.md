# YJK's silly web adventures (YW)

I've first started learning the Web from the high-level - Writing HTML, CSS, JS, and using various Web APIs.

But how about starting from the... other end?

## Warning

This web browser does NOT implement any web security yet (CORS, CSRF protection, etc). You've been warned.

## Current progress

Everything is in WIP and very early stage. For now, this repo is mostly for backing up my source code.

There are also some "tools" inside cmd directory, but they were made for experimenting, and will likely get removed in the future.

## Getting started

If you want to try this out yourself, run:

```
go run ./cmd/yw -url <HTTP/HTTPS URL of the page>
```

Or to build as executable:

```
go build ./cmd/yw
```

## License

-   See `LICENSE` for license of YW.
-   See `LICENSE_WHATWG_SPECS` for details about WHATWG software licenses. In short, all included soruce code portions are BSD-3-Clause.
