import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import PrimeUIShipmentUpdateReweighForm from './PrimeUIShipmentUpdateReweighForm';

import styles from 'components/Office/CustomerContactInfoForm/CustomerContactInfoForm.module.scss';

const PrimeUIShipmentUpdateReweigh = () => {
  const handleSubmit = () => {};

  const handleClose = () => {};

  return (
    <div className={styles.tabContent}>
      <div className={styles.container}>
        <GridContainer className={styles.gridContainer}>
          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }}>
              <PrimeUIShipmentUpdateReweighForm handleSubmit={handleSubmit} handleClose={handleClose} />
            </Grid>
          </Grid>
        </GridContainer>
      </div>
    </div>
  );
};

export default PrimeUIShipmentUpdateReweigh;
