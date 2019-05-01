from locust import (HttpLocust, TaskSet, TaskSequence, task, seq_task)
from locust import events
from bravado.client import SwaggerClient
from bravado.requests_client import RequestsClient
from bravado.exception import HTTPError


class AnonBehavior(TaskSet):

    @task(1)
    def index(self):
        self.client.get("/")


class AnonUser(HttpLocust):
    host = "http://milmovelocal:8080"
    # weight = 5  # 5x more likely than other users
    weight = 1
    task_set = AnonBehavior


class MilMoveUserBehavior(TaskSequence):

    swagger_internal = None
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

    def update_service_member(self, service_member):
        self.user["service_member"] = service_member

    def update_duty_stations(self, duty_stations):
        self.user["duty_stations"] = duty_stations

    def swagger_internal_wrapper(self, callable_operation, *args, **kwargs):
        """
        Swagger client uses requests send() method instead of request(). This means we need to send off
        events to Locust on our own.
        """
        method = callable_operation.operation.http_method.upper()
        path_name = callable_operation.operation.path_name
        response_future = callable_operation(*args, **kwargs)
        try:
            response = response_future.response()
        except HTTPError as e:
            events.request_failure.fire(
                request_type=method,
                name=path_name,
                response_time=0,  # Not clear how to get this
                exception=e,
            )
            return e.swagger_result
        else:
            metadata = response.metadata

            events.request_success.fire(
                request_type=method,
                name=path_name,
                response_time=metadata.elapsed_time,
                response_length=len(metadata.incoming_response.raw_bytes),
            )
            return response.result

    @task(2)
    def load_swagger_file_internal(self):
        self.client.get("/internal/swagger.yaml")

    @task(2)
    def load_swagger_file_public(self):
        self.client.get("/api/v1/swagger.yaml")

    @seq_task(1)
    def login(self):
        resp = self.client.post('/devlocal-auth/create', data={"userType": "milmove"})
        try:
            self.login_gov_user = resp.json()
            self.session_token = self.client.cookies.get('mil_session_token')
            self.requests_client = RequestsClient()
            # Set the session to be the same session as locust uses
            self.requests_client.session = self.client
            # Set the csrf token in the global headers for all requests
            # Don't validate requests or responses because we're using OpenAPI Spec 2.0 which doesn't respect
            # nullable sub-definitions
            self.swagger_internal = SwaggerClient.from_url(
                "http://milmovelocal:8080/internal/swagger.yaml",
                request_headers={'x-csrf-token': self.csrf},
                http_client=self.requests_client,
                config={
                    'validate_requests': False,
                    'validate_responses': False,
                })
        except Exception:
            print(resp.content)

    @seq_task(2)
    def retrieve_user(self):
        resp = self.client.get("/internal/users/logged_in")
        self.user = resp.json()
        # check response for 200

    @seq_task(3)
    def create_service_member(self):
        model = self.swagger_internal.get_model("CreateServiceMemberPayload")
        payload = model(user_id=self.user["id"])
        service_member = self.swagger_internal_wrapper(
            self.swagger_internal.service_members.createServiceMember,
            createServiceMemberPayload=payload)
        self.update_service_member(service_member)

    @seq_task(4)
    def create_your_profile(self):
        model = self.swagger_internal.get_model("PatchServiceMemberPayload")
        payload = model(
            affiliation="NAVY",  # Rotate
            edipi="3333333333",  # Random
            rank="E_5",  # Rotate
            social_security_number="333-33-3333",  # Random
        )
        service_member = self.swagger_internal_wrapper(
            self.swagger_internal.service_members.patchServiceMember,
            serviceMemberId=self.user["service_member"].id,
            patchServiceMemberPayload=payload)
        self.update_service_member(service_member)

    @seq_task(5)
    def create_your_name(self):
        model = self.swagger_internal.get_model("PatchServiceMemberPayload")
        payload = model(
            first_name="Alice",  # Random
            last_name="Bob",  # Random
            middle_name="Carol",
            suffix="",
        )
        service_member = self.swagger_internal_wrapper(
            self.swagger_internal.service_members.patchServiceMember,
            serviceMemberId=self.user["service_member"].id,
            patchServiceMemberPayload=payload)
        self.update_service_member(service_member)

    @seq_task(6)
    def create_your_contact_info(self):
        model = self.swagger_internal.get_model("PatchServiceMemberPayload")
        payload = model(
            email_is_preferred=True,
            personal_email="20190321164732@example.com",
            phone_is_preferred=True,
            secondary_telephone="333-333-3333",
            telephone="333-333-3333",
        )
        service_member = self.swagger_internal_wrapper(
            self.swagger_internal.service_members.patchServiceMember,
            serviceMemberId=self.user["service_member"].id,
            patchServiceMemberPayload=payload)
        self.update_service_member(service_member)

    @seq_task(7)
    def search_for_duty_station(self):
        station_list = ["b", "buck", "buckley"]
        for station in station_list:
            duty_stations = self.swagger_internal_wrapper(
                self.swagger_internal.duty_stations.searchDutyStations,
                search=station)
            self.update_duty_stations(duty_stations)

    @seq_task(8)
    def current_duty_station(self):
        model = self.swagger_internal.get_model("PatchServiceMemberPayload")
        payload = model(
            current_station_id=self.user["duty_stations"][0].id
        )
        service_member = self.swagger_internal_wrapper(
            self.swagger_internal.service_members.patchServiceMember,
            serviceMemberId=self.user["service_member"].id,
            patchServiceMemberPayload=payload)
        self.update_service_member(service_member)

    @seq_task(9)
    def logout(self):
        self.client.post("/auth/logout")
        self.login_gov_user = None
        self.session_token = None
        self.user = {}


class MilMoveUser(HttpLocust):
    host = "http://milmovelocal:8080"
    weight = 1
    task_set = MilMoveUserBehavior
    min_wait = 1000
    max_wait = 5000


class OfficeUserBehavior(TaskSequence):

    swagger_internal = None
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

    def swagger_internal_wrapper(self, callable_operation, *args, **kwargs):
        """
        Swagger client uses requests send() method instead of request(). This means we need to send off
        events to Locust on our own.
        """
        method = callable_operation.operation.http_method.upper()
        path_name = callable_operation.operation.path_name
        response_future = callable_operation(*args, **kwargs)
        try:
            response = response_future.response()
        except HTTPError as e:
            events.request_failure.fire(
                request_type=method,
                name=path_name,
                response_time=0,  # Not clear how to get this
                exception=e,
            )
            return e.swagger_result
        else:
            metadata = response.metadata

            events.request_success.fire(
                request_type=method,
                name=path_name,
                response_time=metadata.elapsed_time,
                response_length=len(metadata.incoming_response.raw_bytes),
            )
            return response.result

    @task(2)
    def load_swagger_file_internal(self):
        self.client.get("/internal/swagger.yaml")

    @task(2)
    def load_swagger_file_public(self):
        self.client.get("/api/v1/swagger.yaml")

    @seq_task(1)
    def login(self):
        resp = self.client.post('/devlocal-auth/create', data={"userType": "office"})
        try:
            self.login_gov_user = resp.json()
            self.session_token = self.client.cookies.get('office_session_token')
            self.requests_client = RequestsClient()
            # Set the session to be the same session as locust uses
            self.requests_client.session = self.client
            # Set the csrf token in the global headers for all requests
            # Don't validate requests or responses because we're using OpenAPI Spec 2.0 which doesn't respect
            # nullable sub-definitions
            self.swagger_internal = SwaggerClient.from_url(
                "http://officelocal:8080/internal/swagger.yaml",
                request_headers={'x-csrf-token': self.csrf},
                http_client=self.requests_client,
                config={
                    'validate_requests': False,
                    'validate_responses': False,
                })
        except Exception:
            print(resp.content)

    @seq_task(2)
    def retrieve_user(self):
        resp = self.client.get("/internal/users/logged_in")
        self.user = resp.json()
        # check response for 200

    @seq_task(3)
    def view_new_moves_queue(self):
        self.swagger_internal_wrapper(
            self.swagger_internal.queues.showQueue,
            queueType="new")

    @seq_task(4)
    def logout(self):
        self.client.post("/auth/logout")
        self.login_gov_user = None
        self.session_token = None
        self.user = {}


class OfficeUser(HttpLocust):
    host = "http://officelocal:8080"
    weight = 1
    task_set = OfficeUserBehavior
    min_wait = 1000
    max_wait = 5000
