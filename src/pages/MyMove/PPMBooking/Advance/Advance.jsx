import React, { useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { generatePath, useHistory, useParams } from 'react-router-dom';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';

import ppmBookingPageStyles from 'pages/MyMove/PPMBooking/PPMBooking.module.scss';
import { customerRoutes } from 'constants/routes';
import AdvanceForm from 'components/Customer/PPM/Booking/Advance/AdvanceForm';
import { shipmentTypes } from 'constants/shipments';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { getResponseError, patchMTOShipment } from 'services/internalApi';
import { updateMTOShipment } from 'store/entities/actions';
import { selectCurrentOrders, selectMTOShipmentById } from 'store/entities/selectors';
import { setFlashMessage } from 'store/flash/actions';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import ScrollToTop from 'components/ScrollToTop';
import { matchesOrdersType } from 'utils/orders';
import { ORDERS_TYPE } from 'constants/orders';

const Advance = () => {
  const [errorMessage, setErrorMessage] = useState();
  const history = useHistory();
  const { moveId, mtoShipmentId, shipmentNumber } = useParams();
  const dispatch = useDispatch();
  const mtoShipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));
  const orders = useSelector((state) => selectCurrentOrders(state));

  const handleBack = () => {
    history.push(generatePath(customerRoutes.SHIPMENT_PPM_ESTIMATED_INCENTIVE_PATH, { moveId, mtoShipmentId }));
  };

  const handleSubmit = async (values, { setSubmitting }) => {
    setErrorMessage(null);

    const hasRequestedAdvance = values.hasRequestedAdvance === 'true';

    const payload = {
      shipmentType: mtoShipment.shipmentType,
      ppmShipment: {
        id: mtoShipment.ppmShipment.id,
        hasRequestedAdvance,
        advanceAmountRequested: hasRequestedAdvance ? values.advanceAmountRequested * 100 : null,
      },
    };

    patchMTOShipment(mtoShipment.id, payload, mtoShipment.eTag)
      .then((response) => {
        setSubmitting(false);
        dispatch(updateMTOShipment(response));
        dispatch(
          setFlashMessage(
            'PPM_ONBOARDING_SUBMIT_SUCCESS',
            'success',
            'Review your info and submit your move request now, or come back and finish later.',
            'Details saved',
          ),
        );
        history.push(generatePath(customerRoutes.MOVE_REVIEW_PATH, { moveId }));
      })
      .catch((err) => {
        setSubmitting(false);

        setErrorMessage(getResponseError(err.response, 'Failed to update MTO shipment due to server error.'));
      });
  };

  const isRetireeOrSeparatee = matchesOrdersType(orders, ORDERS_TYPE.RETIREMENT, ORDERS_TYPE.SEPARATION);

  if (!mtoShipment || !orders) {
    return <LoadingPlaceholder />;
  }

  return (
    <div className={ppmBookingPageStyles.PPMBookingPage}>
      <ScrollToTop otherDep={errorMessage} />
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <ShipmentTag shipmentType={shipmentTypes.PPM} shipmentNumber={shipmentNumber} />
            <h1>Advances</h1>
            {isRetireeOrSeparatee && (
              <Alert slim type="info">
                People leaving the military may not be eligible to receive an advance, based on individual service
                policies. Your counselor can give you more information after you make your request.
              </Alert>
            )}
            {errorMessage && (
              <Alert slim type="error">
                {errorMessage}
              </Alert>
            )}
            <AdvanceForm mtoShipment={mtoShipment} onSubmit={handleSubmit} onBack={handleBack} />
          </Grid>
        </Grid>
      </GridContainer>
    </div>
  );
};

export default Advance;
