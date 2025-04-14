This script automates the entire OAuth2 flow and made it executable. Let me explain how to use it:

How to Use the Script
Basic usage:

```bash
./help_script_test_userinfo.sh
```

This will:

Get an authorization code from the /authorize endpoint
Exchange the code for an access token
Use the token to query the /userinfo endpoint
Optional Arguments
You can customize the script's behavior with these options:

-h, --host HOST: Specify the OAuth2 server host (default: localhost)
-p, --port PORT: Specify the OAuth2 server port (default: 8080)
-c, --client CLIENT_ID: Use a custom client ID (default: test-client)
-s, --secret SECRET: Use a custom client secret (default: test-secret)
-r, --redirect URI: Use a custom redirect URI (default:  http://localhost/callback)
-u, --update: Enable user info updates after retrieval
--help: Show the help message
Updating User Info
To update the user info and test the changes:

This will perform the regular OAuth2 flow and then:

Send a request to the /config endpoint with new user information
Fetch the updated user info to verify the changes worked

The script requires jq for JSON processing - if you don't have it installed, you might need to install it with:

```bash
apt-get install jq
```

The script includes robust error handling, color-coded output for readability, and step-by-step progress indicators to help you understand what's happening at each stage of the process.