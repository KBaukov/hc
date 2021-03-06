CREATE TABLE `devices` (
  `ID` int(11) NOT NULL,
  `type` varchar(15) NOT NULL,
  `NAME` varchar(20) DEFAULT NULL,
  `IP` varchar(15) NOT NULL,
  `ACTIVE_FLAG` varchar(3) NOT NULL,
  `DESCRIPTION` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `PRIMARY_KEY_8` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `kotel` (
  `DEST_TP` double DEFAULT NULL,
  `DEST_TO` double DEFAULT NULL,
  `DEST_KW` int(11) DEFAULT NULL,
  `DEST_TC` double DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `map_sensors` (
  `ID` int(11) NOT NULL,
  `MAP_ID` int(11) NOT NULL,
  `DEVICE_ID` int(11) NOT NULL,
  `TYPE` varchar(20) NOT NULL,
  `XK` decimal(5,3) NOT NULL,
  `YK` decimal(5,3) NOT NULL,
  `DESCRIPTION` varchar(50) DEFAULT NULL,
  `PICT` varchar(50) NOT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `PRIMARY_KEY_1` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `maps` (
  `ID` smallint(6) NOT NULL,
  `TITLE` varchar(20) NOT NULL,
  `PICT` varchar(30) NOT NULL,
  `W` int(11) NOT NULL,
  `H` int(11) NOT NULL,
  `DESCRIPTION` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `PRIMARY_KEY_2` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `sessions` (
  `SESSION_ID` varchar(70) NOT NULL,
  `USER_ID` int(11) NOT NULL,
  `EXP_DATE` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `SESSIONS_SESSION_ID_IDX` (`SESSION_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `users` (
  `ID` int(11) NOT NULL,
  `LOGIN` varchar(15) NOT NULL,
  `PASS` varchar(70) NOT NULL,
  `ACTIVE_FLAG` char(3) NOT NULL,
  `USER_TYPE` varchar(10) DEFAULT NULL,
  `LAST_VISIT` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `PRIMARY_KEY_61` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;