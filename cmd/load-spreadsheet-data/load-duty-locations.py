import os, sys, uuid, pandas as pd

current_dir = os.getcwd()

filename = (
    f"{current_dir}/migrations/app/schema/20230810180036_update_duty_locations.up.sql"
)

f = open(filename, "w")

if len(sys.argv) < 2:
    sys.exit("Input file required.")


# takes a comma-separated string of aliases and an id
# if the id is not null, it removes existing duty_location_names and adds new ones
# otherwise, it adds them
def handle_aliases(aliases, id):
    # remove aliases from existing
    if id is not None:
        f.write("-- Remove existing associated duty_location_names\n")
        f.write(f"DELETE FROM duty_location_names WHERE duty_location_id = '{id}'")
    if not pd.isna(aliases):
        aliases = aliases.split(r"\w?,\w?")
        for alias in aliases:
            dln_id = uuid.uuid4()
            f.write("-- Insert new duty_location_names\n")
            f.write(
                f"INSERT INTO duty_location_names (id, name, duty_location_id, created_at, updated_at) VALUES('{dln_id}', '{alias}', '{id}', now(), now());"
            )


def handle_to(gbloc):
    if not pd.isna(gbloc):
        f.write(
            f"WITH found_to AS (SELECT t.id FROM transportation_offices AS t, addresses AS a WHERE t.address_id = a.id AND t.gbloc='{gbloc}' AND a.city='{city}' AND a.state='{state}' LIMIT 1)\n"
        )


f.write("-- Remove unique index from duty_location_names\n")
f.write("DROP INDEX IF EXISTS duty_location_names_name_idx;\n\n")

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

    precomma_name = name.split(", ")[0]

    if pd.isna(street_address):
        street_address = ""

    if pd.isna(affiliation):
        affiliation = "NULL"
    else:
        affiliation = "_".join([w.upper() for w in affiliation.split()])
        affiliation = f"'{affiliation}'"

    dl_query = f"(SELECT id from duty_locations where name = '{precomma_name}' UNION SELECT dl.id from duty_locations as dl, addresses as a where a.city = '{city}' and a.state = '{state}' and a.postal_code = '{postal_code}' and dl.address_id = a.id LIMIT 1)"

    f.write(f"UPDATE duty_locations SET name = '{name}' WHERE id = {dl_query}\n")
    # TODO: get the values for the new DL
    # TODO: handle address
    handle_to(gbloc)
    f.write(
        "INSERT INTO duty_locations (id, address_id, transportation_office_id, name, affiliation, provides_services_counseling, updated_at, created_at) ON CONFLICT DO NOTHING\n"
    )
    f.write(
        f"VALUES ('{id}', (SELECT id FROM addresses LIMIT 1), '{to_id}', '{name}', {affiliation}, TRUE, now(), now());\n\n"
    )

f.close()
sys.exit()
