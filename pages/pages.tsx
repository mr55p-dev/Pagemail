import { AuthCheck } from "../components/AuthCheck"
import PagesView from "../components/pagesView"

export default function PagesRoute() {
    return(
        <main>
            <AuthCheck>
                <h1 className="text-center text-3xl underline my-3">Your Pages</h1>
                <PagesView />
            </AuthCheck>
        </main>

    )
}