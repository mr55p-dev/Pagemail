import Head from "next/head";
import { useContext } from "react"
import { AuthCheck } from "../components/AuthCheck"
import { UserContext } from "../lib/context"
import { useUserData } from "../lib/hooks";

export default function Account ({ }): JSX.Element {
    const { email, username, newsletter } = useUserData();

    return (
        <AuthCheck>
            <main>
                <Head>
                    <title>Your account</title>
                </Head>
                <h1>Your account information</h1>
                <div>
                    <form>
                        <p>Username: </p><input defaultValue={username} readOnly={true} />
                        <p>email: </p><input defaultValue={email} readOnly={true} />
                        <p>Subscribe to emails: </p><input type="checkbox" defaultChecked={newsletter} />
                    </form>
                </div>
            </main>
        </AuthCheck>
     );
}