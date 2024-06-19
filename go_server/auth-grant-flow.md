Auth Grant Flow (Secure Design)

Server creates OauthLogin URl for the Client
and is sent to the Client

The client then sends the ClientId and Redirects to that URl for the Googl OAuth Server for login
It then receives a Code from the Google OAuth Server

The CLient then shares this code with the backend go server, and the server communicated with the Google OAuth Server to verify the code along with the Google Client Secret and retrieve
Access Token - User fitness data, short time period
Refresh Token - Refresh and creates new Access Token, DONOT SHARE WITH CLIENT
ID Token - Identification Token, metdata about user

None is shared with the client back

Instead Create self signed JWT token to be shared with the client so, it can verify with the server using the JWT token
