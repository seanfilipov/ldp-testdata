COMMENT ON SCHEMA normal IS 'Extra tables used for denormalization';

COMMENT ON SCHEMA loading IS 'Internal area used for data loading';

COMMENT ON TABLE loading.exlock IS
    'Exclusive lock to prevent concurrent data loads';

COMMENT ON TABLE users IS 'User records';

COMMENT ON COLUMN users.user_key IS 'Primary key of user record';
COMMENT ON COLUMN users.user_id IS 'FOLIO ID of the user';
COMMENT ON COLUMN users.username IS 'The user''s unique name, typically used for login';
COMMENT ON COLUMN users.barcode IS 'The library barcode for the user';
COMMENT ON COLUMN users.user_type IS 'The user class';
COMMENT ON COLUMN users.active IS
        'A flag to determine if a user can log in, take out loans, etc.';
COMMENT ON COLUMN users.group_name IS 'The group that the user is associated with';
COMMENT ON COLUMN users.group_description IS 'A description of the associated group';
COMMENT ON COLUMN users.record_effective IS
        'Date and time when the user record becomes effective';

COMMENT ON TABLE loans IS 'Loan transactions';

COMMENT ON COLUMN loans.loan_key IS 'Primary key of loan transaction';
COMMENT ON COLUMN loans.loan_id IS 'FOLIO ID of the loan transaction';
COMMENT ON COLUMN loans.user_key IS 'Foreign key of the user record for this loan';
-- COMMENT ON COLUMN loans.location_key IS '';
-- COMMENT ON COLUMN loans.item_key IS '';
COMMENT ON COLUMN loans.action IS 'Last action performed on the loan';
COMMENT ON COLUMN loans.status_name IS 'Overall status of the loan';
COMMENT ON COLUMN loans.loan_date IS 'Date and time when the loan began';
COMMENT ON COLUMN loans.due_date IS 'Date and time when the item is due to be returned';


