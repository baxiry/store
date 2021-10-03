CREATE TABLE users (
    id INT NOT NULL AUTO_INCREMENT,
    username varchar(100),
    password varchar(100),
    email varchar(100),
    phon varchar(100),
    linkavatar varchar(255),
    primary key (id)
);

CREATE TABLE users (
  id INT NOT NULL AUTO_INCREMENT, username varchar(30), password varchar(30), email varchar(70), phon varchar(15),
  linkavatar varchar(70), primary key (id));


CREATE TABLE products (
    productID INT UNSIGNED NOT NULL AUTO_INCREMENT,
    ownerID INT NOT NULL, -- userID
    title varchar(150),
    -- productCode  CHAR(3) NOT NULL DEFAULT '',
    price 
    catigory
    quantity INT UNSIGNED  NOT NULL DEFAULT 1,
    photos varchar(1000), -- TODO splite photos as table;
    tamestamp TIMESTAMP,
    PRIMARY KEY (productID)
);
   
CREATE TABLE products (productID INT UNSIGNED NOT NULL AUTO_INCREMENT, ownerID INT UNSIGNED  NOT NULL, title varchar(150),
    quantity INT UNSIGNED  NOT NULL DEFAULT 1,photos varchar(1000),tamestamp TIMESTAMP,PRIMARY KEY (productID));
--  ------------- general info -----------------

--  add new column in spicial position:
ALTER TABLE table_name ADD new_column_name column_definition  [FIRST | AFTER column_name];
-- example:
ALTER TABLE products ADD description varchar(500) AFTER title;

-- remove table:
DROP TABLE table_name;

-- remove column:
ALTER TABLE products DROP COLUMN owner;

-- fix string encoding error
ALTER DATABASE <db_name> CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
ALTER TABLE <table_name> CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

--

create table tutorials_tbl( 
	usreid INT NOT NULL AUTO_INCREMENT,
	username VARCHAR(100) NOT NULL,
	password VARCHAR(40) NOT NULL,
	phon VARCHAR(14) NOT NULL,
	email VARCHAR(42)NOT NULL,
	PRIMARY KEY ( userid )
);


DELETE FROM table_name
WHERE condition;

db.QueryRow("DELETE FROM customer WHERE id=?", id)

