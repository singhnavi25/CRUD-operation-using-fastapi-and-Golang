from fastapi import FastAPI
from typing import List
import databases, sqlalchemy
from pydantic import BaseModel, Field

# Enter connection constraints for database connection like database here is postgresql
url = "database://username:password" + "host" + "/database name"
dataBase = databases.Database(url)

# Metadata
metaData = sqlalchemy.MetaData()

employees = sqlalchemy.Table(
    "ManjotEmployees",
    metaData,
    sqlalchemy.Column("id", sqlalchemy.String, primary_key=True),
    sqlalchemy.Column("name", sqlalchemy.String),
    sqlalchemy.Column("password", sqlalchemy.String),
    sqlalchemy.Column("salary", sqlalchemy.String),
    sqlalchemy.Column("age", sqlalchemy.String),
) 

engine = sqlalchemy.create_engine(url)

metaData.create_all(engine)

# these models used to define the returned results.
class AllEmployeesList(BaseModel):
    id: str
    name: str
    password: str
    salary: str
    age: str

  

class insertEmployeeData(BaseModel):
    id: str = Field(..., example="E1001")
    name: str = Field(..., example="Manjot Singh")
    password: str = Field(..., example="3&3e(Ul0@")
    salary: str = Field(..., example="1000")
    age: str = Field(..., example="20")


class updateEmployeeData(BaseModel):
    id: str = Field(..., example="Enter id")
    name: str = Field(..., example="Manjot Singh")
    password: str = Field(..., example="3&3e(Ul0@")
    salary: str = Field(..., example="1000")
    age: str = Field(..., example="20")

class deleteEmployee(BaseModel):
    id: str = Field(..., example = "Enter employee id")

app = FastAPI()

@app.on_event("startup")
async def startup():
    await dataBase.connect()


@app.on_event("shutdown")
async def shutdown():
    await dataBase.disconnect()


@app.get("/employees", response_model = List[AllEmployeesList])
async def alldata():
    selectQuery = employees.select()
    return await dataBase.fetch_all(selectQuery)


@app.post("/employees", response_model = AllEmployeesList)
async def employeeRegister(employee: insertEmployeeData):
    registerQuery = employees.insert().values(
        id = employee.id,
        name = employee.name,
        password = employee.password,
        salary = employee.salary,
        age = employee.age
    )
    await dataBase.execute(registerQuery)
    return {
        "Employee registered" : "Yes",
        **employee.dict() 
    }


@app.get("/employees/{id}", response_model = AllEmployeesList)
async def getEmployeeFromId(id: str):
    selectQuery = employees.select().where(employees.c.id == id)
    return await dataBase.fetch_one(selectQuery)


@app.put("/users", response_model = AllEmployeesList)
async def updateEmployee(employee: updateEmployeeData):
    updateQuery = employees.update().where(employees.c.id == employee.id).values(
        name = employee.name,
        password = employee.password,
        salary = employee.salary,
        age = employee.age
    )
    await dataBase.execute(updateQuery)
    return getEmployeeFromId(employee.id)


@app.delete("/employees/{id}")
async def deleteEmployeeFromData(id:str, employee: deleteEmployee):
    delQuery = employees.delete().where(employees.c.id == id)
    await dataBase.execute(delQuery)
    return "Employee with " + id + " is deleted"