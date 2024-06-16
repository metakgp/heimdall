import { servicesList } from "./constants";

export const validateEmail = (email: string): boolean => {
    return (
        email.endsWith("@kgpian.iitkgp.ac.in") ||
        email.endsWith("@iitkgp.ac.in")
    );
};

export const validRedirect = (link: string | null): string => {
    if (!link) {
        return "/services";
    }

    for (var service of servicesList) {
        if (link.startsWith(service)) {
            return link;
        }
    }

    return "/services";
};
