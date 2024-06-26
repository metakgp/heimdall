import { useState } from "react";
import { validateEmail } from "./utils";
import toast from "react-hot-toast";
import { BACKEND_URL } from "./constants";
import Success from "./Success";

type FormProps = {
    isAuthenticated: boolean
    setIsAuthenticated: React.Dispatch<React.SetStateAction<boolean>>
    email: string
    setEmail: React.Dispatch<React.SetStateAction<string>>
}

const Form = ({isAuthenticated, setIsAuthenticated, email, setEmail }: FormProps) => {
    const [otpRequested, setOtpRequested] = useState(false);
    const [otp, setOtp] = useState("");


    const requestOTP = () => {
        if (validateEmail(email)) {
            const loadingToast = toast.loading("Sending OTP...");
            const formdata = new FormData();
            formdata.append("email", email);
            fetch(`${BACKEND_URL}/get-otp`, {
                method: "POST",
                body: formdata,
            })
                .then((response) => {
                    if (response.ok) {
                        toast.success("OTP sent successfully", {
                            id: loadingToast,
                        });
                        setOtpRequested(true);
                    } else {
                        toast.error("Failed to send OTP. Please try again", {
                            id: loadingToast,
                        });
                    }
                })
                .catch(() => {
                    toast.error("Failed to send OTP. Please try again", {
                        id: loadingToast,
                    });
                });
        } else {
            toast.error("Invalid email address. Please use your kgpian email");
            return;
        }
    };

    const verifyOTP = () => {
        const loadingToast = toast.loading("Verifying OTP...");
        const formdata = new FormData();
        formdata.append("email", email);
        formdata.append("otp", otp);
        fetch(`${BACKEND_URL}/verify-otp`, {
            method: "POST",
            body: formdata,
            credentials: "include",
        })
            .then((response) => {
                if (response.ok) {
                    toast.success("OTP verified successfully", {
                        id: loadingToast,
                    });
                    setIsAuthenticated(true);
                } else {
                    toast.error("Failed to verify OTP. Please try again", {
                        id: loadingToast,
                    });
                }
            })
            .catch(() => {
                toast.error("Failed to verify OTP. Please try again", {
                    id: loadingToast,
                });
            });
    };

    if (isAuthenticated) {
        return <Success email={email} />;
    }

    return (
        <div className="form-container">
            <div className="info">
                <div className="title">Heimdall</div>
                <p>The gatekeeper to Metakgp services</p>
                <p>Please verify using your kgpian email to continue</p>
            </div>
            <div className="form">
                <input
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="Enter your kgpian email"
                />
                {otpRequested ? (
                    <>
                        <input
                            type="text"
                            value={otp}
                            onChange={(e) => setOtp(e.target.value)}
                            placeholder="Enter OTP"
                        />
                        <button onClick={verifyOTP}>Verify</button>
                    </>
                ) : (
                    <button onClick={requestOTP}>Send OTP</button>
                )}
            </div>
        </div>
    );
};

export default Form;
