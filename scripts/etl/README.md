Run the etl with `python3 etl.py <spreadsheetName> <dbConnectionString>.
For example:

python3 etl/etl.py /mnt/c/Users/alex.whitehorn/MDC\ Chemical\ inventory.xlsx "host=localhost port=5432 user=postgres password=postgres dbname=informatics sslmode=disable"

In case this needs to be run into wopr, credentials can be found in Rancher under wopr -> default -> secrets -> coshh-secret.