import { servicesList } from "./constants";

const Services = () => {
    return (
        <div className="services-container">
            <div className="title">MetaKGP Services</div>
            <div className="subtitle">
                Click on any of the links below to visit
            </div>
            {servicesList.map((service) => (
                <a className="service" key={service} href={service}>
                    {service}
                </a>
            ))}
        </div>
    );
};

export default Services;
