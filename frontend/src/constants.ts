export const BACKEND_URL = import.meta.env.PROD
    ? "https://api.heimdall.metakgp.org"
    : "http://localhost:3333";

export const servicesList = [
    "https://naarad.metakgp.org",
    "https://chill.metakgp.org",
];
