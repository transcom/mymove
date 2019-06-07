# -*- coding: utf-8 -*-
import datetime
import pprint
import os
import random
from urllib.parse import urljoin
import uuid

from locust import seq_task
from bravado.client import SwaggerClient
from bravado.requests_client import RequestsClient

from .base import BaseTaskSequence
from .base import InternalAPIMixin
from .base import get_swagger_config
from .base import swagger_request


class MilMoveUserBehavior(BaseTaskSequence, InternalAPIMixin):

    login_gov_user = None
    session_token = None

    fixtures_path = os.path.join(
        os.path.dirname(os.path.realpath(__file__)), "fixtures"
    )

    # User is the LoggedInUserPayload
    user = None
    duty_stations = []
    new_duty_stations = []

    def update_user(self):
        self.user = swagger_request(self.swagger_internal.users.showLoggedInUser)

    def update_service_member(self, service_member):
        self.user.service_member = service_member

    def get_dutystations(self, short_name):
        station_list = [short_name[0], short_name[0:3], short_name]
        duty_stations = None
        for station in station_list:
            duty_stations = swagger_request(
                self.swagger_internal.duty_stations.searchDutyStations, search=station
            )
        return duty_stations

    @seq_task(1)
    def login(self):
        resp = self.client.post("/devlocal-auth/create", data={"userType": "milmove"})
        try:
            self.login_gov_user = resp.json()
            self.session_token = self.client.cookies.get("mil_session_token")
            self.requests_client = RequestsClient()
            # Set the session to be the same session as locust uses
            self.requests_client.session = self.client
            # Set the csrf token in the global headers for all requests
            # Don't validate requests or responses because we're using OpenAPI Spec 2.0
            # which doesn't respect nullable sub-definitions
            self.swagger_internal = SwaggerClient.from_url(
                urljoin(self.parent.host, "internal/swagger.yaml"),
                request_headers={"x-csrf-token": self.csrf},
                http_client=self.requests_client,
                config=get_swagger_config(),
            )
        except Exception:
            print(resp.content)

    @seq_task(2)
    def retrieve_user(self):
        self.update_user()

    @seq_task(3)
    def create_service_member(self):
        model = self.swagger_internal.get_model("CreateServiceMemberPayload")
        payload = model(user_id=self.user["id"])
        service_member = swagger_request(
            self.swagger_internal.service_members.createServiceMember,
            createServiceMemberPayload=payload,
        )
        self.update_service_member(service_member)

    @seq_task(4)
    def create_your_profile(self):
        model = self.swagger_internal.get_model("PatchServiceMemberPayload")
        payload = model(
            affiliation="NAVY",  # Rotate
            edipi="3333333333",  # Random
            social_security_number="333-33-3333",  # Random
            rank="E_5",  # Rotate
        )
        service_member = swagger_request(
            self.swagger_internal.service_members.patchServiceMember,
            serviceMemberId=self.user["service_member"].id,
            patchServiceMemberPayload=payload,
        )
        self.update_service_member(service_member)

    @seq_task(5)
    def create_your_name(self):
        model = self.swagger_internal.get_model("PatchServiceMemberPayload")
        payload = model(
            first_name="Alice",  # Random
            middle_name="Carol",
            last_name="Bob",  # Random
            suffix="",
        )
        service_member = swagger_request(
            self.swagger_internal.service_members.patchServiceMember,
            serviceMemberId=self.user["service_member"].id,
            patchServiceMemberPayload=payload,
        )
        self.update_service_member(service_member)

    @seq_task(6)
    def create_your_contact_info(self):
        model = self.swagger_internal.get_model("PatchServiceMemberPayload")
        payload = model(
            email_is_preferred=True,
            personal_email=self.user["email"],  # Email is derived from logging in
            phone_is_preferred=True,
            secondary_telephone="333-333-3333",
            telephone="333-333-3333",
        )
        service_member = swagger_request(
            self.swagger_internal.service_members.patchServiceMember,
            serviceMemberId=self.user["service_member"].id,
            patchServiceMemberPayload=payload,
        )
        self.update_service_member(service_member)

    @seq_task(7)
    def search_for_duty_station(self):
        self.duty_stations = self.get_dutystations("buckley")

    @seq_task(8)
    def current_duty_station(self):
        model = self.swagger_internal.get_model("PatchServiceMemberPayload")
        payload = model(current_station_id=self.duty_stations[0].id)
        service_member = swagger_request(
            self.swagger_internal.service_members.patchServiceMember,
            serviceMemberId=self.user["service_member"].id,
            patchServiceMemberPayload=payload,
        )
        self.update_service_member(service_member)

    @seq_task(9)
    def current_residence_address(self):
        model = self.swagger_internal.get_model("PatchServiceMemberPayload")
        address = self.swagger_internal.get_model("Address")
        payload = model(
            residential_address=address(
                street_address_1="12345 Fake St",
                city="Aurora",
                state="CO",
                postal_code="80013",
            )
        )
        service_member = swagger_request(
            self.swagger_internal.service_members.patchServiceMember,
            serviceMemberId=self.user["service_member"].id,
            patchServiceMemberPayload=payload,
        )
        self.update_service_member(service_member)

    @seq_task(10)
    def backup_mailing_address(self):
        model = self.swagger_internal.get_model("PatchServiceMemberPayload")
        address = self.swagger_internal.get_model("Address")
        payload = model(
            backup_mailing_address=address(
                street_address_1="12345 Fake St",
                city="Aurora",
                state="CO",
                postal_code="80013",
            )
        )
        service_member = swagger_request(
            self.swagger_internal.service_members.patchServiceMember,
            serviceMemberId=self.user["service_member"].id,
            patchServiceMemberPayload=payload,
        )
        self.update_service_member(service_member)

    @seq_task(11)
    def backup_contact(self):
        model = self.swagger_internal.get_model(
            "CreateServiceMemberBackupContactPayload"
        )
        payload = model(
            name="Alice",
            email="alice@example.com",
            permission="NONE",
            telephone="333-333-3333",
        )
        swagger_request(
            self.swagger_internal.backup_contacts.createServiceMemberBackupContact,
            serviceMemberId=self.user["service_member"].id,
            createBackupContactPayload=payload,
        )

    #
    # At this point the user profile is complete so let's refresh our knowledge
    # of the user's profile.
    #

    @seq_task(12)
    def refresh_user_profile(self):
        self.update_user()

    #
    # Start adding move orders
    #

    @seq_task(13)
    def move_orders(self):
        # Get new duty station
        self.new_duty_stations = self.get_dutystations("travis")

        # Determine a few things randomly
        issue_date = datetime.datetime.now() + datetime.timedelta(
            days=random.randint(0, 30)
        )
        report_by_date = issue_date + datetime.timedelta(days=random.randint(60, 90))
        spouse_has_pro_gear = False
        has_dependents = bool(random.getrandbits(1))
        if has_dependents:
            spouse_has_pro_gear = bool(random.getrandbits(1))

        model = self.swagger_internal.get_model("CreateUpdateOrders")
        payload = model(
            service_member_id=self.user["service_member"].id,
            issue_date=issue_date,
            report_by_date=report_by_date,
            orders_type="PERMANENT_CHANGE_OF_STATION",
            has_dependents=has_dependents,
            spouse_has_pro_gear=spouse_has_pro_gear,
            new_duty_station_id=self.new_duty_stations[0].id,
        )
        self.orders = swagger_request(
            self.swagger_internal.orders.createOrders, createOrders=payload)

    @seq_task(14)
    def upload_orders(self):
        with open(os.path.join(self.fixtures_path, "test.pdf"), "rb") as f:
            swagger_request(
                self.swagger_internal.uploads.createUpload,
                documentId=str(uuid.uuid4()),
                file=f,
            )

    #
    # At this point orders have been uploaded so let's refresh our knowledge
    # of the user's profile.
    #

    @seq_task(15)
    def refresh_user_profile_2(self):
        self.update_user()

    #
    # Create the PPM Move
    #

    @seq_task(16)
    def select_ppm_move(self):
        model = self.swagger_internal.get_model("PatchMovePayload")
        payload = model(selected_move_type="PPM")
        pprint.pprint(self.orders)
        swagger_request(
            self.swagger_internal.moves.patchMove,
            moveId=self.orders[0].moves[0].id,
            patchMovePayload=payload,
        )

    @seq_task(17)
    def ppm_dates_and_locations(self):
        self.original_move_date = datetime.datetime.now() + datetime.timedelta(
            days=random.randint(30, 60)
        )
        swagger_request(
            self.swagger_internal.ppm.showPPMEstimate,
            original_move_date=self.original_move_date,
            origin_zip="80013",
            destination_zip="94535",
            weight_estimate=11500,
        )

    @seq_task(18)
    def create_ppm(self):
        model = self.swagger_internal.get_model("CreatePersonallyProcuredMovePayload")
        payload = model(
            days_in_storage=None,
            destination_postal_code="94535",
            has_additional_postal_code=False,
            has_sit=False,
            original_move_date=self.original_move_date,
            pickup_postal_code="80013",
        )
        swagger_request(
            self.swagger_internal.ppm.createPersonallyProcuredMove,
            moveId=self.orders[0].moves[0].id,
            createPersonallyProcuredMovePayload=payload,
        )

    @seq_task(19)
    def logout(self):
        self.client.post("/auth/logout")
        self.login_gov_user = None
        self.session_token = None
        self.user = {}
