import os, sys, uuid, pandas as pd

current_dir = os.getcwd()

filename = f"{current_dir}/migrations/app/schema/20230809180036_update_transportation_offices.up.sql"

f = open(filename, "w")

if len(sys.argv) < 2:
    sys.exit("Input file required.")

df = pd.read_excel(pd.ExcelFile(sys.argv[1]))
df = df.reset_index()

<<<<<<< HEAD
f.write("-- Generated programmatically by load-transportation-offices.py\n\n")

=======
>>>>>>> 3394ea313a (Delete parent records for duty_locations being deleted.)
for index, row in df.iterrows():
    id = row["id"]
    name = row["name"]
    gbloc = row["GBLOC"]

    street_address_1 = row["street_address_1"]
    street_address_2 = row["street_address_2"]
    city = row["city"]
    state = row["state"]
    # trim to zip5
    postal_code = "%05d" % int(str(row["postal_code"]).split("-")[0])

    if pd.isna(id):
        id = uuid.uuid4()
        address_id = uuid.uuid4()
        f.write("-- Insert the address\n")
        f.write(
            f"INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at) VALUES ('{address_id}', '{street_address_1}', '{street_address_2}', '{city}', '{state}', '{postal_code}', now(), now());\n\n"
        )
        f.write("-- Insert the TO\n")
        f.write(
            f"INSERT INTO transportation_offices (id, address_id, gbloc, name, created_at, updated_at, latitude, longitude) VALUES ('{id}', '{address_id}', '{gbloc}', '{name}', now(), now(), 0, 0);\n\n"
        )
    else:
        f.write("-- Update the TO \n")
        f.write(
            f"UPDATE transportation_offices SET name = '{name}', gbloc = '{gbloc}' WHERE id = '{id}';\n\n"
        )
        f.write("-- Update the address\n")
        f.write(
            f"UPDATE addresses SET street_address_1 = '{street_address_1}', street_address_2 = '{street_address_2}', city = '{city}', state = '{state}', postal_code = '{postal_code}' WHERE id = (SELECT address_id FROM transportation_offices where id = '{id}');\n\n"
        )

f.close()
sys.exit()
