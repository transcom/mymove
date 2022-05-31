import React, { useState } from 'react';
import { generatePath, useHistory, useParams } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';

import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import closingPageStyles from 'pages/MyMove/PPM/Closeout/Closeout.module.scss';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import ScrollToTop from 'components/ScrollToTop';
import { shipmentTypes } from 'constants/shipments';
import AboutForm from 'components/Customer/PPM/Closeout/AboutForm/AboutForm';
import { customerRoutes, generalRoutes } from 'constants/routes';
import { selectMTOShipmentById } from 'store/entities/selectors';
import { validatePostalCode } from 'utils/validation';
import { formatDateForSwagger } from 'shared/dates';
import { getResponseError, patchMTOShipment } from 'services/internalApi';
import { updateMTOShipment } from 'store/entities/actions';

const About = () => {
  const [errorMessage, setErrorMessage] = useState();

  const history = useHistory();
  const { moveId, mtoShipmentId } = useParams();
  const dispatch = useDispatch();

  const mtoShipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));

  const handleBack = () => {
    history.push(generalRoutes.HOME_PATH);
  };

  const handleSubmit = async (values, { setSubmitting }) => {
    setErrorMessage(null);
    const hasReceivedAdvance = values.hasReceivedAdvance === 'true';
    const payload = {
      shipmentType: mtoShipment.shipmentType,
      ppmShipment: {
        id: mtoShipment.ppmShipment.id,
        actualMoveDate: formatDateForSwagger(values.actualMoveDate),
        actualPickupPostalCode: values.actualPickupPostalCode,
        actualDestinationPostalCode: values.actualDestinationPostalCode,
        hasReceivedAdvance,
        advanceAmountReceived: hasReceivedAdvance ? values.advanceAmountReceived * 100 : null,
      },
    };

    patchMTOShipment(mtoShipment.id, payload, mtoShipment.eTag)
      .then((response) => {
        setSubmitting(false);
        dispatch(updateMTOShipment(response));
        history.push(generatePath(customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_PATH, { moveId, mtoShipmentId }));
      })
      .catch((err) => {
        setSubmitting(false);
        setErrorMessage(getResponseError(err.response, 'Failed to update MTO shipment due to server error.'));
      });
  };

  return (
    <div className={ppmPageStyles.ppmPageStyle}>
      <ScrollToTop otherDep={errorMessage} />
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <ShipmentTag shipmentType={shipmentTypes.PPM} />
            <h1>About your PPM</h1>
            {errorMessage && (
              <Alert slim type="error">
                {errorMessage}
              </Alert>
            )}
            <div className={closingPageStyles['closing-section']}>
              <p>Finish moving this PPM before you start documenting it.</p>
              <p>To complete your PPM, you will:</p>
              <ul>
                <li>Upload weight tickets for each trip</li>
                <li>Upload receipts to document any expenses</li>
                <li>Upload receipts if you used short-term storage, so you can request reimbursement</li>
                <li>Upload any other documentation (such as proof of ownership for a trailer, if you used your own)</li>
                <li>Complete your PPM to send it to a counselor for review</li>
              </ul>
            </div>
            <AboutForm
              mtoShipment={mtoShipment}
              onSubmit={handleSubmit}
              onBack={handleBack}
              postalCodeValidator={validatePostalCode}
            />
          </Grid>
        </Grid>
      </GridContainer>
    </div>
  );
};

export default About;
