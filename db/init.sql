CREATE DATABASE RETIREMENTAPP;
USE RETIREMENTAPP;


CREATE TABLE customers (
  name VARCHAR(100),
  customerid VARCHAR(36) NOT NULL PRIMARY KEY,
  password VARCHAR(100)
);

insert into customers (name, customerid, password) values ('drew', UUID(), 'TOPSECRET');

CREATE TABLE sources (
  sourcename VARCHAR(1) NOT NULL PRIMARY KEY,
  sourcetype VARCHAR(36)
);
insert into sources (sourcename, sourcetype) values ('E', 'PRETAX');
insert into sources (sourcename, sourcetype) values ('2', 'ROTH');

CREATE TABLE plans (
  planid VARCHAR(36) NOT NULL PRIMARY KEY,
  planname VARCHAR(200),
  externalid VARCHAR(8)
);
insert into plans (planid, planname, externalid) values (UUID(), 'Penn State 403B Plan', '565657');
insert into plans (planid, planname, externalid) values (UUID(), 'Penn State 457B Plan', '565658');

CREATE TABLE enrollments (
  enrollmentid VARCHAR(36) NOT NULL PRIMARY KEY,
  deductionmethod INTEGER,
  planid VARCHAR(36), 
  FOREIGN KEY (planid) REFERENCES plans(planid),
  status VARCHAR(50),
  createtime DATETIME,
  customerid VARCHAR(36),
  FOREIGN KEY (customerid) REFERENCES customers(customerid)
);

CREATE TABLE deferrals (
  deferralid VARCHAR(36) NOT NULL PRIMARY KEY,
  sourcename VARCHAR(1),
  FOREIGN KEY (sourcename) REFERENCES sources(sourcename),
  deductamount INTEGER,
  createtime DATETIME,
  enrollmentid VARCHAR(36),
  FOREIGN KEY (enrollmentid) REFERENCES enrollments(enrollmentid)
);


