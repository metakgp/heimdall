export const BACKEND_URL = import.meta.env.PROD
    ? "https://heimdall-api.metakgp.org"
    : "http://localhost:3333";

export const servicesList = [
    "https://naarad.metakgp.org/web",
    "https://chill.metakgp.org",
];
