import { useState, useEffect } from "react";
import { validateEmail } from "./utils";
import toast from "react-hot-toast";
import { BACKEND_URL } from "./constants";
import Success from "./Success";

type FormProps = {
    isAuthenticated: boolean;
    setIsAuthenticated: React.Dispatch<React.SetStateAction<boolean>>;
};

const Form = ({ isAuthenticated, setIsAuthenticated }: FormProps) => {
    const [email, setEmail] = useState("");
    const [otpRequested, setOtpRequested] = useState(false);
    const [otp, setOtp] = useState("");
    const [timer, setTimer] = useState("01:00");
    const [timerStart, setTimerStart] = useState<number | null>(null);

    useEffect(() => {
        fetch(`${BACKEND_URL}/validate-jwt`, { credentials: "include" }).then(
            (response) => {
                if (response.ok) {
                    response.json().then((data) => {
                        setEmail(data.email);
                        setIsAuthenticated(true);
                    });
                }
            },
        );
    }, []);

    useEffect(() => {
        if (!otpRequested || timerStart === null) return;

        // const startTime = Date.now();
        const countdown = setInterval(() => {
            const totalSeconds = Math.max(
                60 - Math.floor((Date.now() - timerStart) / 1000),
                0,
            );
            const minutes = Math.floor(totalSeconds / 60);
            const seconds = totalSeconds % 60;
            setTimer(
                (minutes > 9 ? minutes : "0" + minutes) +
                    ":" +
                    (seconds > 9 ? seconds : "0" + seconds),
            );
            if (minutes === 0 && seconds === 0) {
                clearInterval(countdown);
            }
        }, 1000);
        return () => clearInterval(countdown);
    }, [otpRequested, timerStart]);

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
                        setTimerStart(Date.now());
                        setTimer("01:00");
                    } else {
                        response.text().then((msg) => {
                            // Show message to user in case he refreshes page and attempts to request otp again
                            toast.error(msg, {
                                id: loadingToast,
                            });
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
                        {timer === "00:00" ? (
                            <button onClick={requestOTP}>Resend OTP</button>
                        ) : (
                            <button
                                style={{
                                    backgroundColor: "#f2b183",
                                    cursor: "wait",
                                }}
                                disabled={true}
                            >
                                {timer}
                            </button>
                        )}
                    </>
                ) : (
                    <button onClick={requestOTP}>Send OTP</button>
                )}
            </div>
        </div>
    );
};

export default Form;
