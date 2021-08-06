import React from 'react';
import PropTypes from 'prop-types';
import { Tag } from '@trussworks/react-uswds';

import styles from './ShipmentModificationTag.module.scss';

import { shipmentModificationTypes } from 'constants/shipments';

const ShipmentModificationTag = ({ shipmentModificationType }) => (
  <Tag className={styles.ShipmentModificationTag}>{shipmentModificationType}</Tag>
);

ShipmentModificationTag.propTypes = {
  shipmentModificationType: PropTypes.oneOfType([
    PropTypes.string,
    PropTypes.oneOf(Object.keys(shipmentModificationTypes)),
  ]).isRequired,
};

export default ShipmentModificationTag;
