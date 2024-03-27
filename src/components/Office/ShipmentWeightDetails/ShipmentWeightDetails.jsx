import React from 'react';
import classnames from 'classnames';
import PropTypes from 'prop-types';
import { Button, Tag } from '@trussworks/react-uswds';

import DataTableWrapper from '../../DataTableWrapper/index';
import DataTable from '../../DataTable/index';

import styles from './ShipmentWeightDetails.module.scss';

import returnLowestValue from 'utils/returnLowestValue';
import { formatWeight } from 'utils/formatters';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { ShipmentOptionsOneOf } from 'types/shipment';
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';

const ShipmentWeightDetails = ({ estimatedWeight, initialWeight, shipmentInfo, handleRequestReweighModal }) => {
  const emDash = '\u2014';
  const lowestWeight = returnLowestValue(initialWeight, shipmentInfo.reweighWeight);
  const shipmentIsPPM = shipmentInfo.shipmentType === SHIPMENT_OPTIONS.PPM;

  const reweighHeader = (
    <div className={styles.shipmentWeight}>
      <span>Reweigh weight</span>
      {!shipmentInfo.reweighID && (
        <div className={styles.rightAlignButtonWrapper}>
          <Restricted to={permissionTypes.createReweighRequest}>
            <Button type="button" onClick={() => handleRequestReweighModal(shipmentInfo)} unstyled>
              Request reweigh
            </Button>
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
        {!shipmentIsPPM && (
          <DataTable
            columnHeaders={['Actual pro gear weight', 'Actual spouse pro gear weight']}
            dataRow={[
              shipmentInfo.shipmentActualProGearWeight && shipmentInfo.shipmentType !== SHIPMENT_OPTIONS.NTSR
                ? formatWeight(shipmentInfo.shipmentActualProGearWeight)
                : emDash,
              shipmentInfo.shipmentActualSpouseProGearWeight
                ? formatWeight(shipmentInfo.shipmentActualSpouseProGearWeight)
                : emDash,
            ]}
          />
        )}
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
    shipmentActualProGearWeight: PropTypes.number,
    shipmentActualSpouseProGearWeight: PropTypes.number,
  }).isRequired,
  handleRequestReweighModal: PropTypes.func.isRequired,
};

ShipmentWeightDetails.defaultProps = {
  estimatedWeight: null,
  initialWeight: null,
};

export default ShipmentWeightDetails;
