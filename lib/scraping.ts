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