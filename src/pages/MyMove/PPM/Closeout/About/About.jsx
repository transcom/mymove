import React, { useState, useEffect } from 'react';
import { generatePath, useNavigate, useParams } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';

import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import { shipmentTypes } from 'constants/shipments';
import AboutForm from 'components/Shared/PPM/Closeout/AboutForm/AboutForm';
import { customerRoutes } from 'constants/routes';
import { selectMTOShipmentById } from 'store/entities/selectors';
import { formatDateForSwagger } from 'shared/dates';
import { getResponseError, patchMTOShipment, getMTOShipmentsForMove } from 'services/internalApi';
import { updateMTOShipment } from 'store/entities/actions';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { isWeightTicketComplete } from 'utils/shipments';
import { PPM_TYPES } from 'shared/constants';
import { CUSTOMER_ERROR_MESSAGES } from 'constants/errorMessages';
import { APP_NAME } from 'constants/apps';

const About = () => {
  const [errorMessage, setErrorMessage] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  const navigate = useNavigate();
  const { moveId, mtoShipmentId } = useParams();
  const dispatch = useDispatch();

  const mtoShipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));

  const ppmShipment = mtoShipment?.ppmShipment || {};
  const { ppmType } = ppmShipment;
  const appName = APP_NAME.MYMOVE;

  useEffect(() => {
    getMTOShipmentsForMove(moveId)
      .then((response) => {
        dispatch(updateMTOShipment(response.mtoShipments[mtoShipmentId]));
      })
      .catch(() => {
        setErrorMessage('Failed to fetch shipment information');
      })
      .finally(() => {
        setIsLoading(false);
      });
  }, [moveId, mtoShipmentId, dispatch]);

  if (!mtoShipment || isLoading) {
    return <LoadingPlaceholder />;
  }

  const handleBack = () => {
    navigate(generatePath(customerRoutes.MOVE_HOME_PATH, { moveId }));
  };

  const handleSubmit = async (values, { setSubmitting }) => {
    setErrorMessage(null);
    const hasReceivedAdvance = values.hasReceivedAdvance === 'true';
    const payload = {
      shipmentType: mtoShipment.shipmentType,
      ppmShipment: {
        id: mtoShipment.ppmShipment.id,
        actualMoveDate: formatDateForSwagger(values.actualMoveDate),
        pickupAddress: values.pickupAddress,
        hasSecondaryPickupAddress: values.hasSecondaryPickupAddress === 'true',
        secondaryPickupAddress: values.hasSecondaryPickupAddress === 'true' ? values.secondaryPickupAddress : null,
        destinationAddress: values.destinationAddress,
        hasSecondaryDestinationAddress: values.hasSecondaryDestinationAddress === 'true',
        secondaryDestinationAddress:
          values.hasSecondaryDestinationAddress === 'true' ? values.secondaryDestinationAddress : null,
        hasReceivedAdvance,
        advanceAmountReceived: hasReceivedAdvance ? values.advanceAmountReceived * 100 : null,
        w2Address: values.w2Address,
      },
    };

    const handleErrorMessage = (error) => {
      if (error?.response?.status === 412) {
        setErrorMessage(CUSTOMER_ERROR_MESSAGES.PRECONDITION_FAILED);
      } else {
        setErrorMessage(getResponseError(error.response, 'Failed to update PPM shipment due to server error.'));
      }
    };

    patchMTOShipment(mtoShipment.id, payload, mtoShipment.eTag)
      .then((response) => {
        setSubmitting(false);
        dispatch(updateMTOShipment(response));

        let path;
        if (ppmType === PPM_TYPES.SMALL_PACKAGE) {
          path = generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, {
            moveId,
            mtoShipmentId,
          });
        } else if (response.ppmShipment.weightTickets.length === 0) {
          path = generatePath(customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_PATH, {
            moveId,
            mtoShipmentId,
          });
        } else if (!response.ppmShipment.weightTickets.some(isWeightTicketComplete)) {
          path = generatePath(customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_EDIT_PATH, {
            moveId,
            mtoShipmentId,
            weightTicketId: response.ppmShipment.weightTickets[0].id,
          });
        } else {
          path = generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, {
            moveId,
            mtoShipmentId,
          });
        }

        navigate(path);
      })
      .catch((err) => {
        setSubmitting(false);
        handleErrorMessage(err);
      });
  };

  if (!mtoShipment) {
    return <LoadingPlaceholder />;
  }

  return (
    <div className={ppmPageStyles.ppmPageStyle}>
      <NotificationScrollToTop dependency={errorMessage} />
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
            <AboutForm mtoShipment={mtoShipment} onSubmit={handleSubmit} onBack={handleBack} appName={appName} />
          </Grid>
        </Grid>
      </GridContainer>
    </div>
  );
};

export default About;
