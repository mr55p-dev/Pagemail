// Next.js API route support: https://nextjs.org/docs/api-routes/introduction

import { load } from 'cheerio';
import { getAuth } from 'firebase-admin/auth';
import { applicationDefault, getApps } from 'firebase-admin/app';
import { PageMetadata } from '../../lib/typeAliases';


// Error handler
function catchErrorsFrom(handler) {
  return async (req, res) => {
    return handler(req, res)
      .catch((error) => {
        // Error handling!
        if (error === "notSignedIn") {
          return res.status(401).end("Unauthorised")
        } else if (error === "invalidURL") {
          return res.status(400).end("Invalid URL");
        } else if (error === "badGateway") {
          return res.status(200).json({})
        }
        console.error(error);
        return res.status(500).send(error.message || error);
      });
  }
}

// Must be done this way due to library import requirements
const admin = require('firebase-admin');

export default catchErrorsFrom(async (req, res) => {
  // Verify the user
  const uid = await verifyUser(req.headers?.token)

  // Verify the url
  const url = await verifyURL(req.query?.url)

  // Fetch the page metadata
  const meta = await scrapeMeta(url);
  return res.status(200).json(meta);
})

async function verifyURL(encodedURL: string): Promise<URL> {
  let url: URL;
  try {
    return new URL(decodeURIComponent(encodedURL))
  } catch {
    throw "invalidURL"
  }
}

async function verifyUser(token: string): Promise<string> {
  // Initialise firebase
  if (getApps().length === 0) {
    admin.initializeApp({
      credential: applicationDefault()
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

async function scrapeMeta(url: URL): Promise<PageMetadata> {
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
  const scrapedDoc: PageMetadata = {
    title: $("title").first().text(),
    description: getMetatag("description"),
    author: getMetatag("author"),
    image: getMetatag("image")
  }
  return scrapedDoc;
}