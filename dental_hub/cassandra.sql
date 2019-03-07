CREATE TABLE dentists (
   id uuid,
   username varchar,
   email varchar,
   password varchar,
   registrationdate timestamp,
   patients set<uuid>,
   
   PRIMARY KEY (email),
);

CREATE TABLE dentistresetpassword (
   dentistid uuid,
   code varchar,
   expirationdate timestamp,
   PRIMARY KEY (dentistid)
); 

CREATE TABLE dentistsignup (
   email varchar,
   username varchar,
   password varchar,
   verificationid uuid,
   expirationdate timestamp,
   PRIMARY KEY (email)
); 

CREATE TABLE diagnosis (
   id int,
   diagnosisname varchar,
   changestatus int,
   PRIMARY KEY (id)
); 

CREATE TABLE manipulations (
   id int,
   manipulationname varchar,
   changestatus int,
   PRIMARY KEY (id)
); 

CREATE TABLE toothstatus (
   id int,
   status varchar,
   PRIMARY KEY (id)
); 

CREATE TYPE appointment (
    date timestamp,
    dentistid uuid
);

CREATE TABLE schedule (
   dentistid uuid,
   day timestamp,
   appointments list<frozen<appointment>>,
   
   PRIMARY KEY ((dentistid, day))
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

