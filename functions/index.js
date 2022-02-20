const nodemailer = require("nodemailer");
const admin = require('firebase-admin');
const functions = require('firebase-functions');
const { getApps } = require('firebase-admin/app');

async function sendMail({ email, content }) {
    // Addresses is string
    let transporter = nodemailer.createTransport({
        host: process.env.MAIL_HOST,
        port: process.env.MAIL_PORT,
        secure: false, // true for 465, false for other ports
        auth: {
          user: process.env.MAIL_USER,
          pass: process.env.MAIL_PASS
        },
    });

    // send mail with defined transport object
    let info = await transporter.sendMail({
        from: `"PageMail Roundup" <${process.env.MAIL_USER}>`, // sender address
        to: email, // list of receivers
        subject: "Your daily roundup for $today$", // Subject line
        text: "You need HTML to display this email", // plain text body
        html: `<div>${content}</div>`, // html body
    });
    functions.logger.info("Message sent: %s", info.messageId);
    return email;
}

async function contactUser(document) {
    // Get the user properties
    // const contents = document.data();
    functions.logger.debug("Contacting user %s", document.id)
    const db = admin.firestore()

    // Get the user ID and email
    const uid = document.id;
    const email = document.data().email;

    // Create a reference to the specific pages
    const pagesRef = db.collection(`users/${uid}/pages`)
    const yesterday = new Date(new Date().getTime() - (24 * 60 * 60 * 1000));
    const yesterdayFS = admin.firestore.Timestamp.fromDate(yesterday)

    // Get the pages added since yesterday
    const pageSnapshot = await pagesRef.where('timeAdded', '>=', yesterdayFS).limit(10).get()

    // Unwrap all the pages into the page data
    const pages = pageSnapshot.docs.map(doc => doc.data())
    if (pages.empty) {
        return Promise.reject("No pages saved...")
    }

    // Make a html list string of the titles
    const listItems = pages.map( page => `<li>${page.url}</li>` )
    const listString = listItems.join("\n")

    functions.logger.debug(listString)

    // Construct the HTML
    const emailBody = `<body>
        <h1>PageMail daily roundup</h1>
        <ul>${listString}</ul>
    </body>`

    // Send the email!
    return sendMail({email, content: emailBody})

}

exports.helloWorld = functions.https.onRequest(async (req, res) => {
    // Only get the app once
    if (getApps().length === 0) {
        admin.initializeApp(functions.config().firebase);
    }

    // Create an instance of the database object using the admin config
    const db = admin.firestore()

    // Get a reference to the users collection
    const usersRef = db.collection("users")

    // Get all the users which have the newsletter property set true
    const subscribed = await usersRef.where('newsletter', '==', true).get()

    if (subscribed.empty) {
        return res.json({error: "No users are subscribed"})
    }

    // For every user, contact them!
    const mails = subscribed.docs.map(contactUser)

    // Try to settle all these promises
    try {
        const promises = await Promise.allSettled(mails)
        return res.json({msg: promises})
    } catch (error) {
        return res.json({error})
    }
})