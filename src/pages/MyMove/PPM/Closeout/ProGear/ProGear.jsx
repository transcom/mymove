import React from 'react';
import { useHistory } from 'react-router-dom';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import ScrollToTop from 'components/ScrollToTop';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { shipmentTypes } from 'constants/shipments';
import { generalRoutes } from 'constants/routes';
import closingPageStyles from 'pages/MyMove/PPM/Closeout/Closeout.module.scss';
import ProGearForm from 'components/Customer/PPM/Closeout/ProGearForm/ProGearForm';

const handleSubmit = () => {};

const ProGear = () => {
  const history = useHistory();
  const handleBack = () => {
    history.push(generalRoutes.HOME_PATH);
  };
  return (
    <div className={ppmPageStyles.ppmPageStyle}>
      <ScrollToTop />
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <ShipmentTag shipmentType={shipmentTypes.PPM} />
            <h1>Pro-gear</h1>
            <div className={closingPageStyles['closing-section']}>
              <p>
                If you moved pro-gear for yourself or your spouse as part of this PPM, document the total weight here.
                Reminder: This pro-gear should be included in your total weight moved.
              </p>
            </div>
            <ProGearForm onBack={handleBack} onSubmit={handleSubmit} />
          </Grid>
        </Grid>
      </GridContainer>
    </div>
  );
};

export default ProGear;
