"""
Jira ticket: MB-16327
This script takes as input a CSV of Transportation Offices, and outputs a SQL migration file that updates the `addresses` and `transportation_offices` tables according to the provided Excel spreadsheet.
If the row in the spreadsheet contains the TO id, then update the corresponding row in the database with the information in the spreadsheet.
If the row in the spreadsheet doesn't contain the TO id, then create a new TO and its corresponding address.

See the Duty Station Locations Discovery doc for more context (https://dp3.atlassian.net/wiki/spaces/MT/pages/2206662681/Discovery+Duty+Station+Locations)
 """

import os, sys, uuid, pandas as pd
from datetime import datetime



"""
Process transportation office data from the input csv file and save results to the output sql file.

Args:
    input_file (str): Path to the input csv file.
    output_file (str): Path to the output sql file.
"""
def process_transportation_office_csv(input_file: str, output_file: str) -> None:
  try:
    df = pd.read_excel(pd.ExcelFile(input_file))
    df = df.reset_index()

    with open(output_file, 'w') as out:
      # Wrap this in a transaction so if anything goes wrong, the entire operation gets rolled back
      out.write("BEGIN;\n\n")

      for _, row in df.iterrows():
        transportation_office_id = row["id"]
        name = row["name"]
        gbloc = row["GBLOC"]

        street_address_1 = row["street_address_1"]
        street_address_2 = row["street_address_2"]
        city = row["city"]
        state = row["state"]
        # trim to zip5
        postal_code = "%05d" % int(str(row["postal_code"]).split("-")[0])

        # If transportation_office_id is missing, create a new TO and its corresponding address
        if pd.isna(transportation_office_id):
            transportation_office_id = uuid.uuid4()
            address_id = uuid.uuid4()
            out.write("-- Insert the address and TO\n")
            out.write(
                f"INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at) VALUES ('{address_id}', '{street_address_1}', '{street_address_2}', '{city}', '{state}', '{postal_code}', now(), now());\n\n"
            )
            out.write(
                f"INSERT INTO transportation_offices (id, address_id, gbloc, name, created_at, updated_at, latitude, longitude) VALUES ('{transportation_office_id}', '{address_id}', '{gbloc}', '{name}', now(), now(), 0, 0);\n\n"
            )
        else:
            # If transportation_office_id is present, update the TO and its corresponding address
            out.write("-- Update the TO and address\n")
            out.write(
                f"UPDATE transportation_offices SET name = '{name}', gbloc = '{gbloc}' WHERE id = '{transportation_office_id}';\n\n"
            )
            out.write(
                f"UPDATE addresses SET street_address_1 = '{street_address_1}', street_address_2 = '{street_address_2}', city = '{city}', state = '{state}', postal_code = '{postal_code}' WHERE id = (SELECT address_id FROM transportation_offices where id = '{transportation_office_id}');\n\n"
            )

      out.write("END;\n")

  except FileNotFoundError:
      sys.stderr.write(f"File not found: {input_file}")

if __name__ == "__main__":
    if len(sys.argv) < 2:
      sys.exit("Input file required.")

    input_file = sys.argv[1]
    current_dir = os.getcwd()
    current_datetime_formatted = datetime.now().strftime("%Y%m%d%H%M%S")
    output_file = f"{current_dir}/migrations/app/schema/{current_datetime_formatted}_update_transportation_offices.up.sql"

    process_transportation_office_csv(input_file, output_file)
