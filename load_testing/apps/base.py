from locust import TaskSequence
from locust import task
from locust import events
from bravado_core.formatter import SwaggerFormat
from bravado.exception import HTTPError


def get_swagger_config():
    """
    Generate the config used in generating the swagger client from the spec
    """

    # MilMove uses custom formats for some fields. Without wanting to duplicate them here but
    # still wanting to not get warnings about them being undefined the UDFs are created here.
    # See https://bravado-core.readthedocs.io/en/stable/formats.html
    milmove_formats = []
    string_fmt_list = [
        "edipi",
        "ssn",
        "telephone",
        "uuid",
        "x-email",
        "zip",
    ]
    for fmt in string_fmt_list:
        swagger_fmt = SwaggerFormat(
            format=fmt,
            to_wire=str,
            to_python=str,
            validate=lambda x: x,
            description='Converts [wire]string:string <=> python string',
        )
        milmove_formats.append(swagger_fmt)
    swagger_config = {
        'validate_requests': False,
        'validate_responses': False,
        'formats': milmove_formats,
    }
    return swagger_config


def swagger_request(callable_operation, *args, **kwargs):
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


class BaseTaskSequence(TaskSequence):

    csrf = None

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


class InternalAPIMixin(object):
    swagger_internal = None

    @task(2)
    def load_swagger_file_internal(self):
        self.client.get("/internal/swagger.yaml")


class PublicAPIMixin(object):
    swagger_public = None

    @task(2)
    def load_swagger_file_public(self):
        self.client.get("/api/v1/swagger.yaml")
