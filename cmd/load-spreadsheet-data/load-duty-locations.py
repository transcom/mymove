import os, sys, uuid, pandas as pd

current_dir = os.getcwd()

filename = (
    f"{current_dir}/migrations/app/schema/20230810180036_update_duty_locations.up.sql"
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
f.write("-- Remove unique index from duty_location_names\n")
f.write("DROP INDEX IF EXISTS duty_location_names_name_idx;\n\n")
f.write("CREATE INDEX duty_location_names_name_idx ON duty_location_names(name);")

df = pd.read_excel(
    pd.ExcelFile(sys.argv[1]), keep_default_na=False, na_values=["", "nan"]
)
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
    provides_services_counseling = "TRUE"
    aliases = row["Alias"]

    address_id = uuid.uuid4()
    id = uuid.uuid4()

    if pd.isna(street_address):
        street_address = ""

    f.write("-- Insert the address\n")
    f.write(
        f"INSERT INTO addresses (id, street_address_1, city, state, postal_code, created_at, updated_at) VALUES ('{address_id}', '{street_address}', '{city}', '{state}', '{postal_code}', now(), now());\n\n"
    )

    if not pd.isna(gbloc):
        f.write("-- Find the transportation office by GBLOC, City, and State\n")
        f.write(
            f"WITH found_to AS (SELECT t.id FROM transportation_offices AS t, addresses AS a WHERE t.address_id = a.id AND t.gbloc='{gbloc}' AND a.city='{city}' AND a.state='{state}' LIMIT 1)\n\n"
        )
        to_id = "(select id from found_to)"

    if pd.isna(affiliation):
        affiliation = "NULL"
    else:
        affiliation = "_".join([w.upper() for w in affiliation.split()])
        affiliation = f"'{affiliation}'"

    f.write("-- Insert the duty location\n")
    f.write(
        f"INSERT INTO duty_locations (id, address_id, transportation_office_id, name, affiliation, provides_services_counseling, updated_at, created_at) VALUES('{id}', '{address_id}', {to_id}, '{name}', {affiliation}, {provides_services_counseling}, now(), now());\n\n"
    )
    if not pd.isna(aliases):
        aliases = aliases.split(r"\w?,\w?")
        for alias in aliases:
            dln_id = uuid.uuid4()
            f.write("-- Insert search alias as duty_location_name\n")
            f.write(
                f"INSERT INTO duty_location_names (id, name, duty_location_id, created_at, updated_at) VALUES('{dln_id}', '{alias}', '{id}', now(), now());"
            )

f.write("-- Replace foreign key constraint on duty location addresses\n")
f.write(
    "ALTER TABLE duty_locations ADD CONSTRAINT duty_locations_address_id_fkey FOREIGN KEY (address_id) REFERENCES addresses(id);\n"
)

f.close()
sys.exit()
