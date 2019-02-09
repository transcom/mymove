#! /usr/bin/env python3

import urllib

import requests

MILMOVE = 'http://milmovelocal:8080'


def get_csrf_token(client):
    """
    Pull the CSRF token from the website by hitting the root URL.

    The token is set as a cookie with the name `masked_gorilla_csrf`
    """
    client.get(urllib.parse.urljoin(MILMOVE, '/'))
    return client.cookies.get('masked_gorilla_csrf')


def create_user(client, csrf):
    """
    Create a new user for local testing using the CSRF token in the header
    """
    resp = client.post(urllib.parse.urljoin(MILMOVE, 'devlocal-auth/create'), headers={'x-csrf-token': csrf})
    print(resp.content)


if __name__ == "__main__":
    client = requests.session()
    csrf = get_csrf_token(client)
    create_user(client, csrf)
