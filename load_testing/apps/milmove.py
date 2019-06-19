# -*- coding: utf-8 -*-
import datetime
import os
import random
from urllib.parse import urljoin

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
    move = None
    ppm = None
    entitlements = None
    rank = None
    allotment = None
    zip_origin = "32168"
    zip_destination = "78626"

    def update_user(self):
        self.user = swagger_request(self.swagger_internal.users.showLoggedInUser)

    def update_service_member(self, service_member):
        self.user["service_member"] = service_member

    def get_dutystations(self, short_name):
        station_list = [short_name[0], short_name[0:3], short_name]
        duty_stations = None
        for station in station_list:
            duty_stations = swagger_request(
                self.swagger_internal.duty_stations.searchDutyStations, search=station
            )
        return duty_stations

    def get_move_id(self):
        return self.user["service_member"]["orders"][0]["moves"][0]["id"]

    def get_address_origin(self):
        address = self.swagger_internal.get_model("Address")
        return address(
            street_address_1="12345 Fake St",
            city="Crescent City",
            state="FL",
            postal_code=self.zip_origin,
        )

    def get_address_destination(self):
        address = self.swagger_internal.get_model("Address")
        return address(
            street_address_1="12345 Fake St",
            city="Austin",
            state="TX",
            postal_code=self.zip_destination,
        )

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

        # Need the entitlements to get ranks and weight allotments
        self.entitlements = swagger_request(
            self.swagger_internal.entitlements.indexEntitlements
        )
        self.rank = random.choice(list(self.entitlements.keys()))
        self.allotment = self.entitlements[self.rank]

        # Now set the profile
        model = self.swagger_internal.get_model("PatchServiceMemberPayload")
        payload = model(
            affiliation=random.choice(
                ["ARMY", "NAVY", "MARINES", "AIR_FORCE", "COAST_GUARD"]
            ),
            edipi=str(random.randint(10 ** 9, 10 ** 10 - 1)),
            social_security_number="333-33-3333",  # Random
            rank=self.rank,
        )
        service_member = swagger_request(
            self.swagger_internal.service_members.patchServiceMember,
            serviceMemberId=self.user["service_member"]["id"],
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
            serviceMemberId=self.user["service_member"]["id"],
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
            serviceMemberId=self.user["service_member"]["id"],
            patchServiceMemberPayload=payload,
        )
        self.update_service_member(service_member)

    @seq_task(7)
    def search_for_duty_station(self):
        self.duty_stations = self.get_dutystations("eglin")

    @seq_task(8)
    def current_duty_station(self):
        model = self.swagger_internal.get_model("PatchServiceMemberPayload")
        payload = model(current_station_id=self.duty_stations[0]["id"])
        service_member = swagger_request(
            self.swagger_internal.service_members.patchServiceMember,
            serviceMemberId=self.user["service_member"]["id"],
            patchServiceMemberPayload=payload,
        )
        self.update_service_member(service_member)

    @seq_task(9)
    def current_residence_address(self):
        model = self.swagger_internal.get_model("PatchServiceMemberPayload")
        payload = model(residential_address=self.get_address_origin())
        service_member = swagger_request(
            self.swagger_internal.service_members.patchServiceMember,
            serviceMemberId=self.user["service_member"]["id"],
            patchServiceMemberPayload=payload,
        )
        self.update_service_member(service_member)

    @seq_task(10)
    def backup_mailing_address(self):
        model = self.swagger_internal.get_model("PatchServiceMemberPayload")
        payload = model(backup_mailing_address=self.get_address_origin())
        service_member = swagger_request(
            self.swagger_internal.service_members.patchServiceMember,
            serviceMemberId=self.user["service_member"]["id"],
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
            serviceMemberId=self.user["service_member"]["id"],
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
            days=random.randint(0, 15)
        )
        report_by_date = issue_date + datetime.timedelta(days=random.randint(30, 60))
        spouse_has_pro_gear = False
        has_dependents = bool(random.getrandbits(1))
        if has_dependents:
            spouse_has_pro_gear = bool(random.getrandbits(1))

        model = self.swagger_internal.get_model("CreateUpdateOrders")
        payload = model(
            service_member_id=self.user["service_member"]["id"],
            issue_date=issue_date.date(),
            report_by_date=report_by_date.date(),
            orders_type="PERMANENT_CHANGE_OF_STATION",
            has_dependents=has_dependents,
            spouse_has_pro_gear=spouse_has_pro_gear,
            new_duty_station_id=self.new_duty_stations[0]["id"],
        )
        swagger_request(self.swagger_internal.orders.createOrders, createOrders=payload)

    @seq_task(14)
    def upload_orders(self):
        with open(os.path.join(self.fixtures_path, "test.pdf"), "rb") as f:
            swagger_request(self.swagger_internal.uploads.createUpload, file=f)

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
        swagger_request(
            self.swagger_internal.moves.patchMove,
            moveId=self.get_move_id(),
            patchMovePayload=payload,
        )

    @seq_task(17)
    def ppm_dates_and_locations(self):
        self.original_move_date = (
            datetime.datetime.now() + datetime.timedelta(days=random.randint(15, 30))
        ).date()
        swagger_request(
            self.swagger_internal.ppm.showPPMEstimate,
            original_move_date=self.original_move_date,
            origin_zip=self.zip_origin,
            destination_zip=self.zip_destination,
            weight_estimate=11500,  # This appears to be hard coded in the original API call
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
            pickup_postal_code=self.zip_origin,
        )
        self.ppm = swagger_request(
            self.swagger_internal.ppm.createPersonallyProcuredMove,
            moveId=self.get_move_id(),
            createPersonallyProcuredMovePayload=payload,
        )

    @seq_task(19)
    def select_ppm_weight(self):
        # Choose the TShirt size but don't update weight just yet
        model = self.swagger_internal.get_model("PatchPersonallyProcuredMovePayload")
        tshirt_size = random.choice(["S", "M", "L"])
        payload = model(size=tshirt_size, weight_estimate=0)

        # Sometimes the patch doesn't succeed because discount data is missing
        ppm_id = self.ppm["id"]
        new_ppm = swagger_request(
            self.swagger_internal.ppm.patchPersonallyProcuredMove,
            moveId=self.get_move_id(),
            personallyProcuredMoveId=ppm_id,
            patchPersonallyProcuredMovePayload=payload,
        )
        # This could mean that the PPM discount doesn't exist for the provided move dates
        if new_ppm is not None:
            self.ppm = new_ppm

        # Weights are decided by TShirt size
        allotment = self.allotment["total_weight_self"]
        size_weight = {
            "S": {
                "min": int(allotment * 0.01),
                "max": int(allotment * 0.10),
            },  # 1% - 10%
            "M": {
                "min": int(allotment * 0.10),
                "max": int(allotment * 0.25),
            },  # 10% - 25%
            "L": {"min": int(allotment * 0.25), "max": int(allotment)},  # 25% to max
        }

        weight_min = size_weight[tshirt_size]["min"]
        weight_max = size_weight[tshirt_size]["max"]

        # Make initial estimate call
        swagger_request(
            self.swagger_internal.ppm.showPPMEstimate,
            original_move_date=self.original_move_date,
            origin_zip=self.zip_origin,
            destination_zip=self.zip_destination,
            weight_estimate=int((weight_max - weight_min) / 2 + weight_min),
        )
        # Now modify the estimate within a random range
        weight_step = 5
        weight_estimate = random.randrange(weight_min, weight_max, weight_step)
        swagger_request(
            self.swagger_internal.ppm.showPPMEstimate,
            original_move_date=self.original_move_date,
            origin_zip=self.zip_origin,
            destination_zip=self.zip_destination,
            weight_estimate=weight_estimate,
        )
        payload = model(has_requested_advance=False, weight_estimate=weight_estimate)
        new_ppm_2 = swagger_request(
            self.swagger_internal.ppm.patchPersonallyProcuredMove,
            moveId=self.get_move_id(),
            personallyProcuredMoveId=ppm_id,
            patchPersonallyProcuredMovePayload=payload,
        )
        if new_ppm_2 is not None:
            self.ppm = new_ppm_2

    @seq_task(20)
    def update_move(self):
        self.move = swagger_request(
            self.swagger_internal.moves.showMove, moveId=self.get_move_id()
        )

    @seq_task(21)
    def validate_entitlements(self):
        self.move = swagger_request(
            self.swagger_internal.entitlements.validateEntitlement,
            moveId=self.get_move_id(),
        )

    @seq_task(22)
    def signature(self):
        model = self.swagger_internal.get_model("CreateSignedCertificationPayload")
        swagger_request(
            self.swagger_internal.certification.createSignedCertification,
            moveId=self.get_move_id(),
            createSignedCertificationPayload=model(
                date=datetime.datetime.now(),
                signature="ABC",
                certification_text="clatto verata necktie",
            ),
        )

    @seq_task(23)
    def submit_move(self):
        model = self.swagger_internal.get_model("SubmitMoveForApprovalPayload")
        swagger_request(
            self.swagger_internal.moves.submitMoveForApproval,
            moveId=self.get_move_id(),
            submitMoveForApprovalPayload=model(ppm_submit_date=datetime.datetime.now()),
        )

    @seq_task(24)
    def get_transportation_offices(self):
        swagger_request(
            self.swagger_internal.transportation_offices.showDutyStationTransportationOffice,
            dutyStationId=self.duty_stations[0]["id"],
        )
        swagger_request(
            self.swagger_internal.transportation_offices.showDutyStationTransportationOffice,
            dutyStationId=self.new_duty_stations[0]["id"],
        )

    @seq_task(25)
    def logout(self):
        self.client.post("/auth/logout")
        self.login_gov_user = None
        self.session_token = None
        self.user = {}
