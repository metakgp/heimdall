<div id="top"></div>

<!-- PROJECT SHIELDS -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links-->
<div align="center">

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![Wiki][wiki-shield]][wiki-url]

</div>

<!-- PROJECT LOGO -->
<br />
<!-- UPDATE -->
<div align="center">
  <a href="https://github.com/metakgp/heimdall">
    <img width="140" alt="image" src="https://user-images.githubusercontent.com/86282911/206632284-cb260f57-c612-4ab5-b92b-2172c341ab23.png">
  </a>

  <h3 align="center">Heimdall</h3>

  <p align="center">
  <!-- UPDATE -->
    <i>The network checker for IIT KGP</i>
    <br />
    <a href="https://github.com/metakgp/heimdall/issues">Report Bug</a>
    ·
    <a href="https://github.com/metakgp/heimdall/issues">Request Feature</a>
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
<summary>Table of Contents</summary>

- [About The Project](#about-the-project)
  - [Supports](#supports)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
- [Maintainer(s)](#maintainers)
- [Contact](#contact)
- [Additional documentation](#additional-documentation)

</details>

<!-- ABOUT THE PROJECT -->

## About The Project

<!-- UPDATE -->
<div align="center">
  <a href="https://github.com/metakgp/heimdall">
    <img width="80%" alt="image" src="./assets/screenshot.png">
  </a>
</div>

_Heimdall checks the client's IP to know whether the request has originated from inside the IIT Kharagpur network. This helps to ascertain if the client is a current member of the institute and should have access to certain information._

<p align="right">(<a href="#top">back to top</a>)</p>

<div id="supports"></div>

### Supports:

<!-- UPDATE -->

1. Shells
   - `bash`
   - `zsh`
2. OS(s)
   - any `*nix`[`GNU+Linux` and `Unix`]

<p align="right">(<a href="#top">back to top</a>)</p>

## Getting Started

To set up a local instance of the application, follow the steps below.

### Prerequisites

The following dependencies are required to be installed for the project to function properly:

<!-- UPDATE -->

- [go](https://go.dev/)

<p align="right">(<a href="#top">back to top</a>)</p>

### Installation

_Now that the environment has been set up and configured to properly compile and run the project, the next step is to install and configure the project locally on your system._

<!-- UPDATE -->

1. Clone the repository
   ```sh
   git clone https://github.com/metakgp/heimdall.git
   ```
2. Install dependencies
   ```sh
   cd ./heimdall
   go mod download
   ```
3. Execute the script
   ```sh
   go run main.go
   ```

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- BACKGROUND INFORMATION -->

## How does this work?

<!-- UPDATE -->

IIT Kharagpur has its internal campus network which is the primary source of Internet for its students, staff and faculty.

For connection to the outside network (normal internet services), it uses a [NAT Gateway](https://docs.aws.amazon.com/vpc/latest/userguide/vpc-nat-gateway.html) which handles all requests going outside. While doing so, the client IP address in those requests is changed from the internal IP to any one of the pool of public IP addresses assigned to IIT Kharagpur.

So, to check whether a request has originated from inside the IIT Kharagpur network, we just check whether the client's IP address belongs to one of those public IPs.

While just doing this would have sufficed, we do not know whether these Public IPs are permanent or are subject to change over time. We therefore do a Whois lookup of the IP address and check its response to decide whether this IP address belongs to IIT Kharagpur. A screenshot of such a Whois lookup is shown below.

<div align="center">
  <a href="https://github.com/metakgp/heimdall">
    <img width="80%" alt="image" src="./assets/whois.png">
  </a>

  _For complete Whois information check [here](https://whois.domaintools.com/203.110.242.30)._
</div>

<p align="right">(<a href="#top">back to top</a>)</p>

All this time you might be wondering why you need a different server to just check this. Can't we do this in any project where such a feature is required?

Well yes. Provided it has a backend server. This cannot be done in the front-end because the Web Browser does not provide the IP information to the Javascript engine. So for projects that do not need a backend, like [Chillzone](https://github.com/metakgp/chillzone) or [ERP Auto Login](https://github.com/metakgp/iitkgp-erp-auto-login), this simple API call can do the required work.

## Maintainer(s)

<!-- UPDATE -->

- [Chirag Ghosh](https://github.com/chirag-ghosh)

<p align="right">(<a href="#top">back to top</a>)</p>

## Contact

<p>
📫 Metakgp -
<a href="https://bit.ly/metakgp-slack">
  <img align="center" alt="Metakgp's slack invite" width="22px" src="https://raw.githubusercontent.com/edent/SuperTinyIcons/master/images/svg/slack.svg" />
</a>
<a href="mailto:metakgp@gmail.com">
  <img align="center" alt="Metakgp's email " width="22px" src="https://raw.githubusercontent.com/edent/SuperTinyIcons/master/images/svg/gmail.svg" />
</a>
<a href="https://www.facebook.com/metakgp">
  <img align="center" alt="metakgp's Facebook" width="22px" src="https://raw.githubusercontent.com/edent/SuperTinyIcons/master/images/svg/facebook.svg" />
</a>
<a href="https://www.linkedin.com/company/metakgp-org/">
  <img align="center" alt="metakgp's LinkedIn" width="22px" src="https://raw.githubusercontent.com/edent/SuperTinyIcons/master/images/svg/linkedin.svg" />
</a>
<a href="https://twitter.com/metakgp">
  <img align="center" alt="metakgp's Twitter " width="22px" src="https://raw.githubusercontent.com/edent/SuperTinyIcons/master/images/svg/twitter.svg" />
</a>
<a href="https://www.instagram.com/metakgp_/">
  <img align="center" alt="metakgp's Instagram" width="22px" src="https://raw.githubusercontent.com/edent/SuperTinyIcons/master/images/svg/instagram.svg" />
</a>
</p>

<p align="right">(<a href="#top">back to top</a>)</p>

## Additional documentation

- [License](/LICENSE)
- [Code of Conduct](/.github/CODE_OF_CONDUCT.md)
- [Security Policy](/.github/SECURITY.md)
- [Contribution Guidelines](/.github/CONTRIBUTING.md)

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->

[contributors-shield]: https://img.shields.io/github/contributors/metakgp/heimdall.svg?style=for-the-badge
[contributors-url]: https://github.com/metakgp/heimdall/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/metakgp/heimdall.svg?style=for-the-badge
[forks-url]: https://github.com/metakgp/heimdall/network/members
[stars-shield]: https://img.shields.io/github/stars/metakgp/heimdall.svg?style=for-the-badge
[stars-url]: https://github.com/metakgp/heimdall/stargazers
[issues-shield]: https://img.shields.io/github/issues/metakgp/heimdall.svg?style=for-the-badge
[issues-url]: https://github.com/metakgp/heimdall/issues
[license-shield]: https://img.shields.io/github/license/metakgp/heimdall.svg?style=for-the-badge
[license-url]: https://github.com/metakgp/heimdall/blob/master/LICENSE
[wiki-shield]: https://custom-icon-badges.demolab.com/badge/metakgp_wiki-grey?logo=metakgp_logo&style=for-the-badge
[wiki-url]: https://wiki.metakgp.org
