import React from 'react';
import * as PropTypes from 'prop-types';
import { get } from 'lodash';

import styles from './OrdersTable/OrdersTable.module.scss';

import { BackupContactShape } from 'types/backupContact';

const CustomerInfoTable = ({ customerInfo }) => {
  return (
    <div className={styles.OrdersTable}>
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
            <td data-testid="name">{customerInfo.name}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              DoD ID
            </th>
            <td data-testid="dodId">{customerInfo.dodId}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Phone
            </th>
            <td data-testid="phone">{customerInfo.phone}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Email
            </th>
            <td data-testid="email">{customerInfo.email}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Current address
            </th>
            <td data-testid="currentAddress">
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
            <td data-testid="backupContactName">{customerInfo.backupContact?.name}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Backup contact email
            </th>
            <td data-testid="backupContactEmail">{customerInfo.backupContact?.email}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Backup contact phone
            </th>
            <td data-testid="backupContactPhone">
              {customerInfo.backupContact?.phone ? `+1 ${customerInfo.backupContact.phone}` : ''}
            </td>
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
    backupContact: BackupContactShape,
  }).isRequired,
};

export default CustomerInfoTable;
