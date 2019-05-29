#! /usr/bin/env python3

import urllib

import requests


class Demo(object):
    def __init__(
        self, url, user_type="milmove", session_token_name="mil_session_token"
    ):
        self.session = requests.Session()
        self.user_type = user_type
        self.session_token_name = session_token_name
        self.url = url
        self.csrf = None
        self.user = None
        self.token = None
        self.logged_in = None

    def _get_csrf_token(self):
        """
        Pull the CSRF token from the website by hitting the root URL.

        The token is set as a cookie with the name `masked_gorilla_csrf`
        """
        if self.csrf:
            return self.csrf
        self.session.get(urllib.parse.urljoin(self.url, "/"))
        self.csrf = self.session.cookies.get("masked_gorilla_csrf")
        return self.csrf

    def create_user(self):
        """
        Create a new user for local testing using the CSRF token in the header
        """
        resp = self.session.post(
            urllib.parse.urljoin(self.url, "devlocal-auth/create"),
            headers={"x-csrf-token": self._get_csrf_token()},
            data={"userType": self.user_type},
        )
        try:
            self.user = resp.json()
        except Exception:
            print(resp.content)
        self.token = self.session.cookies.get(self.session_token_name)

    def get_logged_in_user(self):
        """
        Get the logged in user information
        """
        resp = self.session.get(
            urllib.parse.urljoin(self.url, "internal/users/logged_in")
        )
        try:
            self.logged_in = resp.json()
        except Exception:
            print(resp.content)


if __name__ == "__main__":

    milmove_demo = Demo(
        "http://milmovelocal:8080",
        user_type="milmove",
        session_token_name="mil_session_token",
    )
    office_demo = Demo(
        "http://officelocal:8080",
        user_type="office",
        session_token_name="office_session_token",
    )
    tsp_demo = Demo(
        "http://tsplocal:8080", user_type="tsp", session_token_name="tsp_session_token"
    )
    for demo in [milmove_demo, office_demo, tsp_demo]:
        print(demo.url)
        demo.create_user()
        print("User", demo.user)
        print("Token", demo.token)
        demo.get_logged_in_user()
        print("LoggedIn", demo.logged_in)
        print()
