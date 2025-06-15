import { ALLOWED_SERVICES } from "./constants";

const Services = () => {
    return (
        <div className="services-container">
            <div className="title">MetaKGP Services</div>
            <div className="subtitle">
                Click on any of the links below to visit
            </div>
            {Object.entries(ALLOWED_SERVICES).map(([url, serviceName]) => (
                <a className="service" key={url} href={url}>
                    {serviceName}
                </a>
            ))}
        </div>
    );
};

export default Services;
