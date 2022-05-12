CREATE TABLE customer_support_remarks
(
    id uuid PRIMARY KEY NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    content text NOT NULL,
    office_user_id uuid NOT NULL,
    move_id uuid NOT NULL,
    CONSTRAINT fk_office_users FOREIGN KEY(office_user_id) REFERENCES office_users(id),
    CONSTRAINT fk_moves FOREIGN KEY (move_id) REFERENCES moves(id)
);

COMMENT on TABLE customer_support_remarks IS 'Store remarks from office users pertaining to moves.';
COMMENT on COLUMN customer_support_remarks.content IS 'Text content of the customer support remark written by an office user.';
COMMENT on COLUMN customer_support_remarks.move_id IS 'The move the customer support remark is associated with.';
COMMENT on COLUMN customer_support_remarks.office_user_id IS 'The office_user who authored the customer support remark.';

