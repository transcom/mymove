import React from 'react';
import * as PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from '../ShipmentCard.module.scss';

import hhgShipmentCardStyles from './HHGShipmentCard.module.scss';

import ShipmentContainer from 'components/Office/ShipmentContainer';
import { SHIPMENT_OPTIONS } from 'shared/constants';

const HHGShipmentCard = ({ shipmentType }) => {
  return (
    <div className={styles.ShipmentCard} data-testid="shipment-display">
      <ShipmentContainer className={styles.container} shipmentType={shipmentType}>
        <div className={styles.ShipmentCardHeader}>
          <div>
            <h4>Shipment 1: HHG</h4>
            <p>#ABC123K-001</p>
          </div>
          <Button className={styles.editBtn} unstyled>
            Edit
          </Button>
        </div>

        <dl className={styles.shipmentCardSubsection}>
          <div className={styles.row}>
            <dt>Requested pickup date</dt>
            <dd>26 Mar 2020</dd>
          </div>
          <div className={styles.row}>
            <dt>Pickup location</dt>
            <dd>
              17 8th St <br />
              New York, NY 111111
            </dd>
          </div>
          <div className={styles.row}>
            <dt>Releasing agent</dt>
            <dd>
              Jo Xi <br />
              (555) 555-5555
              <br />
              jo.xi@example.com
            </dd>
          </div>
          <div className={styles.row}>
            <dt>Requested delivery date</dt>
            <dd>15 Apr 2020</dd>
          </div>
          <div className={styles.row}>
            <dt>Destination</dt>
            <dd>73523</dd>
          </div>
          <div className={styles.row}>
            <dt>Receiving agent</dt>
            <dd>
              Jo Xi <br />
              (555) 555-5555
              <br />
              jo.xi@example.com
            </dd>
          </div>
          <div className={styles.row}>
            <dt>Remarks</dt>
            <dd />
          </div>
        </dl>
        <p className={hhgShipmentCardStyles.remarksCell}>
          This is 500 characters of customer remarks right here. Lorem ipsum dolor sit amet, consectetur adipiscing
          elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud
          exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit
          in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non
          proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
        </p>
      </ShipmentContainer>
    </div>
  );
};

HHGShipmentCard.propTypes = {
  shipmentType: PropTypes.oneOf([
    SHIPMENT_OPTIONS.PPM,
    SHIPMENT_OPTIONS.HHG,
    SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC,
    SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
    SHIPMENT_OPTIONS.NTS,
  ]),
};

HHGShipmentCard.defaultProps = {
  shipmentType: SHIPMENT_OPTIONS.HHG,
};

export default HHGShipmentCard;
