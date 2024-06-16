import { useSearchParams } from "react-router-dom";
import { validRedirect } from "./utils";

const Success = ({ email }: { email: string }) => {
    const [searchParams] = useSearchParams();
    const givenRedirectUrl = searchParams.get("redirect_url");
    const redirectUrl = validRedirect(givenRedirectUrl);

    setTimeout(() => {
        window.open(redirectUrl, "_self");
    }, 3000);

    return (
        <div className="success-container">
            <img src="/green-check.webp" alt="Success" />
            <div className="email">
                Successfully authenticated as
                <br />
                <span>{email}</span>
            </div>
            <a href={redirectUrl}>
                Redirecting to <span>{redirectUrl}</span> in a few seconds
            </a>
        </div>
    );
};

export default Success;
