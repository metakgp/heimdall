export const validateEmail = (email: string): boolean => {
    return (
        email.endsWith("@kgpian.iitkgp.ac.in") ||
        email.endsWith("@iitkgp.ac.in")
    );
};
