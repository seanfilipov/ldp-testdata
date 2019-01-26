INSERT INTO d_users (id, username, barcode, user_type, active, group_name,
        description)
SELECT u.id,
       u.username,
       u.barcode,
       u.user_type,
       u.active,
       g.group_name,
       g.description group_description
    FROM users u
        LEFT JOIN groups g ON u.patron_group_id = g.id;

INSERT INTO d_locations (id, location_name)
SELECT 'id-' || replace(lower(tll.location_name), ' ', '-') id,
       tll.location_name
    FROM (
        SELECT DISTINCT location_name FROM tmp_loans_locations
    ) tll;
	
INSERT INTO f_loans (id, user_id, location_id, item_id, action, status_name,
        loan_date, due_date)
SELECT l.id,
       l.user_id,
       'id-' || replace(lower(tll.location_name), ' ', '-') location_id,
       l.item_id,
       l.action,
       l.status_name,
       l.loan_date,
       l.due_date
    FROM loans l
        LEFT JOIN tmp_loans_locations tll ON l.id = tll.loan_id;

