## Readability

Readability is a service for article text extraction and audio generation.
It consists of endpoints which form the contract between pagemail and readability:

1. `POST /check` - returns 200 if an article can be read or 400 otherwise
2. `POST /extract` - returns the extracted article text
3. `POST /synthesize` - begins a polly speech synthesis task and returns the job id.
4. `GET /status` - returns the status of the provided job id from polly
