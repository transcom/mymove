import React from 'react';
import * as PropTypes from 'prop-types';
import { get } from 'lodash';

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
            <th scope="row" className="text-bold">
              Name
            </th>
            <td data-cy="name">{customerInfo.name}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              DoD ID
            </th>
            <td data-cy="dodId">{customerInfo.dodId}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Phone
            </th>
            <td data-cy="phone">{customerInfo.phone}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Email
            </th>
            <td data-cy="email">{customerInfo.email}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Current address
            </th>
            <td data-cy="currentAddress">
              {`${get(customerInfo, 'currentAddress.street_address_1')}, ${get(
                customerInfo,
                'currentAddress.city',
              )}, ${get(customerInfo, 'currentAddress.state')} ${get(customerInfo, 'currentAddress.postal_code')}`}
            </td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Backup contact name
            </th>
            <td data-cy="backupContactName">{customerInfo.backupContactName}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Backup contact phone
            </th>
            <td data-cy="backupContactPhone">{customerInfo.backupContactPhone}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Backup contact email
            </th>
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
    backupContactName: PropTypes.string,
    backupContactPhone: PropTypes.string,
    backupContactEmail: PropTypes.string,
  }).isRequired,
};

export default CustomerInfoTable;
