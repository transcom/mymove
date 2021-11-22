import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import { generatePath } from 'react-router';
import { useHistory, useParams } from 'react-router-dom';

import PrimeUIShipmentUpdateReweighForm from './PrimeUIShipmentUpdateReweighForm';

import styles from 'components/Office/CustomerContactInfoForm/CustomerContactInfoForm.module.scss';
import { primeSimulatorRoutes } from 'constants/routes';

const PrimeUIShipmentUpdateReweigh = () => {
  const { moveCodeOrID } = useParams();
  const history = useHistory();

  const handleSubmit = () => {};

  const handleClose = () => {
    history.push(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
  };

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
