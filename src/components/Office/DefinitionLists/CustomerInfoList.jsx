import React from 'react';
import * as PropTypes from 'prop-types';

import styles from './OfficeDefinitionLists.module.scss';

import { BackupContactShape } from 'types/backupContact';
import descriptionListStyles from 'styles/descriptionList.module.scss';
import { AddressShape } from 'types/address';
import { formatCustomerContactFullAddress } from 'utils/formatters';
import departmentIndicators from 'constants/departmentIndicators';

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
          <dd data-testid="edipi">{customerInfo.edipi}</dd>
        </div>
        {customerInfo.agency === departmentIndicators.COAST_GUARD && (
          <div className={descriptionListStyles.row}>
            <dt>EMPLID</dt>
            <dd data-testid="emplid">{customerInfo.emplid}</dd>
          </div>
        )}
        <div className={descriptionListStyles.row}>
          <dt>Phone</dt>
          <dd data-testid="phone">{`+1 ${customerInfo.phone}`}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Alt. Phone</dt>
          <dd data-testid="phone">{customerInfo.altPhone ? `+1 ${customerInfo.altPhone}` : '—'}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Email</dt>
          <dd data-testid="email">{customerInfo.email}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Pickup Address</dt>
          <dd data-testid="currentAddress">
            {customerInfo.currentAddress?.streetAddress1
              ? formatCustomerContactFullAddress(customerInfo.currentAddress)
              : '—'}
          </dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Backup address</dt>
          <dd data-testid="backupAddress">
            {customerInfo.backupAddress?.streetAddress1
              ? formatCustomerContactFullAddress(customerInfo.backupAddress)
              : '—'}
          </dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Backup contact name</dt>
          <dd data-testid="backupContactName">
            {customerInfo.backupContact?.name ? customerInfo.backupContact.name : '—'}
          </dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Backup contact email</dt>
          <dd data-testid="backupContactEmail">
            {customerInfo.backupContact?.email ? customerInfo.backupContact.email : '—'}
          </dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Backup contact phone</dt>
          <dd data-testid="backupContactPhone">
            {customerInfo.backupContact?.phone ? `+1 ${customerInfo.backupContact.phone}` : '—'}
          </dd>
        </div>
      </dl>
    </div>
  );
};

CustomerInfoList.propTypes = {
  customerInfo: PropTypes.shape({
    name: PropTypes.string,
    edipi: PropTypes.string,
    phone: PropTypes.string,
    email: PropTypes.string,
    currentAddress: AddressShape,
    backupAddress: AddressShape,
    backupContact: BackupContactShape,
  }).isRequired,
};

export default CustomerInfoList;
