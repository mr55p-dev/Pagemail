import { AuthCheck } from "../components/AuthCheck"
import { AccountView } from "../components/AccountView";

export default function account ({ }): JSX.Element {
    return (
        <AuthCheck>
            <AccountView />
        </AuthCheck>
     );
}