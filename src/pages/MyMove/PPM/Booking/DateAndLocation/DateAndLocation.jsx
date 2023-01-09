import React, { useState } from 'react';
import { useDispatch } from 'react-redux';
import { generatePath, useHistory, useParams } from 'react-router-dom';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';

import DateAndLocationForm from 'components/Customer/PPM/Booking/DateAndLocationForm/DateAndLocationForm';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { customerRoutes, generalRoutes } from 'constants/routes';
import { shipmentTypes } from 'constants/shipments';
import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import { createMTOShipment, patchMove, patchMTOShipment } from 'services/internalApi';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatDateForSwagger } from 'shared/dates';
import { updateMTOShipment, updateMove } from 'store/entities/actions';
import { DutyLocationShape } from 'types';
import { MoveShape, ServiceMemberShape } from 'types/customerShapes';
import { ShipmentShape } from 'types/shipment';
import { validatePostalCode } from 'utils/validation';

const DateAndLocation = ({ mtoShipment, serviceMember, destinationDutyLocation, move }) => {
  const [errorMessage, setErrorMessage] = useState(null);
  const history = useHistory();
  const { moveId, shipmentNumber } = useParams();
  const dispatch = useDispatch();

  const includeCloseoutOffice = serviceMember.affiliation === 'ARMY' || serviceMember.affiliation === 'AIR_FORCE';
  const isNewShipment = !mtoShipment?.id;
  const handleBack = () => {
    if (isNewShipment) {
      history.push(generatePath(customerRoutes.SHIPMENT_SELECT_TYPE_PATH, { moveId }));
    } else {
      history.push(generalRoutes.HOME_PATH);
    }
  };

  const onShipmentSaveSuccess = (response, setSubmitting) => {
    // Update submitting state
    setSubmitting(false);

    // Update the shipment in the store
    dispatch(updateMTOShipment(response));

    // navigate to the next page
    history.push(
      generatePath(customerRoutes.SHIPMENT_PPM_ESTIMATED_WEIGHT_PATH, {
        moveId,
        mtoShipmentId: response.id,
      }),
    );
  };

  const handleSubmit = async (values, { setSubmitting }) => {
    setErrorMessage(null);

    const hasSecondaryPickupPostalCode = values.hasSecondaryPickupPostalCode === 'true';
    const hasSecondaryDestinationPostalCode = values.hasSecondaryDestinationPostalCode === 'true';

    const createOrUpdateShipment = {
      moveTaskOrderID: moveId,
      shipmentType: SHIPMENT_OPTIONS.PPM,
      ppmShipment: {
        pickupPostalCode: values.pickupPostalCode,
        hasSecondaryPickupPostalCode, // I think sending this is necessary so we know if the customer wants to clear their previously secondary ZIPs, or we could send nulls for those fields.
        secondaryPickupPostalCode: hasSecondaryPickupPostalCode ? values.secondaryPickupPostalCode : null,
        destinationPostalCode: values.destinationPostalCode,
        hasSecondaryDestinationPostalCode,
        secondaryDestinationPostalCode: hasSecondaryDestinationPostalCode
          ? values.secondaryDestinationPostalCode
          : null,
        sitExpected: values.sitExpected === 'true',
        expectedDepartureDate: formatDateForSwagger(values.expectedDepartureDate),
      },
    };

    if (isNewShipment) {
      createMTOShipment(createOrUpdateShipment)
        .then((createResponse) => {
          if (includeCloseoutOffice) {
            // Associate the selected closeout office with the move
            patchMove(move.id, { closeoutOfficeId: values.closeoutOffice.id }, move.eTag)
              .then((moveResponse) => {
                // Both create and patch were successful
                dispatch(updateMove(moveResponse));
                onShipmentSaveSuccess(createResponse, setSubmitting);
              })
              .catch(() => {
                setSubmitting(false);
                setErrorMessage('There was an error attempting to update the move closeout office.');
              });
          } else {
            onShipmentSaveSuccess(createResponse, setSubmitting);
          }
        })
        .catch(() => {
          setSubmitting(false);
          setErrorMessage('There was an error attempting to create your shipment.');
        });
    } else {
      createOrUpdateShipment.id = mtoShipment.id;
      createOrUpdateShipment.ppmShipment.id = mtoShipment.ppmShipment?.id;

      patchMTOShipment(mtoShipment.id, createOrUpdateShipment, mtoShipment.eTag)
        .then((shipmentResponse) => {
          if (includeCloseoutOffice) {
            // Associate the selected closeout office with the move
            patchMove(move.id, { closeoutOfficeId: values.closeoutOffice.id }, move.eTag)
              .then((moveResponse) => {
                dispatch(updateMove(moveResponse));
                onShipmentSaveSuccess(shipmentResponse, setSubmitting);
              })
              .catch(() => {
                setSubmitting(false);
                setErrorMessage('There was an error attempting to update the move closeout office.');
              });
          } else {
            onShipmentSaveSuccess(shipmentResponse, setSubmitting);
          }
        })
        .catch(() => {
          setSubmitting(false);
          setErrorMessage('There was an error attempting to update your shipment.');
        });
    }
  };

  return (
    <div className={ppmPageStyles.ppmPageStyle}>
      <NotificationScrollToTop dependency={errorMessage} />
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <ShipmentTag shipmentType={shipmentTypes.PPM} shipmentNumber={shipmentNumber} />
            <h1>PPM date & location</h1>
            {errorMessage && (
              <Alert headingLevel="h4" slim type="error">
                {errorMessage}
              </Alert>
            )}
            <DateAndLocationForm
              mtoShipment={mtoShipment}
              serviceMember={serviceMember}
              destinationDutyLocation={destinationDutyLocation}
              move={move}
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

DateAndLocation.propTypes = {
  mtoShipment: ShipmentShape,
  serviceMember: ServiceMemberShape.isRequired,
  destinationDutyLocation: DutyLocationShape.isRequired,
  move: MoveShape,
};

DateAndLocation.defaultProps = {
  move: {},
  mtoShipment: {},
};

export default DateAndLocation;
