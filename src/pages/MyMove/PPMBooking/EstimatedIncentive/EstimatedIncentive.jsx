import React from 'react';
import { generatePath, useHistory, useParams } from 'react-router-dom';
import { GridContainer, Grid, Button } from '@trussworks/react-uswds';
import { useSelector } from 'react-redux';
import classnames from 'classnames';

import styles from './EstimatedIncentive.module.scss';

import ppmBookingStyles from 'components/Customer/PPMBooking/PPMBooking.module.scss';
import ppmBookingPageStyles from 'pages/MyMove/PPMBooking/PPMBooking.module.scss';
import { shipmentTypes } from 'constants/shipments';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { selectMTOShipmentById } from 'store/entities/selectors';
import { customerRoutes } from 'constants/routes';
import EstimatedIncentiveDetails from 'components/Customer/PPMBooking/EstimatedIncentiveDetails/EstimatedIncentiveDetails';

const EstimatedIncentive = () => {
  const history = useHistory();
  const { moveId, mtoShipmentId, shipmentNumber } = useParams();
  const shipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));
  const handleBack = () => {
    history.goBack();
  };

  const handleNext = () => {
    history.push(generatePath(customerRoutes.SHIPMENT_PPM_ADVANCES_PATH, { moveId, mtoShipmentId }));
  };

  return (
    <div className={classnames(ppmBookingPageStyles.PPMBookingPage, styles.EstimatedIncentive)}>
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <ShipmentTag shipmentType={shipmentTypes.PPM} shipmentNumber={shipmentNumber} />
            <h1>Estimated incentive</h1>
            <EstimatedIncentiveDetails shipment={shipment} />
            <div className={ppmBookingStyles.buttonContainer}>
              <Button className={ppmBookingStyles.backButton} type="button" onClick={handleBack} secondary outline>
                Back
              </Button>
              <Button className={ppmBookingStyles.saveButton} type="button" onClick={handleNext}>
                Next
              </Button>
            </div>
          </Grid>
        </Grid>
      </GridContainer>
    </div>
  );
};

export default EstimatedIncentive;
