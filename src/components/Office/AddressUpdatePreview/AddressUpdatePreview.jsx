import React from 'react';
import classnames from 'classnames';
import { Alert } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './AddressUpdatePreview.module.scss';

import DataTable from 'components/DataTable/index';
import { formatTwoLineAddress } from 'utils/shipmentDisplay';
import DataTableWrapper from 'components/DataTableWrapper';
import { ShipmentAddressUpdateShape } from 'types';
import { MARKET_CODES } from 'shared/constants';

const AddressUpdatePreview = ({ deliveryAddressUpdate, shipment }) => {
  const { originalAddress, newAddress, contractorRemarks } = deliveryAddressUpdate;
  const newSitMileage = deliveryAddressUpdate.newSitDistanceBetween;
  const { marketCode } = shipment;
  return (
    <div>
      <h3 className={styles.previewHeading}>Delivery Address</h3>
      <Alert type="warning" className={styles.alert}>
        {marketCode === MARKET_CODES.DOMESTIC ? (
          <span className={styles.alertContent}>
            If approved, the requested update to the delivery address will change one or all of the following:
            <span className={styles.listItem}>Service area.</span>
            <span className={styles.listItem}>Mileage bracket for direct delivery.</span>
            <span className={styles.listItem}>
              ZIP3 resulting in Domestic Shorthaul (DSH) changing to Domestic Linehaul (DLH) or vice versa.
            </span>
            Approvals will result in updated pricing for this shipment. Customer may be subject to excess costs.
          </span>
        ) : (
          <span className={styles.alertContent}>
            If approved, the requested update to the delivery address will change one or all of the following:
            <span className={styles.listItem}>The rate area for the international shipment destination address.</span>
            <span className={styles.listItem}>Pricing for international shipping & linehaul.</span>
            <span className={styles.listItem}>Pricing for POD Fuel Surcharge (if applicable).</span>
            Approvals will result in updated pricing for this shipment. Customer may be subject to excess costs.
          </span>
        )}
      </Alert>
      {newSitMileage > 50 ? (
        <Alert type="warning" className={styles.alert} id="destSitAlert" data-testid="destSitAlert">
          <span className={styles.alertContent}>
            Approval of this address change request will result in SIT Delivery &gt; 50 Miles.
            <br />
            Updated Mileage for SIT: <strong>{newSitMileage} miles</strong>
          </span>
        </Alert>
      ) : null}
      <DataTableWrapper
        className={classnames('maxw-tablet', 'table--data-point-group', styles.reviewAddressChange)}
        testID="address-change-preview"
      >
        <DataTable
          columnHeaders={['Original Delivery Address', 'Requested Delivery Address']}
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
