import React from 'react';
import { Tag, Button, Grid, GridContainer } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import PropTypes from 'prop-types';

import styles from 'components/Customer/Review/IncompleteShipmentToolTip/IncompleteShipmentToolTip.module.scss';

const IncompleteShipmentToolTip = ({ onClick, shipmentLabel, moveCodeLabel, shipmentTypeLabel }) => {
  return (
    <GridContainer className={styles.gridContainerIncompleteToolTip}>
      <Grid row>
        <Grid col="fill" tablet={{ col: 'auto' }}>
          <Tag>Incomplete</Tag>
        </Grid>
        <Grid col="auto" className={styles.buttonContainer}>
          <Button
            title="Help about incomplete shipment"
            type="button"
            onClick={() => onClick(shipmentLabel, moveCodeLabel, shipmentTypeLabel)}
            unstyled
            className="{styles.buttonRight}"
          >
            <FontAwesomeIcon icon={['far', 'circle-question']} />
          </Button>
        </Grid>
      </Grid>
    </GridContainer>
  );
};

IncompleteShipmentToolTip.propTypes = {
  onClick: PropTypes.func.isRequired,
  shipmentLabel: PropTypes.string.isRequired,
  moveCodeLabel: PropTypes.string.isRequired,
  shipmentTypeLabel: PropTypes.string.isRequired,
};

export default IncompleteShipmentToolTip;
