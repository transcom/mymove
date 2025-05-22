import React from 'react';
import { generatePath, useNavigate, useParams } from 'react-router-dom';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import { useSelector } from 'react-redux';
import classnames from 'classnames';

import styles from './EstimatedIncentive.module.scss';

import formStyles from 'styles/form.module.scss';
import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import { ButtonUsa as Button } from 'shared/standardUI/Buttons/ButtonUsa';
import { shipmentTypes } from 'constants/shipments';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { selectMTOShipmentById } from 'store/entities/selectors';
import { customerRoutes } from 'constants/routes';
import EstimatedIncentiveDetails from 'components/Customer/PPM/Booking/EstimatedIncentiveDetails/EstimatedIncentiveDetails';

const EstimatedIncentive = () => {
  const navigate = useNavigate();
  const { moveId, mtoShipmentId, shipmentNumber } = useParams();
  const shipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));
  const handleBack = () => {
    navigate(generatePath(customerRoutes.SHIPMENT_PPM_ESTIMATED_WEIGHT_PATH, { moveId, mtoShipmentId }));
  };

  const handleNext = () => {
    navigate(generatePath(customerRoutes.SHIPMENT_PPM_ADVANCES_PATH, { moveId, mtoShipmentId }));
  };

  return (
    <div className={classnames(ppmPageStyles.ppmPageStyle, styles.EstimatedIncentive)}>
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <ShipmentTag shipmentType={shipmentTypes.PPM} shipmentNumber={shipmentNumber} />
            <h1>Estimated incentive</h1>
            <EstimatedIncentiveDetails shipment={shipment} />
            <div className={formStyles.buttonContainer}>
              <Button className={formStyles.backButton} type="button" onClick={handleBack} secondary outline>
                Back
              </Button>
              <Button className={formStyles.saveButton} type="button" onClick={handleNext}>
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
