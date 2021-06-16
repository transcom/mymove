/* eslint-disable react/forbid-prop-types */
import React from 'react';
// import { useQueryWithStore } from 'react-admin';
import { ImportButton } from 'react-admin-import-csv';
import PropTypes from 'prop-types';

const avaliableRoles = [
  { roleType: 'customer', name: 'Customer' },
  { roleType: 'transportation_ordering_officer', name: 'Transportation Ordering Officer' },
  { roleType: 'transportation_invoicing_officer', name: 'Transportation Invoicing Officer' },
  { roleType: 'contracting_officer', name: 'Contracting Officer' },
  { roleType: 'ppm_office_users', name: 'PPM Office Users' },
  { roleType: 'services_counselor', name: 'Services Counselor' },
];

const ImportCsvButton = ({ props }) => {
  // call to get transportation office ids
  const validateRow = async (row) => {
    // Verify we have all required fields
    if (!(row.transportationOfficeId && row.firstName && row.lastName && row.roles && row.email && row.telephone)) {
      // throw error
    }

    // Verify the phone format
    const regex = /^[2-9]\d{2}-\d{3}-\d{4}$/;
    if (!regex.test(row.telephone)) {
      // throw error
    }

    // // transportation office id
    // const { offices } = useQueryWithStore({
    //   type: 'getOne',
    //   resource: 'users',
    //   payload: {},
    // });

    // const officeFound = offices.find((office) => {
    //   if (office.id === row.transportationOfficeId) return true;
    //   return false;
    // });

    // if (!officeFound) {
    //   // throw error
    // }

    return row;
  };

  const preCommitCallback = (action, rows) => {
    const alteredRows = [];
    rows.forEach((row) => {
      const copyOfRow = row;
      if (row.roles) {
        const rolesArray = [];

        // Parse roles from string
        const parsedRoles = row.roles.split(',');
        parsedRoles.forEach((parsedRole) => {
          // Remove any whitespace in the role string
          const role = parsedRole.replaceAll(/\s/g, '');
          rolesArray.push(avaliableRoles.find((avaliableRole) => avaliableRole.roleType === role));
        });

        // thow error if row is invalid?

        copyOfRow.roles = rolesArray;
      }
      alteredRows.push(copyOfRow);
    });
    return alteredRows;
  };

  const config = {
    logging: true,
    validateRow,
    preCommitCallback,
    postCommitCallback: () => {
      // console.log('reportItems', { reportItems });
    },
    disableImportOverwrite: true,
  };

  // eslint-disable-next-line react/jsx-props-no-spreading
  return <ImportButton {...props} {...config} />;
};

ImportCsvButton.propTypes = {
  props: PropTypes.object.isRequired,
  config: PropTypes.object.isRequired,
};

export default ImportCsvButton;
