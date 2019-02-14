from locust import HttpLocust, TaskSet, task


class UserBehavior(TaskSet):

    user = None

    def _get_csrf_token(self):
        """
        Pull the CSRF token from the website by hitting the root URL.

        The token is set as a cookie with the name `masked_gorilla_csrf`
        """
        self.client.get('/')
        return self.client.cookies.get('masked_gorilla_csrf')

    def on_start(self):
        """ on_start is called when a Locust start before any task is scheduled """
        self.login()

    def on_stop(self):
        """ on_stop is called when the TaskSet is stopping """
        self.logout()

    def login(self):
        csrf = self._get_csrf_token()
        resp = self.client.post('/devlocal-auth/create',
                                headers={'x-csrf-token': csrf})
        self.user = resp.json()
        resp = self.client.post('/devlocal-auth/login',
                                headers={'x-csrf-token': csrf},
                                data={'id': self.user["id"]})
        print(resp.status_code, resp.content)

    def logout(self):
        self.client.post("/auth/logout")

    @task(1)
    def index(self):
        self.client.get("/")


class WebsiteUser(HttpLocust):
    task_set = UserBehavior
    min_wait = 5000
    max_wait = 9000
