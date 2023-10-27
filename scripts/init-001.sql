DROP DATABASE IF EXISTS informatics;
CREATE DATABASE informatics;
\connect informatics
CREATE SCHEMA coshh;

CREATE TYPE coshh.matter_state AS ENUM ('liquid', 'solid');
CREATE TYPE coshh.hazard AS ENUM (
    'None',
    'Explosive',
    'Flammable',
    'Oxidising',
    'Corrosive',
    'Acute toxicity',
    'Hazardous to the environment',
    'Health hazard/Hazardous to the ozone layer',
    'Serious health hazard',
    'Gas under pressure',
    'Unknown'
    );

CREATE TYPE coshh.storage_temp AS ENUM (
    'Shelf',
    '+4',
    '-20',
    '-80'
    );

CREATE TABLE coshh.chemical
(
    id                SERIAL PRIMARY KEY,
    cas_number        VARCHAR(255),
    chemical_name     VARCHAR(255)       NOT NULL,
    chemical_number   VARCHAR(255),
    matter_state      coshh.matter_state,
    quantity          VARCHAR(255) DEFAULT '0',
    added             DATE,
    expiry            DATE,
    safety_data_sheet VARCHAR(255),
    coshh_link        VARCHAR(255),
    lab_location      VARCHAR(255),
    cupboard          VARCHAR(255),
    chemical_owner    VARCHAR(255),
    storage_temp      coshh.storage_temp NOT NULL,
    is_archived       BOOLEAN            NOT NULL
);

CREATE TABLE coshh.chemical_to_hazard
(
    id     INT,
    hazard coshh.hazard
);
