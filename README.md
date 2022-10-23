## Purpose of the solution.
The Vacation management system is designed to ensure transparency in the accrual and management of vacations for employees of small IT companies. Takes into account the realities of Ukrainian IT.

## Brief description of the business process
#### There are three roles in the system:
1. **Employee** - the actual employee of the company.
1. **FLM** - ***First Line Manager*** - team / department leader who is responsible for coordinating **Employee** vacations.
1. **HR** is a manager who is directly responsible for the accouunting of **Employee** vacations.


Adding new roles is not yet planned.

#### There are three types of vacations in the system:
1. Paid vacation
1. Unpaid vacation
1. Sick leave

### Unpaid vacation
It has no fixed limit in the system. The system records the number of actually used days of unpaid vacation.
### Paid vacation
The system has a single limit for all **Employees** - ***14*** days from January 01 to December 31. It is also possible to charge an additional number of days for each **Employee**. Paid vacation is accrued on ***14/12 = 1.17*** days for each month worked.
The balance can be carried over to the next year

### Sick leave
The system has a single limit for all **Employees** - ***5*** days from January 01 to December 31. The entire limit is available starting January 1st. 
Does not carry over to the next year.

## Brief description of technical details
Authorization and authentication in the system is based on the use of a corporate Google account.

### Installation
1. `go build -o vms main.go`
1. create target folder *\<Vacation Management System\>*
1. copy the executable to *\<Vacation Management System\>*
create subfolders:
    1. *\<Vacation Management System\>*/ini
    1. *\<Vacation Management System\>*/storage
    1. *\<Vacation Management System\>*/ui
1. copy the contents of the ui folder from the repository to the *\<Vacation Management System\>*/ui
1. In the *\<Vacation Management System\>*/ini folder, place the **vs.ini** file created on the basis of example.ini, pre-filling it with configuration information.
    1. **auth** section - Google's OAuth configuration
    1. **session** section - session configuration for gorilla/sessions
1. `cd <Vacation Management System>/storage`
1. `sqlite3 vs. sqlite < <path_to_rep>/storage/init/01.create.tables.sql`
1. Insert at least 1 user with type **HR**<br> 
`INSERT INTO users(email, title, startDate, userType, flm) VALUES('<email>', '<User_title>','2021-12-31','HR', 0)`
1. Run executable
1. Log in to the system under the created user(**HR**)
1. Configure:
    1. Create accounts for company employees
    1. Add an extra days-off
