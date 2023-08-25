import os, sys, pandas as pd
from datetime import datetime

ZERO = '0'

if len(sys.argv) < 2:
    sys.exit("Input file required.")


# Generate to gitignored tmp directory to prevent committing secure data
current_dir = os.getcwd()
destination_directory = 'tmp/generated-secure-migrations'
now = datetime.now()
year = str(now.year).rjust(4, ZERO)
month = str(now.month).rjust(2, ZERO)
day = str(now.day).rjust(2, ZERO)
hour = str(now.hour).rjust(2, ZERO)
minute = str(now.minute).rjust(2, ZERO)
second = str(now.second).rjust(2, ZERO)
filename = f'{year}{month}{day}{hour}{minute}{second}_update_loa_trnsn_ids.up.sql'
secure_migration_filename = (
    f'{current_dir}/{destination_directory}/{filename}'
)

destination_path = f'{current_dir}/{destination_directory}'
if not os.path.exists(destination_path):
    os.makedirs(destination_path)

with open(secure_migration_filename, "w+") as f:
    f.write('-- Update loa_trnsn_id column constraint\n')
    f.write('ALTER TABLE lines_of_accounting ALTER COLUMN loa_trnsn_id TYPE varchar (3);\n')

    # Skip the first and last rows which are just "unclassified"
    input_file = pd.read_excel(sys.argv[1], skiprows=1, skipfooter=1)

    # Missing values should be NULL
    input_file = input_file.fillna('NULL')

    f.write('-- Update lines_of_accounting with updated loa_trnsn_id values, mapped by loa_sys_id\n')
    f.write('UPDATE lines_of_accounting AS loas SET\n')
    f.write('\tloa_trnsn_id = updated.loa_trnsn_id\n')
    f.write('FROM (VALUES\n')

    for index, row in input_file.iterrows():
        loa_sys_id = row['LOA_SYS_ID']

        # Ignore rows where loa_sys_id is missing
        if loa_sys_id == 'NULL':
            continue

        # Add single quotes around non-null values, otherwise, just use NULL
        loa_trnsn_id = row['LOA_TRNSN_ID']
        loa_trnsn_id_write_value = loa_trnsn_id if loa_trnsn_id == 'NULL' else f"'{loa_trnsn_id}'"

        f.write(f'\t({loa_sys_id}, {loa_trnsn_id_write_value})')

        # only add a comma to rows that are not the last row
        if index == len(input_file) - 1:
            f.write('\n')
        else:
            f.write(',\n')

    f.write(') AS updated(loa_sys_id, loa_trnsn_id)\n')
    f.write('WHERE updated.loa_sys_id = loas.loa_sys_id;\n')
sys.exit()
