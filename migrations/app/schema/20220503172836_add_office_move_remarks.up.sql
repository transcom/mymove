CREATE TABLE office_move_remarks
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

COMMENT on TABLE office_move_remarks IS 'Store remarks from office users pertaining to moves.';
COMMENT on COLUMN office_move_remarks.content IS 'Text content of the remark written by an office user.';
COMMENT on COLUMN office_move_remarks.move_id IS 'The move the office remark is associated with.';
COMMENT on COLUMN office_move_remarks.office_user_id IS 'The office_user who authored the office remark.';

