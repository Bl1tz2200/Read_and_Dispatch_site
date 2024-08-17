CREATE DATABASE Users_DB;

USE Users_DB;

CREATE TABLE Users (
	UserName VARCHAR(20) PRIMARY KEY,
    UserPassword VARCHAR(255),
    UserEmail VARCHAR(255) UNIQUE
   );
    
CREATE TABLE Users_Dispatches (
	DispatchID INT PRIMARY KEY,
    UserName VARCHAR(20)
    );
