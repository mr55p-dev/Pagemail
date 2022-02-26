// Next.js API route support: https://nextjs.org/docs/api-routes/introduction

import { load } from 'cheerio';
import { getAuth } from 'firebase-admin/auth';
import { getApps } from 'firebase-admin/app';
import { IPageMetadata } from '../../lib/typeAliases';

// Must be done this way due to library import requirements
const admin = require('firebase-admin');


// Route logic
export default catchErrorsFrom(async (req, res) => {
  console.debug("Inbound request")
  // Verify the user
  const uid = await verifyUser(req.headers?.token)

  console.debug("Request from UID %s accepted", uid)
  // Verify the url
  const url = await verifyURL(req.query?.url)

  console.debug("URL %s accepted", url.toString())
  // Fetch the page metadata
  const meta = await scrapeMeta(url);

  console.debug("Metadata fetched with title: %s", meta?.title)
  return res.status(200).json(meta);
})


// Error handler
function catchErrorsFrom(handler) {
  return async (req, res) => {
    return handler(req, res)
      .catch((error) => {
        // Error handling!
        console.error(error)
        if (error === "notSignedIn") {
          return res.status(401).end("Unauthorised")
        } else if (error === "invalidURL") {
          return res.status(400).end("Invalid URL");
        } else if (error === "badGateway") {
          return res.status(200).json({})
        }
        return res.status(500).send(error.message || error);
      });
  }
}


// URL verification
async function verifyURL(encodedURL: string): Promise<URL> {
  let url: URL;
  try {
    return new URL(decodeURIComponent(encodedURL))
  } catch {
    throw "invalidURL"
  }
}

// User verification
async function verifyUser(token: string): Promise<string> {
  // Initialise firebase
  if (getApps().length === 0) {
    const serviceAccount = JSON.parse(
      process.env.FIREBASE_ADMIN_ACCOUNT_KEY
    );

    admin.initializeApp({
      credential: admin.credential.cert(serviceAccount)
    });
  }

  // Initialise auth
  const auth  = getAuth();

  // Reject the promise if there is no token
  if (!token) { throw "notSignedIn" }

  // Use firebase admin to verify the token and return the user id
  const decodedToken = await auth.verifyIdToken(token)
  return decodedToken.uid
}

// Metadata extraction
async function scrapeMeta(url: URL): Promise<IPageMetadata> {
  // Make the request and check the response
  const pageResponse = await fetch(url.toString());
  if (!pageResponse.ok) { throw "badGateway" }

  // Load the document
  const $ = load(await pageResponse.text());

  // Define the scrape function
  const getMetatag = (name: string) =>
    $(`meta[name=${name}]`).attr('content') ||
    $(`meta[property="og:${name}"]`).attr('content') ||
    $(`meta[property="twitter:${name}"]`).attr('content');

  // Scrape!
  const scrapedDoc: IPageMetadata = {
    url: url.toString(),
    title: $("title").first().text(),
    description: getMetatag("description"),
    author: getMetatag("author"),
    image: getMetatag("image")
  }
  return scrapedDoc;
}