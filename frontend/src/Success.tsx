import { useSearchParams, Link, useNavigate } from "react-router-dom";
import { validRedirect } from "./utils";
import { useEffect } from "react";

const Success = ({ email }: { email: string }) => {
    const [searchParams] = useSearchParams();
    const givenRedirectUrl = searchParams.get("redirect_url");
    const redirectUrl = validRedirect(givenRedirectUrl);
    const navigate = useNavigate();

    useEffect(() => {
        const timer = setTimeout(() => {
            if (redirectUrl === "/services") {
                navigate(redirectUrl);
            } else {
                window.open(redirectUrl, "_self");
            }
        }, 3000);

        return () => clearTimeout(timer)
    }, [redirectUrl])

    return (
        <div className="success-container">
            <img src="/green-check.webp" alt="Success" />
            <div className="email">
                Successfully authenticated as
                <br />
                <span>{email}</span>
            </div>
            <Link to={redirectUrl}>
                Redirecting to <span>{redirectUrl}</span> in a few seconds
            </Link>
        </div>
    );
};

export default Success;
