import { Toaster } from "react-hot-toast";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import Form from "./Form";

function App() {
    return (
        <>
            <BrowserRouter>
                <Routes>
                    <Route path="/" element={<Form />} />
                </Routes>
            </BrowserRouter>
            <Toaster position="bottom-center" reverseOrder={false} />
        </>
    );
}

export default App;
