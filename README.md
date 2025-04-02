# rental-property-management-system

### Run Backend

```shell
  # under the root folder:
  go run backend/main.go
```

### Run Test Python Scripts

```shell
  # under the folder 'tests'
  # init and activate python virtual environment
  python -m venv .
  # for Windows
  ./Scripts/activate 
  # for MacOS: 
  # source bin/activate
  
  # install dependencies
  pip install -r requirements.txt

  # run
  python main.py
```

### Setup Postgres

```shell
  docker run --name some-postgres -p 5433:5432 -e POSTGRES_PASSWORD=270153 -d postgres
```

### Mock Data

1. Login with WeChat, set the role of the user to admin in the database manually
2. Insert house data into the database.
