import React, { useState, useEffect } from 'react';
import { generatePath, useNavigate, useParams } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';
import classnames from 'classnames';

import { isBooleanFlagEnabled } from '../../../../../utils/featureFlags';

import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import closingPageStyles from 'pages/MyMove/PPM/Closeout/Closeout.module.scss';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import { shipmentTypes } from 'constants/shipments';
import AboutForm from 'components/Customer/PPM/Closeout/AboutForm/AboutForm';
import { customerRoutes } from 'constants/routes';
import { selectMTOShipmentById } from 'store/entities/selectors';
import { validatePostalCode } from 'utils/validation';
import { formatDateForSwagger } from 'shared/dates';
import { getResponseError, patchMTOShipment, getMTOShipmentsForMove } from 'services/internalApi';
import { updateMTOShipment } from 'store/entities/actions';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { isWeightTicketComplete } from 'utils/shipments';

const About = () => {
  const [errorMessage, setErrorMessage] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  const navigate = useNavigate();
  const { moveId, mtoShipmentId } = useParams();
  const dispatch = useDispatch();

  const mtoShipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));
  const [multiMove, setMultiMove] = useState(false);

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

    isBooleanFlagEnabled('multi_move').then((enabled) => {
      setMultiMove(enabled);
    });
  }, [moveId, mtoShipmentId, dispatch]);

  if (!mtoShipment || isLoading) {
    return <LoadingPlaceholder />;
  }

  const handleBack = () => {
    if (multiMove) {
      navigate(generatePath(customerRoutes.MOVE_HOME_PATH, { moveId }));
    } else {
      navigate(customerRoutes.MOVE_HOME_PAGE);
    }
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
        w2Address: values.w2Address,
      },
    };

    patchMTOShipment(mtoShipment.id, payload, mtoShipment.eTag)
      .then((response) => {
        setSubmitting(false);
        dispatch(updateMTOShipment(response));

        let path;
        if (response.ppmShipment.weightTickets.length === 0) {
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
        setErrorMessage(getResponseError(err.response, 'Failed to update MTO shipment due to server error.'));
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
            <div className={classnames(closingPageStyles['closing-section'], closingPageStyles['about-ppm'])}>
              <p>Finish moving this PPM before you start documenting it.</p>
              <h2>How to complete your PPM</h2>
              <p>To complete your PPM, you will:</p>
              <ul>
                <li>Upload weight tickets for each trip</li>
                <li>Upload receipts to document any expenses</li>
                <li>Upload receipts if you used short-term storage, so you can request reimbursement</li>
                <li>Upload any other documentation (such as proof of ownership for a trailer, if you used your own)</li>
                <li>Complete your PPM to send it to a counselor for review</li>
              </ul>
              <h2>About your final payment</h2>
              <p>Your final payment will be:</p>
              <ul>
                <li>based on your final incentive</li>
                <li>modified by expenses submitted (authorized expenses reduce your tax burden)</li>
                <li>minus any taxes withheld (the IRS considers your incentive to be taxable income)</li>
                <li>plus any reimbursements you receive</li>
              </ul>
              <p>
                Verified expenses reduce the taxable income you report to the IRS on form W-2. They may not be claimed
                again as moving expenses. Federal tax withholding will be deducted from the profit (entitlement less
                eligible operating expenses.)
              </p>
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
