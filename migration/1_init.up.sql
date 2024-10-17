CREATE TABLE geolocation
(
    ip_address    varchar(255) primary key,
    country_code  varchar(15) NOT NULL,
    country       varchar(255) NOT NULL,
    city          varchar(255) NOT NULL,
    latitude      numeric NOT NULL,
    longitude     numeric NOT NULL,
    -- Don't have enough data about mystery_value type & left it as text (not integer)
    mystery_value text NOT NULL
);
