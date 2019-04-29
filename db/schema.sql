CREATE SEQUENCE users_id_seq1;

CREATE TABLE users(
    id integer default nextval('users_id_seq1') NOT NULL, 
    username text NOT NULL, 
    password text NOT NULL, 
    first_name text NOT NULL, 
    last_name text NOT NULL
);

CREATE SEQUENCE events_id_seq;

CREATE TABLE events(
    id integer default nextval('events_id_seq') NOT NULL,
    title text NOT NULL,
    start_time timestamp without time zone NOT NULL,
    end_time timestamp without time zone NOT NULL,
    location text,
    notes text,
    owner_id integer NOT NULL
);
