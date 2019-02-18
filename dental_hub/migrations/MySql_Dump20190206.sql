CREATE DATABASE  IF NOT EXISTS `dental` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_bin */;
USE `dental`;
-- MySQL dump 10.13  Distrib 8.0.13, for Win64 (x86_64)
--
-- Host: dentalmysql.mysql.database.azure.com    Database: dental
-- ------------------------------------------------------
-- Server version	5.6.39.0

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
 SET NAMES utf8 ;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `dentist`
--

DROP TABLE IF EXISTS `dentist`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `dentist` (
  `Id` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `UserName` varchar(50) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `Email` varchar(50) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `Password` binary(60) NOT NULL,
  `RegistrationDate` datetime(6) DEFAULT NULL,
  PRIMARY KEY (`Id`),
  UNIQUE KEY `Id` (`Id`),
  UNIQUE KEY `IX_Dentist` (`Email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `dentistpatientinvitation`
--

DROP TABLE IF EXISTS `dentistpatientinvitation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `dentistpatientinvitation` (
  `InvitationId` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `DentistId` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `PatientId` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `ExpirationDate` datetime(6) DEFAULT NULL,
  PRIMARY KEY (`InvitationId`),
  UNIQUE KEY `InvitationId` (`InvitationId`),
  UNIQUE KEY `DentistId` (`DentistId`),
  UNIQUE KEY `PatientId` (`PatientId`),
  UNIQUE KEY `NonClusteredIndex-20180511-160851` (`DentistId`,`PatientId`),
  CONSTRAINT `FK_DentistPatientInvitation_Dentist` FOREIGN KEY (`DentistId`) REFERENCES `dentist` (`Id`),
  CONSTRAINT `FK_DentistPatientInvitation_Patient` FOREIGN KEY (`PatientId`) REFERENCES `patient` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `dentistpatientrelation`
--

DROP TABLE IF EXISTS `dentistpatientrelation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `dentistpatientrelation` (
  `DentistId` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `PatientId` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `Disabled` tinyint(1) DEFAULT NULL,
  `InvitationDate` datetime(6) DEFAULT NULL,
  PRIMARY KEY (`DentistId`,`PatientId`),
  UNIQUE KEY `DentistId` (`DentistId`),
  UNIQUE KEY `PatientId` (`PatientId`),
  CONSTRAINT `FK_DentistPatientRelation_Dentist` FOREIGN KEY (`DentistId`) REFERENCES `dentist` (`Id`),
  CONSTRAINT `FK_DentistPatientRelation_Patient` FOREIGN KEY (`PatientId`) REFERENCES `patient` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `diagnosis`
--

DROP TABLE IF EXISTS `diagnosis`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `diagnosis` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `DiagnosisName` varchar(100) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `ChangeStatus` int(11) NOT NULL,
  PRIMARY KEY (`Id`),
  KEY `FK_Diagnosis_ToothStatus` (`ChangeStatus`),
  CONSTRAINT `FK_Diagnosis_ToothStatus` FOREIGN KEY (`ChangeStatus`) REFERENCES `toothstatus` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `manipulations`
--

DROP TABLE IF EXISTS `manipulations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `manipulations` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `ManipulationName` varchar(50) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `ChangeStatus` int(11) NOT NULL,
  PRIMARY KEY (`Id`),
  KEY `FK_Manipulations_ToothStatus` (`ChangeStatus`),
  CONSTRAINT `FK_Manipulations_ToothStatus` FOREIGN KEY (`ChangeStatus`) REFERENCES `toothstatus` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `patient`
--

DROP TABLE IF EXISTS `patient`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `patient` (
  `Id` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `UserName` varchar(50) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `Email` varchar(50) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `Phone` varchar(50) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `Password` varchar(50) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `RegistrationDate` datetime(6) DEFAULT NULL,
  PRIMARY KEY (`Id`),
  UNIQUE KEY `Id` (`Id`),
  UNIQUE KEY `NonClusteredIndex-20180510-110418` (`Email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `patientinfo`
--

DROP TABLE IF EXISTS `patientinfo`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `patientinfo` (
  `Id` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `FirstName` varchar(50) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `MiddleName` varchar(50) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `LastName` varchar(50) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `Email` varchar(50) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `Address` varchar(100) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `PhoneNumber` varchar(100) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `GeneralInfo` longtext CHARACTER SET utf8 COLLATE utf8_bin,
  `RegistrationDate` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `DentistId` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `IsDeleted` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`Id`),
  UNIQUE KEY `Id` (`Id`),
  KEY `IX_PatientInfo` (`DentistId`),
  CONSTRAINT `FK_PatientInfo_Dentist` FOREIGN KEY (`DentistId`) REFERENCES `dentist` (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `resetpassword`
--

DROP TABLE IF EXISTS `resetpassword`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `resetpassword` (
  `DentistId` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `Code` char(6) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `ExpirationDate` datetime(6) NOT NULL,
  PRIMARY KEY (`DentistId`),
  UNIQUE KEY `DentistId` (`DentistId`),
  CONSTRAINT `FK_ResetPassword_Dentist` FOREIGN KEY (`DentistId`) REFERENCES `dentist` (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `resetpasswordpatient`
--

DROP TABLE IF EXISTS `resetpasswordpatient`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `resetpasswordpatient` (
  `PatientId` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `Code` char(6) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `ExpirationDate` datetime(6) NOT NULL,
  PRIMARY KEY (`PatientId`),
  UNIQUE KEY `PatientId` (`PatientId`),
  CONSTRAINT `FK_ResetPassword_Patient` FOREIGN KEY (`PatientId`) REFERENCES `patient` (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `signup`
--

DROP TABLE IF EXISTS `signup`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `signup` (
  `Email` varchar(50) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `UserName` varchar(50) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `Password` binary(60) NOT NULL,
  `VerificationId` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `ExpirationDate` datetime(6) NOT NULL,
  PRIMARY KEY (`Email`),
  UNIQUE KEY `VerificationId` (`VerificationId`),
  UNIQUE KEY `NonClusteredIndex-20180504-112347` (`VerificationId`,`ExpirationDate`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `signuppatient`
--

DROP TABLE IF EXISTS `signuppatient`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `signuppatient` (
  `Email` varchar(50) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `UserName` varchar(50) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `Password` binary(60) NOT NULL,
  `VerificationId` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `ExpirationDate` datetime(6) NOT NULL,
  PRIMARY KEY (`Email`),
  UNIQUE KEY `VerificationId` (`VerificationId`),
  KEY `NonClusteredIndex-20180510-111413` (`VerificationId`,`ExpirationDate`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `toothcurrentstatus`
--

DROP TABLE IF EXISTS `toothcurrentstatus`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `toothcurrentstatus` (
  `ToothID` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `ToothNo` char(2) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `HasCorona` tinyint(1) DEFAULT '0',
  `StatusId` int(11) DEFAULT NULL,
  `PatientId` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  PRIMARY KEY (`ToothID`),
  UNIQUE KEY `ToothNoPatientIndex` (`ToothNo`,`PatientId`),
  KEY `FK_ToothCurrentStatus_ToothStatus` (`StatusId`),
  KEY `FK_ToothPatientInfo_idx` (`PatientId`),
  CONSTRAINT `FK_ToothCurrentStatus_ToothStatus` FOREIGN KEY (`StatusId`) REFERENCES `toothstatus` (`id`),
  CONSTRAINT `FK_ToothPatientInfo` FOREIGN KEY (`PatientId`) REFERENCES `patientinfo` (`Id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `toothdiagnosis`
--

DROP TABLE IF EXISTS `toothdiagnosis`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `toothdiagnosis` (
  `Id` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `DiagnosisId` int(11) NOT NULL,
  `Date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `ToothId` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `IsDeleted` tinyint(1) DEFAULT NULL,
  `DateDeleted` datetime(6) DEFAULT NULL,
  PRIMARY KEY (`Id`),
  UNIQUE KEY `Id` (`Id`),
  KEY `ToothDiagnosisIndex` (`ToothId`,`IsDeleted`),
  KEY `FK_ToothDiagnosis_Diagnosis` (`DiagnosisId`),
  CONSTRAINT `FK_ToothDiagnosis_Diagnosis` FOREIGN KEY (`DiagnosisId`) REFERENCES `diagnosis` (`Id`),
  CONSTRAINT `FK_ToothDiagnosis_ToothCurrentStatus` FOREIGN KEY (`ToothId`) REFERENCES `toothcurrentstatus` (`ToothID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `toothmanipulation`
--

DROP TABLE IF EXISTS `toothmanipulation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `toothmanipulation` (
  `Id` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `ManipulationId` int(11) NOT NULL,
  `Date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `ToothId` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `IsDeleted` tinyint(1) DEFAULT NULL,
  `DateDeleted` datetime(6) DEFAULT NULL,
  PRIMARY KEY (`Id`),
  UNIQUE KEY `Id` (`Id`),
  KEY `ToothManipulationIndex` (`ToothId`,`IsDeleted`),
  KEY `FK_ToothManipulation_Manipulations` (`ManipulationId`),
  CONSTRAINT `FK_ToothManipulation_Manipulations` FOREIGN KEY (`ManipulationId`) REFERENCES `manipulations` (`Id`),
  CONSTRAINT `FK_ToothManipulation_ToothCurrentStatus` FOREIGN KEY (`ToothId`) REFERENCES `toothcurrentstatus` (`ToothID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `toothstatus`
--

DROP TABLE IF EXISTS `toothstatus`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `toothstatus` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Status` varchar(50) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping routines for database 'dental'
--
/*!50003 DROP PROCEDURE IF EXISTS `activate_invitation` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = '' */ ;
DELIMITER ;;
CREATE DEFINER=`sa`@`%` PROCEDURE `activate_invitation`(
	in invitationId char(36)
)
BEGIN
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN 
		ROLLBACK;
		RESIGNAL;
	END;
    
    START TRANSACTION;
    
	INSERT INTO `DentistPatientRelation` (`DentistId`, `PatientId`, `Disabled`, `InvitationDate`)
    SELECT `DentistId`, `PatientId`, 0, SYSDATE()
    FROM `DentistPatientInvitation`
    WHERE `InvitationId` = invitationId and `ExpirationDate` > SYSDATE();
    
    IF (SELECT row_count()) > 0 THEN
		DELETE FROM `DentistPatientInvitation`
		WHERE `InvitationId` = invitationId;
    END IF;
    
    COMMIT;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `add_diagnosis` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = '' */ ;
DELIMITER ;;
CREATE DEFINER=`sa`@`%` PROCEDURE `add_diagnosis`(
	in toothDiagnosisID char(36),
    in patientID char(36),
    in toothNo char(2),
    in diagnosisID int
)
BEGIN
	DECLARE toothID char(36);

	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN 
		ROLLBACK;
		RESIGNAL;
	END;
    
    START TRANSACTION;
    
    SET toothID =  (SELECT t.`ToothID`
					FROM `ToothCurrentStatus` t
					WHERE (t.`PatientId` = patientID and t.`ToothNo` = toothNo)
                    LIMIT 1 );
	
    IF (toothID is null) THEN
		SET toothID = UUID();
		INSERT INTO `ToothCurrentStatus` (`ToothID`, `ToothNo`, `PatientId`)
		VALUES(toothID, toothNo, patientID);
    
    END IF;
    
    INSERT INTO `ToothDiagnosis` (`Id`, `DiagnosisId`, `Date`, `ToothId`)
	VALUES (toothDiagnosisID, diagnosisID, SYSDATE(), toothID);

	COMMIT;	
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `add_manupulation` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = '' */ ;
DELIMITER ;;
CREATE DEFINER=`sa`@`%` PROCEDURE `add_manupulation`(
	in toothManipulationID char(36),
    in patientID char(36),
    in toothNo nchar(2),
    in manipulationID int
)
BEGIN
	DECLARE toothID varchar(64);
    
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN 
		ROLLBACK;
		RESIGNAL;
	END;
    
    START TRANSACTION;

	SELECT t.`ToothID`
    INTO toothID
	FROM `toothcurrentstatus` t
	WHERE (t.`PatientId`=patientID and t.`ToothNo`=toothNo)
	LIMIT 1;
    
    IF (toothID is null) THEN
    
		SET toothID = UUID();
		INSERT INTO `ToothCurrentStatus` (`ToothID`, `ToothNo`, `PatientId`)
		VALUES(toothID, toothNo, patientID);
    END IF;
    
    INSERT INTO `ToothManipulation` (`Id`, `ManipulationId`, `Date`, `ToothId`)
	VALUES (toothManipulationID, manipulationID, SYSDATE(), toothID);
    
    COMMIT;
                    
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `add_password_reset_confirmation_code` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = '' */ ;
DELIMITER ;;
CREATE DEFINER=`sa`@`%` PROCEDURE `add_password_reset_confirmation_code`(
	in p_email nvarchar(50),
    in p_code nchar(6)
)
BEGIN
	delete from `ResetPassword` where `DentistId` in (select `id` from `Dentist` where `Email` = p_email);
    
    INSERT INTO `ResetPassword` (`DentistId`, `Code`, `ExpirationDate`) 
    SELECT `Id`, p_code, DATE_ADD(SYSDATE(), INTERVAL 3 HOUR)
    from `Dentist` where `Email` = p_email;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `add_patient_password_reset_confirmation_code` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = '' */ ;
DELIMITER ;;
CREATE DEFINER=`sa`@`%` PROCEDURE `add_patient_password_reset_confirmation_code`(
	in p_email nvarchar(50),
    in p_code nchar(6)
)
BEGIN
	delete from ResetPasswordPatient where PaitentId in (select `id` from `Patient` where `Email` = p_email);
    
    INSERT INTO `ResetPasswordPatient` (PatientId, Code, ExpirationDate) 
    SELECT `Id`, p_code, DATE_ADD(SYSDATE(), INTERVAL 3 HOUR)
    from `Patient` where `Email` = p_email;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `get_teethData` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = '' */ ;
DELIMITER ;;
CREATE DEFINER=`sa`@`%` PROCEDURE `get_teethData`(
	in patientID char(36)
)
BEGIN

	SELECT t1.Id , t1.DiagnosisId, t1.Date, t2.ToothNo from `ToothDiagnosis` t1 
		inner join `ToothCurrentStatus` t2 
			on t1.ToothId = t2.ToothId 
	WHERE t2.PatientId = patientID and t1.IsDeleted is null
	ORDER BY t1.Date;
    
	SELECT t1.Id, t1.ManipulationId, t1.Date, t2.ToothNo from `ToothManipulation` t1 
		inner join `ToothCurrentStatus` t2 
			on t1.ToothId = t2.ToothId 
	WHERE t2.PatientId = patientID and t1.IsDeleted is null
	ORDER BY t1.Date;

END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `invite_patient` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = '' */ ;
DELIMITER ;;
CREATE DEFINER=`sa`@`%` PROCEDURE `invite_patient`(
	in dentistId char(36),
    in patientEmail nvarchar(50)
)
BEGIN
	DECLARE patientId char(36);
    DECLARE invitationId char(36);
    
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN 
		ROLLBACK;
		RESIGNAL;
	END;
    
    START TRANSACTION;
    
    SET patientId = (SELECT `Id` FROM `Patient` WHERE `Email` = patientEmail);
    
    IF (patientId is not null) THEN
		DELETE FROM `DentistPatientInvitation` WHERE `DentistId`= dentistId AND `PatientId`= patientId;
        
        SET invitationId = UUID();
        
        INSERT INTO `DentistPatientInvitation`(`InvitationId`, `DentistId`, `PatientId`, `ExpirationDate` ) 
		VALUES (invitationId, dentistId, patientId , DATE_ADD(SYSDATE(), INTERVAL 24 HOUR));
        
        SELECT invitationId;
    
    END IF;
    
    COMMIT;

END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `remove_tooth_diagnosis` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = '' */ ;
DELIMITER ;;
CREATE DEFINER=`sa`@`%` PROCEDURE `remove_tooth_diagnosis`(
	in toothDiagnosisID char(36)
)
BEGIN

	UPDATE `ToothDiagnosis` 
	SET `IsDeleted` = 1, `DateDeleted` = SYSDATE()
	WHERE `Id` =toothDiagnosisID;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `remove_tooth_manipulation` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = '' */ ;
DELIMITER ;;
CREATE DEFINER=`sa`@`%` PROCEDURE `remove_tooth_manipulation`(
	in toothManipulationID char(36)
)
BEGIN
	UPDATE `ToothManipulation`
	SET `IsDeleted` = 1, `DateDeleted` = SYSDATE()
	WHERE `Id` = toothManipulationID;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `reset_password_patient_sp` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = '' */ ;
DELIMITER ;;
CREATE DEFINER=`sa`@`%` PROCEDURE `reset_password_patient_sp`(
	in password binary(60),
    in email nvarchar(50),
    in code nchar(6)
)
BEGIN
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN 
		ROLLBACK;
		RESIGNAL;
	END;
    
    START TRANSACTION;
    
    UPDATE `Patient` d
    JOIN `ResetPasswordPatient` r
		ON r.PatientId = d.id
	SET `Password` = password
	WHERE d.Email = email and r.Code = code and r.ExpirationDate > SYSDATE();
    
    IF (SELECT row_count()) > 0 THEN
		DELETE FROM `ResetPasswordPatient` 
        WHERE `PatientId` IN (SELECT id FROM `Patient` WHERE `Email` = email);
    END IF;

	COMMIT;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `reset_password_sp` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = '' */ ;
DELIMITER ;;
CREATE DEFINER=`sa`@`%` PROCEDURE `reset_password_sp`(
	in password binary(60),
    in email nvarchar(50),
    in code nchar(6)
)
BEGIN
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN 
		ROLLBACK;
		RESIGNAL;
	END;
    
    START TRANSACTION;
    
    UPDATE `Dentist` d
		JOIN `ResetPassword` r
			ON r.DentistId = d.id
	SET `Password` = password
	WHERE d.Email = email and r.Code = code and r.ExpirationDate > SYSDATE();
    
    IF (SELECT row_count()) > 0 THEN
		DELETE FROM `ResetPassword` 
        WHERE DentistId IN (SELECT `id` FROM `Dentist` WHERE `Email` = email);
    END IF;
    
	COMMIT;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `signup_activate` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = '' */ ;
DELIMITER ;;
CREATE DEFINER=`sa`@`%` PROCEDURE `signup_activate`(
	in verificationId nvarchar(50)
)
BEGIN
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN 
		ROLLBACK;
		RESIGNAL;
	END;
    
    START TRANSACTION;
    
    INSERT INTO `Dentist` (`Id`, `UserName`, `Email`, `Password`)
	SELECT uuid(), `UserName`, `Email`, `Password`
	FROM `SignUp`
	WHERE `VerificationId` = verificationId and `ExpirationDate` > SYSDATE();
    
    IF (SELECT row_count()) > 0 THEN
		DELETE FROM `SignUp`
		WHERE `VerificationId` = verificationId;
    END IF;
    
    SELECT row_count();

	COMMIT;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `signup_patient_activate` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = '' */ ;
DELIMITER ;;
CREATE DEFINER=`sa`@`%` PROCEDURE `signup_patient_activate`(
	in verificationId nvarchar(50)
)
BEGIN
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN 
		ROLLBACK;
		RESIGNAL;
	END;
    
    START TRANSACTION;
    
    INSERT INTO `Patient` (`Id`, `UserName`, `Email`, `Password`)
	SELECT uuid(), `UserName`, `Email`, `Password`
	FROM `SignUpPatient`
	WHERE `VerificationId` = verificationId and `ExpirationDate` > SYSDATE();
    
    IF (SELECT row_count()) > 0 THEN
		DELETE FROM `SignUpPatient`
		WHERE `VerificationId` = verificationId;
    END IF;

	COMMIT;
    
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `signup_patient_register` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = '' */ ;
DELIMITER ;;
CREATE DEFINER=`sa`@`%` PROCEDURE `signup_patient_register`(
	in email nvarchar(50),
    in userName nvarchar(50),
    in password binary(60)
)
BEGIN
	DECLARE verificationId char(36) DEFAULT uuid();

    DELETE FROM `SignUpPatient` WHERE `Email` = email;
    
    INSERT INTO `SignUpPatient` (`Email`, `UserName`, `Password`, `VerificationId`, `ExpirationDate`) 
	VALUES (email, userName, password, verificationId, DATE_ADD(SYSDATE(), INTERVAL 3 HOUR));
    
    SELECT verificationId;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `signup_register` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8 */ ;
/*!50003 SET character_set_results = utf8 */ ;
/*!50003 SET collation_connection  = utf8_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = '' */ ;
DELIMITER ;;
CREATE DEFINER=`sa`@`%` PROCEDURE `signup_register`(
	in email nvarchar(50),
    in userName nvarchar(50),
    in password binary(60)
)
BEGIN
	DECLARE verificationId char(36) DEFAULT uuid();
    
	DELETE FROM `SignUp` WHERE `Email` = email;
    
    INSERT INTO `SignUp` (`Email`, `UserName`, `Password`, `VerificationId`, `ExpirationDate`) 
	VALUES (email, userName, password, verificationId, DATE_ADD(SYSDATE(), INTERVAL 3 HOUR));
    
    SELECT verificationId;

END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2019-02-06 11:25:08
