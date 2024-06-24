# Frontend

## Heimdall

Heimdall checks the client's IP to know whether the request has originated from inside the IIT Kharagpur network and verifies their institute email ID. This helps to ascertain if the client is a current member of the institute and should have access to certain information.

### Running locally

First install [nodejs](https://nodejs.org/en/download/package-manager). Then install `pnpm` by running `npm install -g pnpm`. 

Then follow the given steps to start the development server:

1. Clone the repository
   ```sh
   git clone https://github.com/metakgp/heimdall.git
   ```
2. Install dependencies
   ```sh
   cd heimdall/frontend
   pnpm install
   ```
3. Start the server
   ```sh
   npm run dev
   ```

This setup will launch the frontend. To start backend server also, please follow the instructions [here](https://github.com/metakgp/heimdall/blob/master/README.md#getting-started)