import { Toaster } from "react-hot-toast";
import { BrowserRouter, Route, Routes, Navigate } from "react-router-dom";
import Form from "./Form";
import Services from "./Services";
import { useState } from "react";

function App() {
    const [isAuthenticated, setIsAuthenticated] = useState(false);

    return (
        <>
            <BrowserRouter>
                <Routes>
                    <Route
                        path="/"
                        element={
                            <Form
                                isAuthenticated={isAuthenticated}
                                setIsAuthenticated={setIsAuthenticated}
                            />
                        }
                    />
                    <Route
                        path="/services"
                        element={
                            isAuthenticated ? <Services /> : <Navigate to="/" />
                        }
                    />
                </Routes>
            </BrowserRouter>
            <Toaster position="bottom-center" reverseOrder={false} />
        </>
    );
}

export default App;
