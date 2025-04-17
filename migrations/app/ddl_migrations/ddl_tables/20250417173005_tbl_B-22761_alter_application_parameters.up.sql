--B-22761 Maria Traskowsky adding flagSentToGexStaleInterval for use in flag_sent_to_gex_for_review function
INSERT into application_parameters (id, parameter_name, parameter_value)
VALUES (
        'b98715f0-0621-470b-9498-b054c967e7e7',
        'flagSentToGexStaleInterval',
        '12 hours'
    );