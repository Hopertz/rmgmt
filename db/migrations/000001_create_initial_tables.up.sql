CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE EXTENSION citext;

CREATE TABLE IF NOT EXISTS admin (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    activated bool NOT NULL,
    version UUID NOT NULL DEFAULT uuid_generate_v4()
);

CREATE TABLE IF NOT EXISTS token (
    hash bytea PRIMARY KEY,
    id UUID NOT NULL REFERENCES admin(id) ON DELETE CASCADE,
    expiry timestamp(0) with time zone NOT NULL,
    scope text NOT NULL
    
);

CREATE TABLE IF NOT EXISTS house (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    location citext NOT NULL,
    block citext NOT NULL,
    partition SMALLINT NOT NULL,
    occupied BOOL NOT NULL,
    version UUID NOT NULL DEFAULT uuid_generate_v4()
); 


CREATE TABLE IF NOT EXISTS tenant(
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    phone TEXT NOT NULL,
    house_id UUID NOT NULL REFERENCES house(id) ON DELETE CASCADE,
    personal_id_type TEXT NOT NULL DEFAULT '',
    personal_id TEXT NOT NULL DEFAULT '',
    photo TEXT NOT NULL DEFAULT '',
    active BOOL NOT NULL ,
    sos DATE NOT NULL,
    eos DATE NOT NULL,
    version UUID NOT NULL DEFAULT uuid_generate_v4()

); 

CREATE TABLE IF NOT EXISTS payment (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    tenant_id  UUID NOT NULL REFERENCES tenant(id) ON DELETE CASCADE,
    period int NOT NULL,
    start_date DATE NOT NULL,
    renewed BOOL NOT NULL,
    end_date DATE NOT NULL,
    version UUID NOT NULL DEFAULT uuid_generate_v4()
);

