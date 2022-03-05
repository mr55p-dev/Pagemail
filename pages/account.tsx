import Head from "next/head";
import { useContext } from "react"
import { AuthCheck } from "../components/AuthCheck"
import { UserContext } from "../lib/context"
import { useUserData } from "../lib/hooks";

export default function Account ({ }): JSX.Element {
    const { email, username, newsletter } = useUserData();

    return (
        <AuthCheck>
            <main className="p-3">
                <Head>
                    <title>Your account</title>
                </Head>
                <h1 className="page-heading">Account information</h1>
                <div>
                    <form className="grid grid-rows-5 grid-cols-12 gap-2">
                        <p className="col-span-12 md:col-span-4">Username: </p>
                            <input className="col-span-12 md:col-span-8" defaultValue={username} readOnly={true} />
                        <p className="col-span-12 md:col-span-4">Email address: </p>
                            <input className="col-span-12 md:col-span-8" defaultValue={email} readOnly={true} />
                        <p className="col-span-10 md:col-span-4">Subscribe to emails: </p>
                            <input className=" col-span-2 md:col-span-8" type="checkbox" defaultChecked={newsletter} />
                    </form>
                </div>
            </main>
        </AuthCheck>
     );
}