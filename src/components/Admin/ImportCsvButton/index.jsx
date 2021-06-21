import React from 'react';
import { useNotify } from 'react-admin';
import { ImportButton } from 'react-admin-import-csv';
import PropTypes from 'prop-types';

import { adminOfficeRoles } from 'constants/userRoles';

// Note: There is not a test file or story for ImportCsvButton beacuse this component HAS to render within a react-admin app
const ImportCsvButton = (props) => {
  const notify = useNotify();
  const { resource } = props;

  const validateRow = async (row) => {
    // Verify we have all required fields
    if (!(row.transportationOfficeId && row.firstName && row.lastName && row.roles && row.email && row.telephone)) {
      const err = new Error(
        `Validation Error: Row does not contain all required fields.
        Required fields are firstName, lastName, email, telephone, roles, and transportationOfficeId \n
        Row Information: ${Object.values(row)}`,
      );
      notify(err.message);
      throw err;
    }

    // Verify the phone format
    const regex = /^[2-9]\d{2}-\d{3}-\d{4}$/;
    if (!regex.test(row.telephone)) {
      const err = new Error(
        `Validation Error: Row contains improperly formatted telephone number. Required format is xxx-xxx-xxxx. \n
        Row Information: ${Object.values(row)}`,
      );
      notify(err.message);
      throw err;
    }

    return row;
  };

  const preCommitCallback = (action, rows) => {
    const alteredRows = [];
    rows.forEach((row) => {
      const copyOfRow = row;
      if (row.roles) {
        const rolesArray = [];

        // Parse roles from string at ","
        const parsedRoles = row.roles.split(',');
        parsedRoles.forEach((parsedRole) => {
          let roleNotFound = true;
          // Remove any whitespace in the role string
          const role = parsedRole.replaceAll(/\s/g, '');
          adminOfficeRoles.forEach((adminOfficeRole) => {
            if (adminOfficeRole.roleType === role) {
              rolesArray.push(adminOfficeRole);
              roleNotFound = false;
            }
          });
          if (roleNotFound) {
            const err = new Error(
              `Processing Error: Invalid roles provided for row. \n Row Information: ${Object.values(row)}`,
            );
            notify(err.message);
            throw err;
          }
        });
        copyOfRow.roles = rolesArray;
      } else {
        const err = new Error(
          `Processing Error: Unable to parse roles for row. \n Row Information: ${Object.values(row)}`,
        );
        notify(err.message);
        throw err;
      }
      alteredRows.push(copyOfRow);
    });
    return alteredRows;
  };

  const postCommitCallback = (reportItems) => {
    reportItems.forEach((reportItem) => {
      if (reportItem.err) {
        return notify(
          `${reportItem.err.name} ${reportItem.err.status}: ${reportItem.err.message}.  \n ${reportItem.err.body.detail}`,
        );
      }
      return null;
    });
  };

  const config = {
    logging: false,
    validateRow,
    preCommitCallback,
    postCommitCallback,
    disableImportOverwrite: true,
  };

  return (
    // eslint-disable-next-line react/jsx-props-no-spreading
    <ImportButton resource={resource} {...props} {...config} />
  );
};

ImportCsvButton.propTypes = {
  resource: PropTypes.string.isRequired,
};

export default ImportCsvButton;
