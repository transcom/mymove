import React from 'react';
// import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Alert } from '@trussworks/react-uswds';
import * as PropTypes from 'prop-types';

import DataTable from 'components/DataTable/index';
// import DataTableWrapper from 'components/DataTableWrapper/index';
import { formatAddress } from 'utils/shipmentDisplay';
import { AddressShape } from 'types';
// import styles from 'components/Office/ShipmentForm/ShipmentForm';

const ReviewSitAddressChange = ({ sitAddressUpdate }) => {
  const { oldAddress, newAddress, contractorRemarks, distance } = sitAddressUpdate;

  return (
    <div>
      <Alert type="warning">
        Requested final SIT delivery address is {distance} miles from the initial SIT delivery address. Approvals over
        50 miles will result in updated pricing for this shipment.
      </Alert>
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
    </div>
  );
};

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
