import React from 'react';
import classnames from 'classnames';
import { Alert } from '@trussworks/react-uswds';
// import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import { formatAddress } from '../../../utils/shipmentDisplay';
import DataTableWrapper from '../../DataTableWrapper/index';
import DataTable from '../../DataTable/index';

import styles from './ServiceItemUpdateModal.module.scss';

import AddressFields from 'components/form/AddressFields/AddressFields';

const EditSitAddressChangeForm = ({ initialAddress }) => {
  return (
    <div className={styles.editSitAddressChangeForm}>
      <DataTableWrapper className={classnames('maxw-tablet', 'table--data-point-group', styles.initialAddress)}>
        <DataTable columnHeaders={['Initial SIT delivery address']} dataRow={[formatAddress(initialAddress)]} />
      </DataTableWrapper>
      <h3 className={styles.destinationAddressTitle}>Final SIT delivery</h3>
      <Alert type="warning">Approvals over 50 miles will result in updated pricing for this shipment.</Alert>
      <div data-testid="editAddressForm">
        <AddressFields name="newAddress" />
      </div>
    </div>
  );
};

export default EditSitAddressChangeForm;
