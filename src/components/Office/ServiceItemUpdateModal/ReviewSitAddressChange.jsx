import React from 'react';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Alert } from '@trussworks/react-uswds';
import * as PropTypes from 'prop-types';

import styles from './ServiceItemUpdateModal.module.scss';

import DataTable from 'components/DataTable/index';
import { formatAddress } from 'utils/shipmentDisplay';
import { AddressShape } from 'types';
import DataTableWrapper from 'components/DataTableWrapper';

/**
 * @component
 * @description This modal serves the purpose of allowing a TOO to review a SIT address update in order to approve or deny the change.
 * @param {ReviewSitAddressChangeProps}
 * @returns {React.ReactElement}
 */
const ReviewSitAddressChange = ({ sitAddressUpdate }) => {
  const { oldAddress, newAddress, contractorRemarks, distance } = sitAddressUpdate;

  return (
    <>
      <Alert type="warning" style={{ marginBottom: '1rem' }} data-testid="distanceAlert">
        Requested final SIT delivery address is {distance} miles from the initial SIT delivery address. Approvals over
        50 miles will result in updated pricing for this shipment.
      </Alert>
      <DataTableWrapper className={classnames('maxw-tablet', 'table--data-point-group', styles.reviewAddressChange)}>
        <DataTable
          columnHeaders={['Initial SIT delivery address', 'Requested final SIT delivery address']}
          dataRow={[formatAddress(oldAddress), formatAddress(newAddress)]}
          icon={<FontAwesomeIcon icon="arrow-right" />}
        />
        <DataTable
          columnHeaders={['Update request details']}
          dataRow={[
            <>
              <b>Contractor remarks:</b> {contractorRemarks}
            </>,
          ]}
        />
      </DataTableWrapper>
    </>
  );
};

/**
 * @typedef ReviewSitAddressChangeProps
 * @prop{AddressShape} oldAddress
 * @prop{AddressShape} newAddress
 * @prop{string} contractorRemarks
 * @prop{string} distance
 */
ReviewSitAddressChange.propTypes = {
  oldAddress: AddressShape,
  newAddress: AddressShape,
  contractorRemarks: PropTypes.string,
  distance: PropTypes.string,
};

ReviewSitAddressChange.defaultProps = {
  oldAddress: {},
  newAddress: {},
  contractorRemarks: '',
  distance: '',
};

export default ReviewSitAddressChange;
