import React from 'react';
import classnames from 'classnames';
import PropTypes from 'prop-types';
import { Button, Tag } from '@trussworks/react-uswds';
import { connect } from 'react-redux';

import DataTableWrapper from '../../DataTableWrapper/index';
import DataTable from '../../DataTable/index';

import styles from './ShipmentWeightDetails.module.scss';

import returnLowestValue from 'utils/returnLowestValue';
import { formatWeight } from 'utils/formatters';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { ShipmentOptionsOneOf } from 'types/shipment';
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';
import { withContext } from 'shared/AppContext';
import { roleTypes } from 'constants/userRoles';

const ShipmentWeightDetails = ({
  estimatedWeight,
  initialWeight,
  shipmentInfo,
  handleRequestReweighModal,
  activeRole,
}) => {
  const emDash = '\u2014';
  const lowestWeight = returnLowestValue(initialWeight, shipmentInfo.reweighWeight);

  const reweighHeader = (
    <div className={styles.shipmentWeight}>
      <span>Reweigh weight</span>
      {!shipmentInfo.reweighID && (
        <div className={styles.rightAlignButtonWrapper}>
          <Restricted to={permissionTypes.createReweighRequest}>
            {activeRole !== roleTypes.SERVICES_COUNSELOR && (
              <Button type="button" onClick={() => handleRequestReweighModal(shipmentInfo)} unstyled>
                Request reweigh
              </Button>
            )}
          </Restricted>
        </div>
      )}
      {shipmentInfo.reweighID && !shipmentInfo.reweighWeight && <Tag>reweigh requested</Tag>}
      {shipmentInfo.reweighWeight && <Tag>reweighed</Tag>}
    </div>
  );

  return (
    <div className={classnames('maxw-tablet', styles.ShipmentWeightDetails)}>
      <DataTableWrapper className={classnames('table--data-point-group')}>
        <DataTable
          columnHeaders={['Estimated weight', 'Initial weight']}
          dataRow={[
            estimatedWeight && shipmentInfo.shipmentType !== SHIPMENT_OPTIONS.NTSR
              ? formatWeight(estimatedWeight)
              : emDash,
            initialWeight ? formatWeight(initialWeight) : emDash,
          ]}
        />
        <DataTable
          columnHeaders={[reweighHeader, 'Actual shipment weight']}
          dataRow={[
            shipmentInfo.reweighWeight ? formatWeight(shipmentInfo.reweighWeight) : emDash,
            lowestWeight ? formatWeight(lowestWeight) : emDash,
          ]}
        />
      </DataTableWrapper>
    </div>
  );
};

ShipmentWeightDetails.propTypes = {
  estimatedWeight: PropTypes.number,
  initialWeight: PropTypes.number,
  shipmentInfo: PropTypes.shape({
    shipmentID: PropTypes.string,
    ifMatchEtag: PropTypes.string,
    reweighID: PropTypes.string,
    reweighWeight: PropTypes.number,
    shipmentType: ShipmentOptionsOneOf.isRequired,
  }).isRequired,
  handleRequestReweighModal: PropTypes.func.isRequired,
};

ShipmentWeightDetails.defaultProps = {
  estimatedWeight: null,
  initialWeight: null,
};

// Checks user role such that Service Counselors cannot modify service items while on MTO read-only page
const mapStateToProps = (state) => {
  return {
    activeRole: state.auth.activeRole,
  };
};

export default withContext(connect(mapStateToProps)(ShipmentWeightDetails));
