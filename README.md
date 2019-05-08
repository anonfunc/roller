# roller
A dice rolling service.

# Basic
## Build / Run locally:

    mage build
    ./roller

# Serverless

## Deploy via SAM / AWS Lambda

    # Configure / login AWS and SAM CLIs
    mage -v deploySAM

## SAM Local

    mage -v runSAM


# Containers

## Run via Docker:

    mage -v docker dockerRun
