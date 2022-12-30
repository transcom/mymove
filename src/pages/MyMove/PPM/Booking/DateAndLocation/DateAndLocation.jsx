import React, { useState } from 'react';
import { useDispatch } from 'react-redux';
import { generatePath, useNavigate, useParams } from 'react-router-dom';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';

import DateAndLocationForm from 'components/Customer/PPM/Booking/DateAndLocationForm/DateAndLocationForm';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { customerRoutes, generalRoutes } from 'constants/routes';
import { shipmentTypes } from 'constants/shipments';
import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import { createMTOShipment, patchMTOShipment } from 'services/internalApi';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatDateForSwagger } from 'shared/dates';
import { updateMTOShipment } from 'store/entities/actions';
import { DutyLocationShape } from 'types';
import { ServiceMemberShape } from 'types/customerShapes';
import { ShipmentShape } from 'types/shipment';
import { validatePostalCode } from 'utils/validation';

const DateAndLocation = ({ mtoShipment, serviceMember, destinationDutyLocation }) => {
  const [errorMessage, setErrorMessage] = useState(null);
  const navigate = useNavigate();
  const { moveId, shipmentNumber } = useParams();
  const dispatch = useDispatch();

  const isNewShipment = !mtoShipment?.id;
  const handleBack = () => {
    if (isNewShipment) {
      navigate(generatePath(customerRoutes.SHIPMENT_SELECT_TYPE_PATH, { moveId }));
    } else {
      navigate(generalRoutes.HOME_PATH);
    }
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
        .then((response) => {
          setSubmitting(false);
          dispatch(updateMTOShipment(response));
          navigate(
            generatePath(customerRoutes.SHIPMENT_PPM_ESTIMATED_WEIGHT_PATH, {
              moveId,
              mtoShipmentId: response.id,
            }),
          );
        })
        .catch(() => {
          setSubmitting(false);
          setErrorMessage('There was an error attempting to create your shipment.');
        });
    } else {
      createOrUpdateShipment.id = mtoShipment.id;
      createOrUpdateShipment.ppmShipment.id = mtoShipment.ppmShipment?.id;

      patchMTOShipment(mtoShipment.id, createOrUpdateShipment, mtoShipment.eTag)
        .then((response) => {
          setSubmitting(false);
          dispatch(updateMTOShipment(response));
          navigate(
            generatePath(customerRoutes.SHIPMENT_PPM_ESTIMATED_WEIGHT_PATH, {
              moveId,
              mtoShipmentId: response.id,
            }),
          );
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
};

DateAndLocation.defaultProps = {
  mtoShipment: {},
};

export default DateAndLocation;
