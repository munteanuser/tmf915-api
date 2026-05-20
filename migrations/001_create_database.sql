-- Run this on the master database first to create TMF915
IF NOT EXISTS (SELECT name FROM sys.databases WHERE name = N'TMF915')
BEGIN
    CREATE DATABASE TMF915;
END
GO
USE TMF915;
GO
