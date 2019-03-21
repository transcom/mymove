from locust import (HttpLocust, TaskSet, TaskSequence, task, seq_task)
from bravado.client import SwaggerClient
from bravado.requests_client import RequestsClient


class AnonBehavior(TaskSet):

    @task(1)
    def index(self):
        self.client.get("/")


class AnonUser(HttpLocust):
    host = "http://milmovelocal:8080"
    weight = 5  # 5x more likely than other users
    task_set = AnonBehavior


class UserBehavior(TaskSequence):

    swagger = None
    login_gov_user = None
    csrf = None
    session_token = None
    user = {}

    def _get_csrf_token(self):
        """
        Pull the CSRF token from the website by hitting the root URL.

        The token is set as a cookie with the name `masked_gorilla_csrf`
        """
        self.client.get('/')
        self.csrf = self.client.cookies.get('masked_gorilla_csrf')
        self.client.headers.update({'x-csrf-token': self.csrf})

    def on_start(self):
        """ on_start is called when a Locust start before any task is scheduled """
        self._get_csrf_token()

    def on_stop(self):
        """ on_stop is called when the TaskSet is stopping """
        pass

    def update_service_member(self, service_member_future):
        service_member_response = service_member_future.response()
        self.user["service_member"] = service_member_response.result

    def update_duty_stations(self, duty_stations_future):
        duty_stations_response = duty_stations_future.response()
        print(duty_stations_response.result)
        self.user["duty_stations"] = duty_stations_response.result

    @seq_task(1)
    def login(self):
        resp = self.client.post('/devlocal-auth/create')
        try:
            self.login_gov_user = resp.json()
            self.session_token = self.client.cookies.get('mil_session_token')
            self.requests_client = RequestsClient()
            # Set the session to be the same session as locust uses
            self.requests_client.session = self.client
            # Set the csrf token in the global headers for all requests
            # Don't validate requests or responses because we're using OpenAPI Spec 2.0 which doesn't respect
            # nullable sub-definitions
            self.swagger = SwaggerClient.from_url(
                "http://milmovelocal:8080/internal/swagger.yaml",
                request_headers={'x-csrf-token': self.csrf},
                http_client=self.requests_client,
                config={
                    'validate_requests': False,
                    'validate_responses': False,
                })
        except Exception:
            print('Headers:', self.client.headers)
            print(resp.content)

    @seq_task(2)
    def retrieve_user(self):
        resp = self.client.get("/internal/users/logged_in")
        self.user = resp.json()
        # check response for 200

    @seq_task(3)
    def create_service_member(self):
        model = self.swagger.get_model("CreateServiceMemberPayload")
        payload = model(user_id=self.user["id"])
        service_member_future = self.swagger.service_members.createServiceMember(
            createServiceMemberPayload=payload)
        self.update_service_member(service_member_future)

    @seq_task(4)
    def create_your_profile(self):
        model = self.swagger.get_model("PatchServiceMemberPayload")
        payload = model(
            affiliation="NAVY",  # Rotate
            edipi="3333333333",  # Random
            rank="E_5",  # Rotate
            social_security_number="333-33-3333",  # Random
        )
        service_member_future = self.swagger.service_members.patchServiceMember(
            serviceMemberId=self.user["service_member"].id,
            patchServiceMemberPayload=payload)
        self.update_service_member(service_member_future)

    @seq_task(5)
    def create_your_name(self):
        model = self.swagger.get_model("PatchServiceMemberPayload")
        payload = model(
            first_name="Alice",  # Random
            last_name="Bob",  # Random
            middle_name="Carol",
            suffix="",
        )
        service_member_future = self.swagger.service_members.patchServiceMember(
            serviceMemberId=self.user["service_member"].id,
            patchServiceMemberPayload=payload)
        self.update_service_member(service_member_future)

    @seq_task(6)
    def create_your_contact_info(self):
        model = self.swagger.get_model("PatchServiceMemberPayload")
        payload = model(
            email_is_preferred=True,
            personal_email="20190321164732@example.com",
            phone_is_preferred=True,
            secondary_telephone="333-333-3333",
            telephone="333-333-3333",
        )
        service_member_future = self.swagger.service_members.patchServiceMember(
            serviceMemberId=self.user["service_member"].id,
            patchServiceMemberPayload=payload)
        self.update_service_member(service_member_future)

    @seq_task(7)
    def search_for_duty_station_1(self):
        station_list = ["b", "buck", "buckley"]
        for station in station_list:
            duty_stations_future = self.swagger.duty_stations.searchDutyStations(
                search=station)
            self.update_duty_stations(duty_stations_future)

    @seq_task(8)
    def current_duty_station(self):
        model = self.swagger.get_model("PatchServiceMemberPayload")
        payload = model(
            current_station_id=self.user["duty_stations"][0].id
        )
        service_member_future = self.swagger.service_members.patchServiceMember(
            serviceMemberId=self.user["service_member"].id,
            patchServiceMemberPayload=payload)
        self.update_service_member(service_member_future)

    @seq_task(9)
    def logout(self):
        self.client.post("/auth/logout")
        self.login_gov_user = None
        self.session_token = None
        self.user = {}


class MilMoveUser(HttpLocust):
    host = "http://milmovelocal:8080"
    weight = 1
    task_set = UserBehavior
    min_wait = 1000
    max_wait = 5000
