import { useState } from "react";

function App() {
    const [otpRequested, setOtpRequested] = useState(false);

    const requestOTP = () => {
        setOtpRequested(true);
    };

    return (
        <div className="form-container">
            <div className="info">
                <div className="title">Heimdall</div>
                <p>The gatekeeper to Metakgp services</p>
                <p>Please verify using your kgpian email to continue</p>
            </div>
            <div className="form">
                <input type="email" placeholder="Enter your kgpian email" />
                {otpRequested ? (
                    <>
                        <input type="text" placeholder="Enter OTP" />
                        <button>Verify</button>
                    </>
                ) : (
                    <button onClick={requestOTP}>Send OTP</button>
                )}
            </div>
        </div>
    );
}

export default App;
