import React from 'react';
import * as PropTypes from 'prop-types';
import { get } from 'lodash';

import styles from './OfficeDefinitionLists.module.scss';

import { BackupContactShape } from 'types/backupContact';
import descriptionListStyles from 'styles/descriptionList.module.scss';
import { ResidentialAddressShape } from 'types/address';

const CustomerInfoList = ({ customerInfo }) => {
  return (
    <div className={styles.OfficeDefinitionLists}>
      <dl className={descriptionListStyles.descriptionList}>
        <div className={descriptionListStyles.row}>
          <dt>Name</dt>
          <dd data-testid="name">{customerInfo.name}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>DoD ID</dt>
          <dd data-testid="dodId">{customerInfo.dodId}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Phone</dt>
          <dd data-testid="phone">{customerInfo.phone}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Email</dt>
          <dd data-testid="email">{customerInfo.email}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Current Address</dt>
          <dd data-testid="currentAddress">
            {`${get(customerInfo, 'currentAddress.street_address_1')}, ${get(
              customerInfo,
              'currentAddress.city',
            )}, ${get(customerInfo, 'currentAddress.state')} ${get(customerInfo, 'currentAddress.postal_code')}`}
          </dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Backup contact name</dt>
          <dd data-testid="backupContactName">{customerInfo.backupContact?.name}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Backup contact email</dt>
          <dd data-testid="backupContactEmail">{customerInfo.backupContact?.email}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Backup contact phone</dt>
          <dd data-testid="backupContactPhone">
            {customerInfo.backupContact?.phone ? `+1 ${customerInfo.backupContact.phone}` : ''}
          </dd>
        </div>
      </dl>
    </div>
  );
};

CustomerInfoList.propTypes = {
  customerInfo: PropTypes.shape({
    name: PropTypes.string,
    dodId: PropTypes.string,
    phone: PropTypes.string,
    email: PropTypes.string,
    currentAddress: ResidentialAddressShape,
    backupContact: BackupContactShape,
  }).isRequired,
};

export default CustomerInfoList;
