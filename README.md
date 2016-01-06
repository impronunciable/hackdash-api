# HackDash API

API for [HackDash.org](https://hackdash.org).
Join the discussion on [Slack](http://hackdash-slack.herokuapp.com/)

## Requirements

- Docker
- Docker compose

## Installation

    $ docker-compose build
    $ docker-compose up

## How to get a token and call the api

Go to the [Auth0 Playground](http://auth0.github.io/playground/) and replace the domain and the client_id for the same used in this api. Make sure you are dumping the token to the console in the callback, it should look like this:

```
var domain = '{tenant_name}.auth0.com';
​
var clientID = '...client_id...'; 
​
var lock = new Auth0Lock(clientID, domain);
lock.show({
  authParams: { scope: 'openid first_name family_name email picture' },
}, function (err, profile, token) {
  console.log(err, token);
});
```

To call the API, copy the token and add the authorization header to the request:

`Authorization: Bearer ...token...`