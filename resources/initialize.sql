-- storidb.accountrepo definition

CREATE TABLE `account` (
                           `account_id` varchar(20) NOT NULL,
                           `name` varchar(100) DEFAULT NULL,
                           `email` varchar(50) DEFAULT NULL,
                           PRIMARY KEY (`account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


-- storidb.account_transaction definition

CREATE TABLE `account_transaction` (
                                       `account_id` varchar(20) NOT NULL,
                                       `txn_id` int(10) unsigned NOT NULL,
                                       `amount_credit` decimal(10,2) NOT NULL,
                                       `amount_debit` decimal(10,2) NOT NULL,
                                       `date` datetime NOT NULL,
                                       PRIMARY KEY (`account_id`,`txn_id`),
                                       CONSTRAINT `account_transaction_FK` FOREIGN KEY (`account_id`) REFERENCES `account` (`account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- INSERTS

INSERT INTO `account`
(account_id, name, email)
VALUES('GHM54789345', 'Günter Hagedorn', 'gunterhm@gmail.com');
