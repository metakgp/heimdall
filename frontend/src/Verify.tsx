import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { BACKEND_URL } from "./constants";
import toast from "react-hot-toast";
import Success from "./Success";

const Verify = () => {
    const [email, setEmail] = useState("");

    const navigate = useNavigate();
    const [searchParams] = useSearchParams();
    const rootRedirectUrl = `/?${searchParams.toString()}`;

    useEffect(() => {
        const jwt = localStorage.getItem("jwt");
        if (!jwt) {
            navigate(rootRedirectUrl);
        }

        const handleInvalidJwt = () => {
            toast.error("Invalid JWT. Please try again");
            localStorage.removeItem("jwt");
            navigate(rootRedirectUrl);
        };

        fetch(`${BACKEND_URL}/validate-jwt`, {
            headers: {
                Authorization: `Bearer ${jwt}`,
            },
        })
            .then((response) =>
                response
                    .json()
                    .then((data) => {
                        if (data.email) {
                            setEmail(data.email);
                        } else {
                            handleInvalidJwt();
                        }
                    })
                    .catch(handleInvalidJwt),
            )
            .catch(() => {
                toast.error("Failed to validate JWT. Please try again");
                navigate(rootRedirectUrl);
            });
    }, []);

    if (email === "") return <div className="verifying">Verifying</div>;

    return <Success email={email} />;
};

export default Verify;
