import os, sys, uuid, pandas as pd

current_dir = os.getcwd()

filename = (
    f"{current_dir}/migrations/app/schema/20230809180036_update_duty_locations.up.sql"
)

f = open(filename, "w")

if len(sys.argv) < 2:
    sys.exit("Input file required.")

f.write("-- Temporarily remove foreign key constraint\n")
f.write("ALTER TABLE duty_locations DROP CONSTRAINT duty_locations_address_id_fkey;\n")
f.write("DELETE from addresses where id in (select address_id from duty_locations);\n")
f.write("-- Remove existing duty location names\n")
f.write("DELETE FROM duty_location_names;\n\n")
f.write("-- Remove existing duty locations\n")
f.write("DELETE FROM duty_locations;\n\n")

df = pd.read_excel(pd.ExcelFile(sys.argv[1]), keep_default_na=False, na_values=[""])
df = df.reset_index()

for index, row in df.iterrows():
    name = row["Duty Location Name"]
    gbloc = row["gbloc"]
    affiliation = row["Affiliation"]
    street_address = row["Street Address"]
    city = row["City"]
    state = row["State"]
    postal_code = "%05d" % row["Postal Code"]
    to_id = "NULL"

    address_id = uuid.uuid4()
    id = uuid.uuid4()
    f.write("-- Insert the address\n")
    f.write(
        f"INSERT INTO addresses (id, street_address_1, city, state, postal_code, created_at, updated_at) VALUES ('{address_id}', '{street_address}', '{city}', '{state}', '{postal_code}', now(), now());\n\n"
    )

    if not pd.isna(gbloc):
        f.write("-- Find the transportation office by GBLOC, City, and State\n")
        f.write(
            f"WITH found_to AS (SELECT t.id FROM transportation_offices AS t, addresses AS a WHERE t.address_id = a.id AND t.gbloc='{gbloc}' AND a.city='{city}' AND a.state='{state}' AND t.shipping_office_id IS NOT NULL LIMIT 1)\n\n"
        )
        to_id = "(select id from found_to)"

    f.write("-- Insert the duty location\n")
    f.write(
        f"INSERT INTO duty_locations (id, address_id, transportation_office_id, name, affiliation, updated_at, created_at) VALUES('{id}', '{address_id}', {to_id}, '{name}', '{affiliation}', now(), now());\n\n"
    )

f.write("-- Replace foreign key constraint on duty location addresses\n")
f.write(
    "ALTER TABLE duty_locations ADD CONSTRAINT duty_locations_address_id_fkey FOREIGN KEY (address_id) REFERENCES addresses(id);\n"
)

f.close()
sys.exit()
