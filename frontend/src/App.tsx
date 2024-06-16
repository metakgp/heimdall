import { Toaster } from "react-hot-toast";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import Form from "./Form";
import Services from "./Services";
import Verify from "./Verify";

function App() {
    return (
        <>
            <BrowserRouter>
                <Routes>
                    <Route path="/" element={<Form />} />
                    <Route path="/verify" element={<Verify />} />
                    <Route path="/services" element={<Services />} />
                </Routes>
            </BrowserRouter>
            <Toaster position="bottom-center" reverseOrder={false} />
        </>
    );
}

export default App;
