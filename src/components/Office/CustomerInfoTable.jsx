import React from 'react';
import * as PropTypes from 'prop-types';

const CustomerInfoTable = ({ customerInfo }) => {
  return (
    <div>
      <div className="stackedtable-header">
        <div>
          <h4>Customer info</h4>
        </div>
      </div>
      <table className="table--stacked">
        <colgroup>
          <col style={{ width: '25%' }} />
          <col style={{ width: '75%' }} />
        </colgroup>
        <tbody>
          <tr>
            <th scope="row">Name</th>
            <td data-cy="name">{customerInfo.name}</td>
          </tr>
          <tr>
            <th scope="row">DoD ID</th>
            <td data-cy="dodId">{customerInfo.dodId}</td>
          </tr>
          <tr>
            <th scope="row">Phone</th>
            <td data-cy="phone">{customerInfo.phone}</td>
          </tr>
          <tr>
            <th scope="row">Email</th>
            <td data-cy="email">{customerInfo.email}</td>
          </tr>
          <tr>
            <th scope="row">Current address</th>
            <td data-cy="currentAddress">{`${customerInfo.currentAddress.street_address_1}, ${customerInfo.currentAddress.city}, ${customerInfo.currentAddress.state} ${customerInfo.currentAddress.postal_code}`}</td>
          </tr>
          <tr>
            <th scope="row">Destination address</th>
            <td data-cy="destinationAddress">{`${customerInfo.destinationAddress.street_address_1}, ${customerInfo.destinationAddress.city}, ${customerInfo.destinationAddress.state} ${customerInfo.destinationAddress.postal_code}`}</td>
          </tr>
          <tr>
            <th scope="row">Backup contact name</th>
            <td data-cy="backupContactName">{customerInfo.backupContactName}</td>
          </tr>
          <tr>
            <th scope="row">Backup contact phone</th>
            <td data-cy="backupContactPhone">{customerInfo.backupContactPhone}</td>
          </tr>
          <tr>
            <th scope="row">Backup contact email</th>
            <td data-cy="backupContactEmail">{customerInfo.backupContactEmail}</td>
          </tr>
        </tbody>
      </table>
    </div>
  );
};

CustomerInfoTable.propTypes = {
  customerInfo: PropTypes.shape({
    name: PropTypes.string,
    dodId: PropTypes.string,
    phone: PropTypes.string,
    email: PropTypes.string,
    currentAddress: PropTypes.shape({
      street_address_1: PropTypes.string,
      city: PropTypes.string,
      state: PropTypes.string,
      postal_code: PropTypes.string,
    }),
    destinationAddress: PropTypes.shape({
      street_address_1: PropTypes.string,
      city: PropTypes.string,
      state: PropTypes.string,
      postal_code: PropTypes.string,
    }),
    backupContactName: PropTypes.string,
    backupContactPhone: PropTypes.string,
    backupContactEmail: PropTypes.string,
  }).isRequired,
};

export default CustomerInfoTable;
