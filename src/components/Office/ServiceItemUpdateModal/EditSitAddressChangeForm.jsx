import React from 'react';
import classnames from 'classnames';
import { Alert } from '@trussworks/react-uswds';

import { formatAddressForSitAddressChangeForm } from '../../../utils/shipmentDisplay';
import DataTableWrapper from '../../DataTableWrapper/index';
import DataTable from '../../DataTable/index';

import styles from './ServiceItemUpdateModal.module.scss';

import { AddressShape } from 'types/address';
import AddressFields from 'components/form/AddressFields/AddressFields';

/**
 * @component
 * @description This is the form specific to for when a TOO edits a SIT destination address. It inluded the Initial address box, a distance alert, and a form to edit the address.
 * @param {EditSitAddressFormProps}
 * @returns {React.ReactElement}
 */
const EditSitAddressChangeForm = ({ initialAddress }) => {
  return (
    <div className={styles.editSitAddressChangeForm}>
      <DataTableWrapper className={classnames('maxw-tablet', 'table--data-point-group', styles.initialAddress)}>
        <DataTable
          columnHeaders={['Initial SIT delivery address']}
          dataRow={[formatAddressForSitAddressChangeForm(initialAddress)]}
        />
      </DataTableWrapper>
      <h3 className={styles.destinationAddressTitle}>Final SIT delivery</h3>
      <Alert type="warning">Approvals over 50 miles will result in updated pricing for this shipment.</Alert>
      <div data-testid="editAddressForm">
        <AddressFields name="newAddress" />
      </div>
    </div>
  );
};

/**
 * @typedef EditSitAddressFormProps
 * @prop{AddressShape} initialAddress
 */
EditSitAddressChangeForm.propTypes = {
  initialAddress: AddressShape.isRequired,
};
export default EditSitAddressChangeForm;
