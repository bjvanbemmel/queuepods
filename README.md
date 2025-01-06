# Queuepod with monitoring solution

## Getting started
First, create the necessary environment file to configure and run the software suite:
```bash
cp .env.example .env
```
Change the values to whatever you deem useful for your usecase. Do know that defauls values have been provided, so you can always leave them empty (just make sure the `.env` file exists).

In order to run this software solution, you must first verify whether or not your docker-engine is running in privileged mode.
If it is, then you can simply open a shell within the root directory of this repository and run the following command:
```bash
docker compose up -d --build
```
If, however, your docker-engine is not running under root, then you must execute this command as root:
```bash
sudo docker compose up -d --build
```
This is because the python (queuepod) package requires access to the `/dev/ttyACM0` device, which requires elevated privileges.
