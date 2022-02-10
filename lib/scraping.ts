export async function scrapePageMetadata(pageURL: URL, token: string) {
    let metadata;

    const apiAddress = new URL(window.location.origin)
    apiAddress.pathname = "/api/scrape";
    apiAddress.searchParams.set("url", encodeURIComponent(pageURL.toString()))

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

    const body = await resp.json()
    metadata = body;

    return metadata;
}