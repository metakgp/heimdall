export const BACKEND_URL = import.meta.env.PROD
    ? "https://heimdall-api.metakgp.org"
    : "http://localhost:3333";

export type AllowedServices = {
    [key: string]: string;
};

export const ALLOWED_SERVICES: AllowedServices = {
    "https://chill.metakgp.org": "Chillzone",
    "https://naarad.metakgp.org": "Naarad",
};
