#! /usr/bin/env python3

import urllib

import requests


class Demo(object):

    def __init__(self, url):
        self.session = requests.Session()
        self.url = url
        self.csrf = None
        self.user = None
        self.token = None

    def _get_csrf_token(self):
        """
        Pull the CSRF token from the website by hitting the root URL.

        The token is set as a cookie with the name `masked_gorilla_csrf`
        """
        if self.csrf:
            return self.csrf
        self.session.get(urllib.parse.urljoin(self.url, '/'))
        self.csrf = self.session.cookies.get('masked_gorilla_csrf')
        return self.csrf

    def create_user(self):
        """
        Create a new user for local testing using the CSRF token in the header
        """
        resp = self.session.post(urllib.parse.urljoin(self.url, 'devlocal-auth/create'),
                                 headers={'x-csrf-token': self._get_csrf_token()},
                                 data={"userType": "milmove"})
        try:
            self.user = resp.json()
        except Exception:
            print(resp.content)
        self.token = self.session.cookies.get('mil_session_token')


if __name__ == "__main__":

    demo = Demo('http://milmovelocal:8080')
    demo.create_user()
    print('User', demo.user)
    print('Token', demo.token)
