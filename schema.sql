CREATE TABLE users (
  id   SERIAL  PRIMARY KEY NOT NULL ,
  name  text    NOT NULL,
  phone_number  text   UNIQUE    NOT NULL,
  otp   text    NULL,
  otp_expiration_time   timestamp
);
