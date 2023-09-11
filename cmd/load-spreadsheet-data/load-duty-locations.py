import os, sys, uuid
from duty_location_data import *

current_dir = os.getcwd()

filename = (
    f"{current_dir}/migrations/app/schema/20230810180036_update_duty_locations.up.sql"
)

f = open(filename, "w")


# takes a comma-separated string of aliases and an id
# writes SQL queries that remove duty_location_names that match each alias
# and that create new duty_location_names for the duty_location id given
def handle_aliases(aliases, id):
    # remove aliases from existing
    if id is not None:
        f.write("-- Remove existing associated duty_location_names\n")
        aliases = aliases.split(r"\w?,\w?")
        for alias in aliases:
            dln_id = uuid.uuid4()
            f.write("-- Delete old duty_location_names\n")
            f.write(f"DELETE FROM duty_location_names where name = '{alias}';\n")
            f.write("-- Insert new duty_location_names\n")
            f.write(
                f"INSERT INTO duty_location_names (id, name, duty_location_id, created_at, updated_at) VALUES('{dln_id}', '{alias}', '{id}', now(), now());\n"
            )


# takes a duty_location id and deletes all of its recursive parent records with an FK constraint
# this is done considering that any existing records in these tables are safe to delete, and is not necessarily a model for future duty location updates
def delete_dl_and_parents(id):
    for t in [
        "archived_move_documents",
        "archived_signed_certifications",
        "archived_personally_procured_moves",
        "customer_support_remarks",
        "evaluation_reports",
        "mto_service_items",
        "mto_shipments",
        "payment_requests",
        "signed_certifications",
        "personally_procured_moves",
        "webhook_notifications",
    ]:
        f.write(
            f"DELETE from {t} where move_id = (SELECT id from moves where orders_id = (SELECT id from orders where origin_duty_location_id = '{dl_id}' or new_duty_location_id = '{dl_id}'));\n"
        )
    f.write(
        f"DELETE from moves where orders_id = (SELECT id from orders where origin_duty_location_id = '{dl_id}' or new_duty_location_id = '{dl_id}');\n"
    )
    f.write(
        f"DELETE from orders where origin_duty_location_id = '{dl_id}' or new_duty_location_id = '{dl_id}';\n"
    )
    for t in ["documents", "archived_access_codes", "notifications"]:
        f.write(
            f"DELETE from {t} where service_member_id = (SELECT id from service_members where duty_location_id = '{dl_id}');\n"
        )
    f.write(f"DELETE from service_members where duty_location_id = '{dl_id}';\n")
    f.write(f"DELETE from duty_location_names where duty_location_id = '{dl_id}';\n")
    f.write(f"DELETE from duty_locations where id = '{dl_id}';\n\n")


for delete in deletes:
    # e.g. "Adak, AK 99546",
    f.write(f"--DELETE\n")
    dl_id = delete[1]
    delete_dl_and_parents(dl_id)


for rename in renames:
    # Fort Gordon => Fort Eisenhower
    f.write(
        f"UPDATE duty_locations SET name = '{rename[1]}' WHERE name = '{rename[0]}';\n\n"
    )

for new in news:
    (
        id,
        name,
        address_id,
        gbloc,
        affiliation,
        street_address,
        city,
        state,
        postal_code,
        aliases,
    ) = new

    affiliation = f"'{affiliation}'" if affiliation else "NULL"
    gbloc = gbloc if gbloc else "NULL"

    f.write("--NEW\n")
    address_query = f"(SELECT id from addresses where street_address_1 = '{street_address}' and city = '{city}' and state = '{state}' and postal_code = '{postal_code}' LIMIT 1)"
    to_query = f"(SELECT t.id FROM transportation_offices AS t, addresses AS a WHERE t.address_id = a.id AND t.gbloc='{gbloc}' AND a.state='{state}' LIMIT 1)"

    f.write(
        f"INSERT INTO addresses (id, street_address_1, city, state, postal_code, created_at, updated_at) VALUES ('{address_id}', '{street_address}', '{city}', '{state}', '{postal_code}', now(), now());\n"
    )

    f.write(
        "INSERT INTO duty_locations (id, address_id, transportation_office_id, name, affiliation, provides_services_counseling, updated_at, created_at) "
    )
    f.write(
        f"VALUES ('{id}', '{address_id}', {to_query}, '{name}', {affiliation}, TRUE, now(), now());\n\n"
    )

for merge in merges:
    # e.g. Aberdeen Proving Ground
    old = merge[0]
    new = merge[1]
    new_dl_id = merge[2]
    old_dl_id = f"(SELECT id from duty_locations where name = '{old}')"
    f.write("-- MERGE\n")
    f.write(
        f"UPDATE orders SET origin_duty_location_id='{new_dl_id}' WHERE origin_duty_location_id={old_dl_id};\n"
    )
    f.write(
        f"UPDATE orders SET new_duty_location_id='{new_dl_id}' WHERE new_duty_location_id={old_dl_id};\n"
    )
    f.write(
        f"UPDATE service_members SET duty_location_id='{new_dl_id}' WHERE duty_location_id={old_dl_id};\n"
    )
    f.write(
        f"UPDATE duty_location_names SET duty_location_id = '{new_dl_id}' WHERE duty_location_id={old_dl_id};\n"
    )
    f.write(
        f"UPDATE duty_locations SET transportation_office_id=to_id, affiliation=aff FROM (SELECT transportation_office_id as to_id, affiliation as aff FROM duty_locations where id={old_dl_id}) dl WHERE id='{new_dl_id}';\n"
    )
    delete_dl_and_parents(old_dl_id)
    handle_aliases(old, new_dl_id)

f.close()
sys.exit()
