CREATE TABLE IF NOT EXISTS users
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    username    VARCHAR(50)  NOT NULL,
    password    VARCHAR(100) NOT NULL,
    email       VARCHAR(100) NOT NULL,
    create_time DATETIME     NOT NULL,
    update_time DATETIME,
    UNIQUE INDEX (username)
);

CREATE TABLE IF NOT EXISTS user_role_type
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    code        VARCHAR(50) NOT NULL,
    disabled    tinyint     NOT NULL DEFAULT 0,
    create_time DATETIME    NOT NULL,
    update_time DATETIME
);

CREATE TABLE IF NOT EXISTS user_roles
(
    id                INT AUTO_INCREMENT PRIMARY KEY,
    user_id           INT      NOT NULL,
    user_role_type_id INT      NOT NULL,
    create_time       DATETIME NOT NULL,
    update_time       DATETIME,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (user_role_type_id) REFERENCES user_role_type (id)
);


CREATE TABLE IF NOT EXISTS finance_category
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    code        VARCHAR(50) NOT NULL,
    create_time DATETIME    NOT NULL,
    update_time DATETIME,
    UNIQUE INDEX (code)
);

INSERT IGNORE INTO finance_category(code, create_time)
VALUES ('飲食', NOW()),
       ('交通', NOW()),
       ('娛樂', NOW()),
       ('醫療', NOW()),
       ('生活', NOW()),
       ('其他', NOW());


CREATE TABLE IF NOT EXISTS user_finance_category
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    users_id    INT         NOT NULL,
    code        VARCHAR(50) NOT NULL,
    create_time DATETIME    NOT NULL,
    update_time DATETIME
);

CREATE TABLE IF NOT EXISTS user_finance_record
(
    id                       INT AUTO_INCREMENT PRIMARY KEY,
    users_id                 INT      NOT NULL,
    user_finance_category_id INT      NOT NULL,
    price                    INT      NOT NULL,
    spend_date               DATE     NOT NULL,
    create_time              DATETIME NOT NULL,
    update_time              DATETIME,
    FOREIGN KEY (users_id) REFERENCES users (id),
    FOREIGN KEY (user_finance_category_id) REFERENCES user_finance_category (id)
);



