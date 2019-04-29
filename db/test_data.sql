-- password is password
INSERT INTO users(username, password, first_name, last_name) 
VALUES('test', '$2a$14$y9LYV/xaUkkwM8oEWvj1lOGwb28H4sCouwutJipY4va/QWqCeEYQ.', 'Test', 'User');

INSERT INTO events(title, start_time, end_time, owner_id)
VALUES ('Event 1', '2019-04-18 2:00:00', '2019-04-18 3:00:00', 1);

INSERT INTO events(title, start_time, end_time, owner_id)
VALUES ('Event 2', '2019-04-01 5:00:00', '2019-04-01 7:00:00', 1);

INSERT INTO events(title, start_time, end_time, owner_id)
VALUES ('Event 3', '2019-03-22 12:00:00', '2019-03-23 12:00:00', 1);

INSERT INTO events(title, start_time, end_time, owner_id)
VALUES ('Event 4', '2019-03-02 9:00:00', '2019-03-02 10:00:00', 1);

INSERT INTO events(title, start_time, end_time, owner_id)
VALUES ('Event 5', '2019-02-22 12:00:00', '2019-03-01 12:00:00', 1);