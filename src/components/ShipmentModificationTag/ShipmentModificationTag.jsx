import React from 'react';
import PropTypes from 'prop-types';
import { Tag } from '@trussworks/react-uswds';

import styles from './ShipmentModificationTag.module.scss';

import { shipmentModificationTypes } from 'constants/shipments';

const ShipmentModificationTag = ({ shipmentModificationType }) =>
  shipmentModificationType !== null && <Tag className={styles.ShipmentModificationTag}>{shipmentModificationType}</Tag>;

ShipmentModificationTag.propTypes = {
  shipmentModificationType: PropTypes.oneOf(Object.keys(shipmentModificationTypes)),
};
ShipmentModificationTag.defaultProps = {
  shipmentModificationType: null,
};

export default ShipmentModificationTag;
