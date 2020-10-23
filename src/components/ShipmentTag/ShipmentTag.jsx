/* eslint-disable security/detect-object-injection */
import React from 'react';
import PropTypes from 'prop-types';
import { Tag } from '@trussworks/react-uswds';

import styles from './ShipmentTag.module.scss';

import { shipmentTypes as shipmentTypeLabels } from 'content/shipments';
import { shipmentTypes } from 'constants/shipments';

const ShipmentTag = ({ shipmentType, shipmentNumber }) => (
  <Tag className={`${styles.ShipmentTag} ${styles[`${shipmentTypes[shipmentType]}`]}`}>
    {shipmentTypeLabels[shipmentType]}
    {shipmentNumber && ` #${shipmentNumber}`}
  </Tag>
);

ShipmentTag.propTypes = {
  shipmentType: PropTypes.oneOf(Object.keys(shipmentTypes)).isRequired,
  shipmentNumber: PropTypes.string,
};

ShipmentTag.defaultProps = {
  shipmentNumber: null,
};

export default ShipmentTag;
