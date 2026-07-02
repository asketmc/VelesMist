# API Security

VelesMist uses read-only Steam Community public inventory endpoints.

Rules:

- no authenticated Steam endpoints;
- no API keys, cookies, sessions, passwords, or Steam Guard codes;
- no POST, PUT, PATCH, or DELETE requests to Steam services;
- bounded HTTP response reads;
- HTTP client timeout is mandatory;
- upstream failures are typed and mapped to exit code `3`.

Any new network target, authenticated endpoint, write method, or third-party pricing provider requires security review.
