import React from 'react';
import classnames from 'classnames';
import { Alert } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './AddressUpdatePreview.module.scss';

import DataTable from 'components/DataTable/index';
import { formatTwoLineAddress } from 'utils/shipmentDisplay';
import DataTableWrapper from 'components/DataTableWrapper';

const AddressUpdatePreview = ({ deliveryAddressUpdate }) => {
  const { originalAddress, newAddress, contractorRemarks } = deliveryAddressUpdate;
  return (
    <div>
      <h3 className={styles.previewHeading}>Delivery location</h3>

      <Alert type="warning">
        If approved, the requested update to the delivery location will change one or all of the following:
        <span className={styles.listItem}>Service area.</span>
        <span className={styles.listItem}>Mileage bracket for direct delivery.</span>
        <span className={styles.listItem}>
          ZIP3 resulting in Domestic Shorthaul (DSH) changing to Domestic Linehaul (DLH) or vice versa.
        </span>
        Approvals will result in updated pricing for this shipment. Customer may be subject to excess costs.
      </Alert>

      <DataTableWrapper className={classnames('maxw-tablet', 'table--data-point-group', styles.reviewAddressChange)}>
        <DataTable
          columnHeaders={['Original delivery location', 'Requested delivery location']}
          dataRow={[formatTwoLineAddress(originalAddress), formatTwoLineAddress(newAddress)]}
          icon={<FontAwesomeIcon icon="arrow-right" />}
        />
        <DataTable
          columnHeaders={['Update request details']}
          custClass={styles.contractorRemarks}
          dataRow={[
            <>
              <b>Contractor remarks:</b> {contractorRemarks}
            </>,
          ]}
        />
      </DataTableWrapper>
    </div>
  );
};

export default AddressUpdatePreview;
