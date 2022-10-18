DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS vacations;
DROP TABLE IF EXISTS extra_days;
DROP TABLE IF EXISTS user_type_vacation;
--;
--;
-- Create a table users
CREATE TABLE users(
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL DEFAULT "",
    userType TEXT NOT NULL DEFAULT "Employee",
    flm INTEGER NOT NULL DEFAULT 0,
    email TEXT NOT NULL DEFAULT "",
    startDate NUMERIC DEFAULT (date('now', 'localtime')),
    extraDays NUMERIC NOT NULL DEFAULT 0,
    spillover NUMERIC NOT NULL DEFAULT 0,
    IsActive NUMERIC NOT NULL DEFAULT 1,
    created_at NUMERIC DEFAULT (datetime('now', 'localtime')),
    updated_at NUMERIC DEFAULT (datetime('now', 'localtime')),
    FOREIGN KEY (flm) REFERENCES users(id)
);
-- Create table extra_days
CREATE TABLE extra_days(
    extra_day NUMERIC NOT NULL DEFAULT 0,
    description TEXT NOT NULL DEFAULT "",
    created_at NUMERIC DEFAULT (datetime('now', 'localtime')),
    updated_at NUMERIC DEFAULT (datetime('now', 'localtime'))
);
--- create table for mapping available type vacation to user
CREATE TABLE user_type_vacation(
    userId INTEGER NOT NULL,
    TypeVacationId INTEGER NOT NULL,
    FOREIGN KEY(userId) REFERENCES users(id),
    FOREIGN KEY (TypeVacationId) REFERENCES type_vacation(id)
);
--- Create table vacations
CREATE TABLE vacations(
    id INTEGER PRIMARY KEY,
    userId INTEGER NOT NULL,
    typeId INTEGER NOT NULL,
    startDate NUMERIC DEFAULT (date('now', 'localtime')),
    duration NUMERIC DEFAULT 1,
    status INTEGER NOT NULL  DEFAULT 0,
    partOfBd INTEGER NOT NULL  DEFAULT 0,
    created_at NUMERIC DEFAULT (datetime('now', 'localtime')),
    updated_at NUMERIC DEFAULT (datetime('now', 'localtime')),
    CONSTRAINT vacations_user_FK FOREIGN KEY (userId) REFERENCES users(id)
);