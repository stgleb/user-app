CREATE DATABASE IF NOT EXISTS userapp;
use userapp;

CREATE TABLE IF NOT EXISTS users (
    Id           varchar(255) NOT NULL,
    Email        varchar(255) NOT NULL,
    Name         varchar(255) DEFAULT '',
    Address      varchar(255) DEFAULT '',
    Telephone    varchar(255) DEFAULT '',
    PasswordHash varchar(255) DEFAULT '',
    PRIMARY KEY(Id)
)   ENGINE = InnoDB DEFAULT CHARSET = latin1;

CREATE TABLE IF NOT EXISTS tokens (
    Value varchar(255) NOT NULL,
    Email varchar(255) NOT NULL,
    Used BOOL DEFAULT false NOT NULL,
    PRIMARY KEY (Value)
)   ENGINE = InnoDB DEFAULT CHARSET = latin1;