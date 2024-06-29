import { useSearchParams, Link, useNavigate } from "react-router-dom";
import { validRedirect } from "./utils";
import { useEffect } from "react";

const Success = ({ email }: { email: string }) => {
    const [searchParams] = useSearchParams();
    const givenRedirectUrl = searchParams.get("redirect_url");
    const { serviceName, serviceLink } = validRedirect(givenRedirectUrl);
    console.log(serviceName, serviceLink)
    const navigate = useNavigate();

    useEffect(() => {
        const timer = setTimeout(() => {
            if (serviceLink === "/services") {
                navigate(serviceLink);
            } else {
                window.open(serviceLink, "_self");
            }
        }, 3000);

        return () => clearTimeout(timer);
    }, [serviceLink]);

    return (
        <div className="success-container">
            <img src="/green-check.webp" alt="Success" />
            <div className="email">
                Successfully authenticated as
                <br />
                <span>{email}</span>
            </div>
            <Link to={serviceLink}>
                Redirecting to <span>{serviceName}</span> in a few seconds
            </Link>
        </div>
    );
};

export default Success;
