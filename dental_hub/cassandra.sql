CREATE TABLE dentists (
   username varchar,
   email varchar,
   password varchar,
   registrationdate timestamp,
   patients list<varchar>,
   
   PRIMARY KEY (email),
);

CREATE TABLE dentists (
   id uuid,
   username varchar,
   email varchar,
   password varchar,
   registrationdate timestamp,
   patients list<varchar>,
   
   PRIMARY KEY (email),
);

CREATE TABLE dentists_by_id (
    id uuid PRIMARY KEY,
    email text,
);

CREATE TABLE dentistresetpassword (
   dentistemail varchar,
   code varchar,
   expirationdate timestamp,
   PRIMARY KEY (dentistemail)
); 

CREATE TABLE dentistsignup (
   email varchar,
   username varchar,
   password varchar,
   verificationid uuid,
   expirationdate timestamp,
   PRIMARY KEY (verificationid, expirationdate)
); 

CREATE TABLE diagnosis (
   partitionid varchar,	
   id int,
   diagnosisname varchar,
   changestatus int,
   PRIMARY KEY (partitionid, id)
); 

CREATE TABLE manipulations (
   partitionid varchar,
   id int,
   manipulationname varchar,
   changestatus int,
   PRIMARY KEY (partitionid, id)
); 

CREATE TABLE toothstatus (
   partitionid varchar,
   id int,
   status varchar,
   PRIMARY KEY (partitionid,id)
); 

CREATE TYPE appointment (
    date timestamp,
    dentistemail varchar
);

CREATE TABLE schedule (
   dentistemail varchar,
   day timestamp,
   appointments list<frozen<appointment>>,
   
   PRIMARY KEY ((dentistemail, day))
);

CREATE TYPE toothaction (
    recordid uuid,
    operationid int,
    date timestamp
);

CREATE TYPE tooth (
    toothno varchar,
    diagnosislist set<frozen<toothaction>>,
    manipulationlist set<frozen<toothaction>>
);

CREATE TABLE patients (
   id uuid,
   email varchar,
   firstname varchar,
   middlename varchar,
   lastname varchar,
   address varchar,
   phonenumber varchar,
   generalinfo varchar, 
   registrationdate timestamp,
   dentists list<varchar>,
   teeth list<frozen<tooth>>,
   
   PRIMARY KEY (email)
);

CREATE TABLE patient_by_id (
    id uuid PRIMARY KEY,
    email text
)

