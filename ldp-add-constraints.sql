START TRANSACTION;


ALTER TABLE users
        ADD CONSTRAINT users_pkey
        PRIMARY KEY (user_id);


ALTER TABLE locations
        ADD CONSTRAINT locations_pkey
        PRIMARY KEY (location_id);


ALTER TABLE loans
        ADD CONSTRAINT loans_pkey
        PRIMARY KEY (loan_id);

ALTER TABLE loans
        ADD CONSTRAINT loans_user_id_fkey
        FOREIGN KEY (user_id)
        REFERENCES users(user_id);

ALTER TABLE loans
        ADD CONSTRAINT loans_location_id_fkey
        FOREIGN KEY (location_id)
        REFERENCES locations(location_id);


COMMIT;

