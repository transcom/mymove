import React, { useState } from 'react';
import { generatePath, useHistory, useParams } from 'react-router-dom';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';

import styles from 'pages/MyMove/PPMBooking/DateAndLocation/DateAndLocation.module.scss';
import { MtoShipmentShape, ServiceMemberShape } from 'types/customerShapes';
import { DutyStationShape } from 'types';
import DateAndLocationForm from 'components/Customer/PPMBooking/DateAndLocationForm/DateAndLocationForm';
import { validatePostalCode } from 'utils/validation';
import { customerRoutes, generalRoutes } from 'constants/routes';
import { createMTOShipment, patchMTOShipment } from 'services/internalApi';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { updateMTOShipment } from 'sagas/entities';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { shipmentTypes } from 'constants/shipments';
import { formatDateForSwagger } from 'shared/dates';

const DateAndLocation = ({ mtoShipment, serviceMember, destinationDutyLocation }) => {
  const [errorMessage, setErrorMessage] = useState();
  const history = useHistory();
  const { moveId, shipmentNumber } = useParams();

  const isNewShipment = !mtoShipment?.id;

  const handleBack = () => {
    if (isNewShipment) {
      history.push(generatePath(customerRoutes.SHIPMENT_SELECT_TYPE_PATH, { moveId }));
    }

    history.push(generalRoutes.HOME_PATH);
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
          updateMTOShipment(response);
          history.push(
            generatePath(customerRoutes.SHIPMENT_PPM_ESTIMATED_INCENTIVE_PATH, {
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
          updateMTOShipment(response);
          history.push(
            generatePath(customerRoutes.SHIPMENT_PPM_ESTIMATED_INCENTIVE_PATH, {
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
    <div className={styles.DateAndLocation}>
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <ShipmentTag shipmentType={shipmentTypes.PPM} shipmentNumber={shipmentNumber} />
            <h1>PPM date & location</h1>
            {errorMessage && (
              <Alert slim type="error">
                {errorMessage}
              </Alert>
            )}
            <DateAndLocationForm
              mtoShipment={mtoShipment}
              serviceMember={serviceMember}
              destinationDutyStation={destinationDutyLocation}
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
  mtoShipment: MtoShipmentShape,
  serviceMember: ServiceMemberShape.isRequired,
  destinationDutyLocation: DutyStationShape.isRequired,
};

DateAndLocation.defaultProps = {
  mtoShipment: {},
};

export default DateAndLocation;
