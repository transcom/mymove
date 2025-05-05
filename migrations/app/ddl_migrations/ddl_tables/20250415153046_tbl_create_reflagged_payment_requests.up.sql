--B-22761 Maria Traskowsky added reflagged_payment_requests table for use in flag_sent_to_gex_for_review function
CREATE TABLE IF NOT EXISTS reflagged_payment_requests (
    payment_request_number TEXT NOT NULL PRIMARY KEY,
    reflagged_count INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    created_at TIMESTAMP NOT NULL DEFAULT now()
);