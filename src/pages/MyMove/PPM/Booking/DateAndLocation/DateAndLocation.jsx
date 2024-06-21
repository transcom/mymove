import React, { useState, useEffect } from 'react';
import { useDispatch } from 'react-redux';
import { generatePath, useNavigate, useParams } from 'react-router-dom';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';

import { isBooleanFlagEnabled } from '../../../../../utils/featureFlags';

import DateAndLocationForm from 'components/Customer/PPM/Booking/DateAndLocationForm/DateAndLocationForm';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { customerRoutes, generalRoutes } from 'constants/routes';
import { shipmentTypes } from 'constants/shipments';
import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import { createMTOShipment, getAllMoves, patchMove, patchMTOShipment } from 'services/internalApi';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatDateForSwagger } from 'shared/dates';
import { updateMTOShipment, updateMove, updateAllMoves } from 'store/entities/actions';
import { DutyLocationShape } from 'types';
import { MoveShape, ServiceMemberShape } from 'types/customerShapes';
import { ShipmentShape } from 'types/shipment';
import SERVICE_MEMBER_AGENCIES from 'content/serviceMemberAgencies';
import { validatePostalCode } from 'utils/validation';
import { formatAddressForAPI } from 'utils/formatMtoShipment';

const DateAndLocation = ({ mtoShipment, serviceMember, destinationDutyLocation, move }) => {
  const [errorMessage, setErrorMessage] = useState(null);
  const [multiMove, setMultiMove] = useState(false);
  const navigate = useNavigate();
  const { moveId, shipmentNumber } = useParams();
  const dispatch = useDispatch();

  const includeCloseoutOffice =
    serviceMember.affiliation === SERVICE_MEMBER_AGENCIES.ARMY ||
    serviceMember.affiliation === SERVICE_MEMBER_AGENCIES.AIR_FORCE ||
    serviceMember.affiliation === SERVICE_MEMBER_AGENCIES.SPACE_FORCE;
  const isNewShipment = !mtoShipment?.id;

  useEffect(() => {
    isBooleanFlagEnabled('multi_move').then((enabled) => {
      setMultiMove(enabled);
    });
  }, []);

  const handleBack = () => {
    if (isNewShipment) {
      navigate(generatePath(customerRoutes.SHIPMENT_SELECT_TYPE_PATH, { moveId }));
    } else if (multiMove) {
      navigate(generatePath(customerRoutes.MOVE_HOME_PATH, { moveId }));
    } else {
      navigate(generalRoutes.HOME_PATH);
    }
  };

  const onShipmentSaveSuccess = (response, setSubmitting) => {
    // Update submitting state
    setSubmitting(false);

    // Update the shipment in the store
    dispatch(updateMTOShipment(response));

    // navigate to the next page
    navigate(
      generatePath(customerRoutes.SHIPMENT_PPM_ESTIMATED_WEIGHT_PATH, {
        moveId,
        mtoShipmentId: response.id,
      }),
    );
  };

  const handleSubmit = async (values, { setSubmitting }) => {
    setErrorMessage(null);

    const hasSecondaryPickupAddress = values.hasSecondaryPickupAddress === 'true';
    const hasSecondaryDestinationAddress = values.hasSecondaryDestinationAddress === 'true';

    const createOrUpdateShipment = {
      moveTaskOrderID: moveId,
      shipmentType: SHIPMENT_OPTIONS.PPM,
      ppmShipment: {
        pickupAddress: formatAddressForAPI(values.pickupAddress.address),
        hasSecondaryPickupAddress, // I think sending this is necessary so we know if the customer wants to clear their previously secondary ZIPs, or we could send nulls for those fields.
        destinationAddress: formatAddressForAPI(values.destinationAddress.address),
        hasSecondaryDestinationAddress,
        sitExpected: values.sitExpected === 'true',
        expectedDepartureDate: formatDateForSwagger(values.expectedDepartureDate),
      },
    };

    if (hasSecondaryPickupAddress && values.secondaryPickupAddress?.address) {
      createOrUpdateShipment.ppmShipment.secondaryPickupAddress = formatAddressForAPI(
        values.secondaryPickupAddress.address,
      );
    }

    if (hasSecondaryDestinationAddress && values.secondaryDestinationAddress?.address) {
      createOrUpdateShipment.ppmShipment.secondaryDestinationAddress = formatAddressForAPI(
        values.secondaryDestinationAddress.address,
      );
    }

    if (isNewShipment) {
      createMTOShipment(createOrUpdateShipment)
        .then((shipmentResponse) => {
          if (includeCloseoutOffice) {
            // Associate the selected closeout office with the move
            patchMove(move.id, { closeoutOfficeId: values.closeoutOffice.id }, move.eTag)
              .then((moveResponse) => {
                // Both create and patch were successful
                dispatch(updateMove(moveResponse));
                onShipmentSaveSuccess(shipmentResponse, setSubmitting);
              })
              .catch(() => {
                setSubmitting(false);
                // Still need to update the shipment in the store since it had a successful create
                dispatch(updateMTOShipment(shipmentResponse));
                setErrorMessage('There was an error attempting to create the move closeout office.');
              });
          } else {
            onShipmentSaveSuccess(shipmentResponse, setSubmitting);
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
              .then(async () => {
                const allMoves = await getAllMoves(serviceMember.id);
                dispatch(updateAllMoves(allMoves));
              })
              .catch(() => {
                setSubmitting(false);
                // Still need to update the shipment in the store since it had a successful update
                dispatch(updateMTOShipment(shipmentResponse));
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
