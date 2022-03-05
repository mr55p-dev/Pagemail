import Head from "next/head"
import { AuthCheck } from "../components/AuthCheck"
import PagesView from "../components/pagesView"

export default function PagesRoute() {
    return(
        <main className="">
            <Head>
                <title>Saved pages</title>
            </Head>
            <AuthCheck>
                <h1 className="page-heading">Your Pages</h1>
                <PagesView />
            </AuthCheck>
        </main>

    )
}