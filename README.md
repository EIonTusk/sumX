# sumX  
A small webservice that summarizes tweets of any X account.

## Prerequisites
* npm
* docker
## Quick Start

```bash
# Clone project and go into it
$ git clone https://github.com/EIonTusk/sumX.git
$ cd sumX/

# Install required node packages
$ npm install

# Create the .env files
$ cp .env.template .env
$ cp ui/.env.template ui/.env

# Add your X and Huggingface Bearer token and edit the postgres settings in `sumX/.env`
# Edit the variable in `sumX/ui/.env`

# Build/Start backend and db
$ sudo docker compose up --build

# Run the dev server
$ npm run dev

```
### Advanced usage
Instead of running the dev server using `npm run dev` you may consider building it using `npm run build` and using a web server like nginx or apache2.
