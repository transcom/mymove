import React from 'react';
import { useNotify } from 'react-admin';
import { ImportButton } from 'react-admin-import-csv';
import PropTypes from 'prop-types';

import {
  checkRequiredFields,
  checkTelephone,
  checkValidRolesWithPrivileges,
  parseRoles,
  parsePrivileges,
} from './validation';

// Note: There is not a test file or story for ImportOfficeUserButton beacuse this component HAS to render within a react-admin app
const ImportOfficeUserButton = (props) => {
  const notify = useNotify();
  const { resource } = props;

  const validateRow = async (row) => {
    // Verify we have all required fields and that the telephone is valid
    const validation = [checkRequiredFields, checkTelephone];
    validation.forEach((check) => {
      try {
        // eslint-disable-next-line react/jsx-props-no-spreading
        check({ ...row });
      } catch (err) {
        notify(`${err.message} \n Row Information: ${Object.values(row)}`);
        throw err;
      }
    });

    return row;
  };

  const preCommitCallback = (action, rows) => {
    const alteredRows = [];
    rows.forEach((row) => {
      const copyOfRow = row;
      try {
        if (checkValidRolesWithPrivileges(row)) {
          const parsedRolesArray = parseRoles(row.roles);
          copyOfRow.roles = parsedRolesArray;

          const parsedPrivilegesArray = row.privileges ? parsePrivileges(row.privileges) : null;
          copyOfRow.privileges = parsedPrivilegesArray;
        }
      } catch (err) {
        notify(`${err.message} \n Row Information: ${Object.values(row)}`);
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

ImportOfficeUserButton.propTypes = {
  resource: PropTypes.string.isRequired,
};

export default ImportOfficeUserButton;
