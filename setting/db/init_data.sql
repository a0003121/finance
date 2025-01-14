use go;

INSERT IGNORE INTO user_role_type(code, create_time)
VALUES ('ADMIN', NOW()),
       ('USER', NOW());

INSERT IGNORE INTO finance_category(code, create_time)
VALUES ('飲食', NOW()),
       ('交通', NOW()),
       ('娛樂', NOW()),
       ('醫療', NOW()),
       ('生活', NOW()),
       ('其他', NOW());