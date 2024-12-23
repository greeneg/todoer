--
-- File generated with SQLiteStudio v3.4.4 on Mon Dec 23 13:24:35 2024
--
-- Text encoding used: UTF-8
--
PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- Table: Statuses
DROP TABLE IF EXISTS Statuses;

CREATE TABLE IF NOT EXISTS Statuses (
    Id         INTEGER PRIMARY KEY AUTOINCREMENT
                       UNIQUE
                       NOT NULL,
    StatusName STRING  UNIQUE
                       NOT NULL
);

INSERT INTO Statuses (
                         Id,
                         StatusName
                     )
                     VALUES (
                         1,
                         'new'
                     );

INSERT INTO Statuses (
                         Id,
                         StatusName
                     )
                     VALUES (
                         2,
                         'inprogress'
                     );

INSERT INTO Statuses (
                         Id,
                         StatusName
                     )
                     VALUES (
                         3,
                         'completed'
                     );


-- Table: Todos
DROP TABLE IF EXISTS Todos;

CREATE TABLE IF NOT EXISTS Todos (
    Id          INTEGER PRIMARY KEY AUTOINCREMENT
                        UNIQUE
                        NOT NULL,
    Description STRING  NOT NULL,
    Status      INTEGER REFERENCES Statuses (Id) 
                        NOT NULL
);


-- Table: Users
DROP TABLE IF EXISTS Users;

CREATE TABLE IF NOT EXISTS Users (
    Id              INTEGER  PRIMARY KEY AUTOINCREMENT
                             UNIQUE
                             NOT NULL,
    UserName        STRING   UNIQUE
                             NOT NULL,
    FullName        STRING   NOT NULL,
    PasswordHash    STRING   NOT NULL,
    Status          STRING   NOT NULL
                             DEFAULT enabled,
    CreationDate    DATETIME NOT NULL
                             DEFAULT (CURRENT_TIMESTAMP),
    LastChangedDate DATETIME NOT NULL
                             DEFAULT (CURRENT_TIMESTAMP) 
);

INSERT INTO Users (
                      Id,
                      UserName,
                      FullName,
                      PasswordHash,
                      Status,
                      CreationDate,
                      LastChangedDate
                  )
                  VALUES (
                      1,
                      'greeneg',
                      'Gary Greene',
                      'b584c299313f39097e3ba9c40a4859e3855496fd946905e3dec3c7bef177739e1e0d8dac2844831cf1388c2a6d91ff37829211216bf0b710cc1225388e690cf6',
                      'enabled',
                      '2024-12-23 17:59:03',
                      '2024-12-23 17:59:03'
                  );


COMMIT TRANSACTION;
PRAGMA foreign_keys = on;
