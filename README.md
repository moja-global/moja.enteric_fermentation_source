# Demo FLINT application

building a website in a go program
accepts login and authenticates against correct auth system
shows datalayers available

it uses gobuffalo/packr2 to bundle the files into an exec

To use this, you will need to create a credentials file (in our case, it is ef_credentials.json) from Google Cloud -> API & Services -> Credentials. In "Credentials", create a credentials files (new API key and Oauth 2.0 Client IDs) and put the domain of your localhost or website into the whitelist with a wildcard "*" at the end. Then, you can download the Oauth 2.0 Client IDs as a json file, which is also the ef_credentials.json you needed.
