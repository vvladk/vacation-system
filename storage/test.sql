INSERT INTO vacations
(userId, typeId, startDate, duration)
VALUES(1, 1, date('now', 'localtime'), 1);

ALTER TABLE vacations ADD partOfBd INTEGER NOT NULL  DEFAULT 0;