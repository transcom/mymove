import React from 'react';
import classnames from 'classnames';
import PropTypes from 'prop-types';
import { Button, Tag } from '@trussworks/react-uswds';

import DataPointGroup from '../../DataPointGroup/index';
import DataPoint from '../../DataPoint/index';

import styles from './ShipmentWeightDetails.module.scss';

import returnLowestValue from 'utils/returnLowestValue';
import { formatWeight } from 'shared/formatters';

const ShipmentWeightDetails = ({ estimatedWeight, actualWeight, shipmentInfo, handleRequestReweighModal }) => {
  const lowestWeight = returnLowestValue(actualWeight, shipmentInfo.reweighWeight);
  const reweighHeader = (
    <div className={styles.shipmentWeight}>
      <span>Shipment weight</span>
      {!shipmentInfo.reweighID && (
        <div className={styles.rightAlignButtonWrapper}>
          <Button type="button" onClick={() => handleRequestReweighModal(shipmentInfo)} unstyled>
            Request reweigh
          </Button>
        </div>
      )}
      {shipmentInfo.reweighID && !shipmentInfo.reweighWeight && <Tag>reweigh requested</Tag>}
      {shipmentInfo.reweighWeight && <Tag>reweighed</Tag>}
    </div>
  );
  return (
    <div className={classnames('maxw-tablet', styles.ShipmentWeightDetails)}>
      <DataPointGroup className="maxw-mobile">
        <DataPoint
          columnHeaders={['Estimated weight']}
          dataRow={estimatedWeight ? [formatWeight(estimatedWeight)] : ['']}
        />
      </DataPointGroup>
      <DataPointGroup className="maxw-mobile">
        <DataPoint columnHeaders={[reweighHeader]} dataRow={lowestWeight ? [formatWeight(lowestWeight)] : ['']} />
      </DataPointGroup>
    </div>
  );
};

ShipmentWeightDetails.propTypes = {
  estimatedWeight: PropTypes.number,
  actualWeight: PropTypes.number,
  shipmentInfo: PropTypes.shape({
    shipmentID: PropTypes.string,
    ifMatchEtag: PropTypes.string,
    reweighID: PropTypes.string,
    reweighWeight: PropTypes.number,
  }).isRequired,
  handleRequestReweighModal: PropTypes.func.isRequired,
};

ShipmentWeightDetails.defaultProps = {
  estimatedWeight: null,
  actualWeight: null,
};

export default ShipmentWeightDetails;
