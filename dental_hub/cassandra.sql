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

CREATE TABLE schedule (
   dentistemail varchar,
   day timestamp,
   appointments text,
   
   PRIMARY KEY ((dentistemail, day))
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
   teeth_json text,
   
   PRIMARY KEY (email)
);

CREATE TABLE patients_by_id (
    id uuid PRIMARY KEY,
    email text
);

CREATE TABLE tag (
   
   patientid uuid ,
   tagkey varchar,
   tagvalue varchar,
   imageid varchar,
   
   PRIMARY KEY ((patientid), tagkey, tagvalue, imageid)
);

INSERT INTO manipulations(partitionid, id, manipulationname, changestatus) VALUES('manipulationsKey', 1, 'Caries himiopolimer', 6);
INSERT INTO manipulations(partitionid, id, manipulationname, changestatus) VALUES('manipulationsKey', 2, 'Caries fotopolimer', 6);

INSERT INTO diagnosis(partitionid, id, diagnosisname, changestatus) VALUES('diagnosisKey', 1, 'Caries superficialis', 3);
INSERT INTO diagnosis(partitionid, id, diagnosisname, changestatus) VALUES('diagnosisKey', 2, 'Caries media', 3);

INSERT INTO toothstatus(partitionid, id, status) VALUES('toothStatusKey', 1, 'Unchenged');
INSERT INTO toothstatus(partitionid, id, status) VALUES('toothStatusKey', 2, 'Extractio');
INSERT INTO toothstatus(partitionid, id, status) VALUES('toothStatusKey', 3, 'Caries');
INSERT INTO toothstatus(partitionid, id, status) VALUES('toothStatusKey', 4, 'Pulpitis');
INSERT INTO toothstatus(partitionid, id, status) VALUES('toothStatusKey', 5, 'Periodontitis');
INSERT INTO toothstatus(partitionid, id, status) VALUES('toothStatusKey', 6, 'Obturacio');



