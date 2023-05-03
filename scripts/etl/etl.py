import sys
import pandas as pd
import numpy as np
import psycopg2
import psycopg2.extras

def get_chemical(lab, chemical) -> dict:

    lab = lab.strip()
    if lab == 'Lab 1' or lab == 'Lab 2':
        return  {
            "chemicalName": chemical[0],
            "chemicalNumber": chemical[1],
            "casNumber": chemical[2],
            "state": chemical[3],
            "quantity": chemical[4],
            "dateAdded": chemical[5],
            "expiryDate": chemical[6],
            "hazards": chemical[7],
            "sds": chemical[8],
            "coshh": chemical[9]
        }

    if lab == 'WMIC':
        return  {
            "chemicalName": chemical[0],
            "chemicalNumber": chemical[12],
            "casNumber": chemical[2],
            "state": chemical[3],
            "quantity": chemical[4],
            "dateAdded": chemical[5],
            "expiryDate": chemical[6],
            "hazards": chemical[7],
            "coshh": chemical[8],
            "sds": None
        }

    if lab == 'Lab 5':
        return  {
            "chemicalName": chemical[0],
            "chemicalNumber": None,
            "casNumber": chemical[1],
            "state": chemical[2],
            "quantity": chemical[3],
            "dateAdded": chemical[4],
            "expiryDate": chemical[5],
            "hazards": chemical[6],
            "coshh": chemical[7],
            "sds": None
        }

    return  {
        "chemicalName": chemical[0],
        "chemicalNumber": None,
        "casNumber": chemical[1],
        "state": chemical[2],
        "quantity": chemical[3],
        "dateAdded": chemical[4],
        "expiryDate": chemical[5],
        "hazards": chemical[6],
        "sds": chemical[7],
        "coshh": chemical[8],
    }

def insert_chemical(data_frame_row, lab_location):

    chemical = get_chemical(lab_location, data_frame_row)
    if chemical["state"] is not None:
        chemical["state"] = chemical["state"].lower()
    
    sql = """
        INSERT INTO coshh.chemical(
            cas_number,
            chemical_name,
            chemical_number,
            matter_state,
            quantity,
            added,
            expiry,
            safety_data_sheet,
            coshh_link,
            lab_location,
            storage_temp,
            is_archived 
        ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
        RETURNING id
    """

    cur = conn.cursor()

    try:
        cur.execute(sql, (
            chemical["casNumber"],
            chemical["chemicalName"],
            chemical["chemicalNumber"],
            chemical["state"] if chemical["state"] is not None and str(chemical["state"]).strip() != 'liq' else 'liquid',
            chemical["quantity"],
            chemical["dateAdded"],
            chemical["expiryDate"],
            chemical["sds"],
            chemical["coshh"],
            lab_location,
            'Shelf',
            False,
        ))
    except Exception as error:
        print("Failed to insert ", chemical["chemicalName"], ": ", error)
        conn.commit()
        cur.close()
        return 

    id = cur.fetchone()[0]
    
    insert_hazard(id)

    conn.commit()
    cur.close()

def insert_hazard(id):

    cur = conn.cursor()
    cur.execute("INSERT INTO coshh.chemical_to_hazard (id, hazard) VALUES (%s, %s)", (id, 'Unknown'))


def insert_chemicals(file_path):
    excel = pd.ExcelFile(file_path)

    for sheet_name in excel.sheet_names:

        if sheet_name == 'Chemical list': 
            continue

        df = pd.read_excel(excel, sheet_name=sheet_name)
        df = df.replace({np.nan: None})

        for i, row in df.iterrows():
            if i > 3 and row[0] is not None and row[0] is not pd.NaT: # skip headers; skip rows where chemical name is empty
                insert_chemical(row, sheet_name)

    conn.close()

file_path = sys.argv[1]
conn_string = sys.argv[2]
conn = None

try: 
    conn = psycopg2.connect(conn_string)
except (Exception) as error: 
    print(error)


insert_chemicals(file_path)
