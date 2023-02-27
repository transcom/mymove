This endpoint updates a single Move by ID. This allows the Admin User to change
the `show` field on the selected field to either `True` or `False`. A "shown"
Move will appear to all users as normal, a "hidden" Move will not be returned or
editable using any other endpoint (besides those in the Support API), and thus
effectively deactivated. Do not use this endpoint directly as it is meant to be
used with the Admin UI exclusively.
