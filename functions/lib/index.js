"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.onUserCreation = void 0;
const functions = require("firebase-functions");
const app_1 = require("firebase-admin/app");
const firestore_1 = require("firebase-admin/firestore");
exports.onUserCreation = functions.auth.user().onCreate((user) => {
    functions.logger.info(`New user signup ${user.displayName} (${user.email})`);
    // Send welcome email
    // Write to the users table
    const app = app_1.initializeApp();
    const db = firestore_1.getFirestore(app);
    const userDoc = db.collection("users").doc(`${user.uid}`);
    userDoc.set({
        name: user.displayName,
        uid: user.uid,
    });
});
//# sourceMappingURL=index.js.map