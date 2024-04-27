# Preface

I tried and be extra verbose in the code to give an idea of my thought process and some of the assumptions I made, but I'll try and add the key points here as well, as well as any extra bits I missed in the code


## Project Structure
 - **bin**
    Output of the `make build` command. Builds for *Windows*, *Mac*, and *Linux*
 - **cmd**
    Houses the `main.go` file and the **api** package. **api** holds the server itself, as well as the middleware (which in this case is just a simple route logger)
 - **services**
    The only service here is **risk**, but it holds the route handler functions and risk structure. In production I probably would have added another layer for storage, but given I'm just using a map and pretending it's a real db this seems fine
 - **utils**
    For now just holds `JSONEncode` to cut out some boilerplate, but could just as easily hold a common `Encode` in the future to allow for different protocols to be used.


## Notes
This has honestly just turned into an excuse for me to test out the new go 1.22 http package

I would have liked to add an easy way to update the server params (and any future ones we get). Depeding on how this is used I usually either just use flags (for cmd line type stuff), or env vars (with godotenv) for docker and cloud based stuff. 

The testing for now is very sparse. I just added some tests to the route handlers, but obviously in a real production setup I'd want to expand those and add a lot more tests for the other elements that I don't test at all currently. Not to mention it's just unit testing for now, but I feel like integration/end-to-end is out of the scope here unless I dockerized the app and set something really fancy up. 

The makefile is also pretty barebones, I imagine we'd want a lot more utility functions for fmt/tidy/linting as well as docker builds and deployment (depending on the cloud arch)

I make the assumption (mostly given this is just HTTP for now), that this is internal and don't do all of the validation on objects I otherwise would. Using the validator package really cuts down on that for the JSON validation, but there are a lot of spots in the could that could use extra validation so we know we aren't grabbing data that is malformed.

Obviously, this is still pretty bare bones. If I had infinite time and wanted to really snaz this up I'd probably
 - dockerize this for easy deployment
 - add an auth subrouter
 - HTTPS (as there really isn't many cases you'd ever want a producdtion HTTP server IMO)
 - Add filtering, pagination, sorting, etc to the `GET \risks\` endpoint
 - some database layer for persistance
