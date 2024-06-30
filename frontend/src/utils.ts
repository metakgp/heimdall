import { ALLOWED_SERVICES } from "./constants";

export const validateEmail = (email: string): boolean => {
    return (
        email.endsWith("@kgpian.iitkgp.ac.in") ||
        email.endsWith("@iitkgp.ac.in")
    );
};

export const validRedirect = (link: string | null): { serviceName: string, serviceLink: string } => {
    if (!link) {
        return { serviceName: "Services Page", serviceLink: "/services" };
    }

    for (const serviceUrl in ALLOWED_SERVICES) {
        if (link.startsWith(serviceUrl)) {
            return { serviceName: ALLOWED_SERVICES[serviceUrl], serviceLink: link };
        }
    }

    return { serviceName: "Services Page", serviceLink: "/services" };
};
