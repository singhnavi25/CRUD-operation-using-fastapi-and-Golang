import psycopg2
from fastapi import FastAPI
from pydantic import BaseModel, Field
from typing import List, Dict

cur = None
conn = None
table_name = "BestEmployees"

# these models used to define the returned results.
class Employee_Data(BaseModel):
    id: str = Field(..., example="E1001")
    name: str = Field(..., example="Manjot Singh")
    password: str = Field(..., example="3&3e(Ul0@")
    salary: str = Field(..., example="1000")
    age: str = Field(..., example="20")

class Update_Employee_Data(BaseModel):
    name: str = Field(..., example="Enter name")
    password: str = Field(..., example="Enter password")
    salary: str = Field(..., example="Enter salary")
    age: str = Field(..., example="Enter age")

def get_data_as_json(lt):
    ans = []

    for row in lt:
        temp = {'id' : row[0], 'name': row[1], 'password': row[2], 'salary': row[3], 'age': row[4]}
        ans.append(temp)
    return ans

app = FastAPI()

@app.on_event("startup")
async def startup():
    global conn
    conn = psycopg2.connect(
        database="", # database name
        user = "", # user name
        password = "", #password
        host = "", # host
        port = "" # port numbe enabled by you
    )
    global cur
    cur = conn.cursor()

    # Checking a table is present or not in the database 
    # if table present, it returns the table is exist otherwise it create the table
    cur.execute(f"SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = '{table_name.lower()}'")
    print(cur.fetchall())
    if bool(cur.rowcount) :
        print('Table is already exist')
    else:
        # Enter the table schema that you want to declare
        cur.execute(f'''CREATE TABLE {table_name}(
                ID VARCHAR(10) PRIMARY KEY,
                NAME VARCHAR(30),
                PASSWORD VARCHAR(30),
                SALARY VARCHAR(15),
                AGE VARCHAR(3)
                )''')
        conn.commit()
        print('created')
    # cur.execute(f"insert into {table_name} values('E1001', 'Manjot Singh', 'man&6@W', '10000', '20')")
    # cur.execute(f'select * from {table_name}')
    # print(cur.fetchall())
    # cur.execute(f'''SELECT * FROM {table_name} LIMIT 0''')
    # print(cur.description)
    
    # print(cur.description) # get metadata of database tables


@app.on_event("shutdown")
async def shutdown():
    await conn.commit()
    await conn.close()

@app.get('/getAllEmployees', response_model = List[Employee_Data])
async def all_employees():
    cur.execute(f"SELECT * FROM {table_name}")
    return get_data_as_json(cur.fetchall())

@app.get('/getEmployeeByID', response_model = List[Employee_Data])
async def get_employee_by_id(id: str):
    cur.execute(f"SELECT * FROM {table_name} where id='{id}'")
    return get_data_as_json(cur.fetchall())

@app.post('/registerEmployee', response_model = List[Employee_Data])
async def register_employee(employee : Employee_Data):
    cur.execute(f"INSERT INTO {table_name}(ID, NAME, PASSWORD, SALARY, AGE) VALUES(%s, %s, %s, %s, %s)", 
    (employee.id, employee.name, employee.password, employee.salary, employee.age))

    conn.commit()

    return await get_employee_by_id(employee.id)
    
@app.put('/updaterEmployeeData', response_model = List[Employee_Data])
async def update_employee_data(id: str, data_changed : Update_Employee_Data):
    cur.execute(f"UPDATE {table_name} SET NAME = %s, PASSWORD = %s, SALARY = %s, AGE = %s WHERE ID = %s", 
    (data_changed.name, data_changed.password, data_changed.salary, data_changed.age, id))

    conn.commit()
    return await get_employee_by_id(id)

@app.delete('/employeeDelete/{id}')
def delete_employee(id: str):
    # data = get_employee_by_id(id)
    cur.execute(f"DELETE FROM {table_name} WHERE ID = %s", (id,))
    conn.commit()
    return "Employee with " + id + " is deleted"
