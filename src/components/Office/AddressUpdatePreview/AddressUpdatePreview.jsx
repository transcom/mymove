import React from 'react';
import classnames from 'classnames';
import { Alert } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './AddressUpdatePreview.module.scss';

import DataTable from 'components/DataTable/index';
import { formatTwoLineAddress } from 'utils/shipmentDisplay';
import DataTableWrapper from 'components/DataTableWrapper';
import { ShipmentAddressUpdateShape } from 'types';

const AddressUpdatePreview = ({ deliveryAddressUpdate, destSitServiceItems }) => {
  const { originalAddress, newAddress, contractorRemarks } = deliveryAddressUpdate;
  return (
    <div>
      <h3 className={styles.previewHeading}>Delivery location</h3>

      <Alert type="warning" className={styles.alert}>
        <span className={styles.alertContent}>
          If approved, the requested update to the delivery location will change one or all of the following:
          <span className={styles.listItem}>Service area.</span>
          <span className={styles.listItem}>Mileage bracket for direct delivery.</span>
          <span className={styles.listItem}>
            ZIP3 resulting in Domestic Shorthaul (DSH) changing to Domestic Linehaul (DLH) or vice versa.
          </span>
          Approvals will result in updated pricing for this shipment. Customer may be subject to excess costs.
        </span>
      </Alert>
      {destSitServiceItems.length > 0 ? (
        <Alert type="info" className={styles.alert} id="destSitAlert" data-testid="destSitAlert">
          <span className={styles.alertContent}>
            This shipment contains {destSitServiceItems.length} destination SIT service items. If approved, this could
            change the following:{' '}
            <span className={styles.listItem}>
              SIT delivery &gt; 50 miles <strong>or</strong> SIT delivery &le; 50 miles.
            </span>
            <span className={styles.listItem}>Service area.</span>
            <span className={styles.listItem}>Mileage bracket (for Direct Delivery).</span>
            <span className={styles.listItem}>Weight bracket change.</span>
            Approvals will result in updated pricing for the service item and require TOO approval. Customer may be
            subject to excess costs.
          </span>
        </Alert>
      ) : null}

      <DataTableWrapper
        className={classnames('maxw-tablet', 'table--data-point-group', styles.reviewAddressChange)}
        testID="address-change-preview"
      >
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

AddressUpdatePreview.propTypes = {
  deliveryAddressUpdate: ShipmentAddressUpdateShape.isRequired,
};
