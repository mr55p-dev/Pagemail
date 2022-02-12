export function validateURL(inputString: string): URL {
    // First coerce the input into a URL
    if (!(inputString.startsWith("https://") || inputString.startsWith("http://"))) {
        console.log("Prepending https://")
        inputString = `https://${inputString}`;
    }
    console.log(`Constructing URL from ${inputString}`)
    const inputURL = new URL(inputString)

    console.log(inputURL)

    // Allow only http or https
    if (!["http:", "https:"].includes(inputURL.protocol)) {
        throw "Invalid URL protocol"
    }
    if (!(/^\S+\.\S+$/.test(inputURL.hostname))) {
        throw "Invalid URL hostname"
    }

    return inputURL
}

export async function scrapePageMetadata(pageURL: URL, token: string) {

    // Set the API to point to the widnwo origin
    const apiAddress = new URL(window.location.origin)

    // Modify the path and query parameters
    apiAddress.pathname = "/api/scrape";
    apiAddress.searchParams.set("url", encodeURIComponent(pageURL.toString()))

    // Get a response
    const resp = await fetch(apiAddress.toString(), {
        method: "GET",
        mode: "same-origin",
        credentials: "same-origin",
        headers: {
            token: token
        }
    })

    if (!resp.ok) {
        throw Error(resp.statusText)
    }

    return await resp.json();
}