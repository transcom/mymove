import React, { useState, useEffect } from 'react';
import { useDispatch } from 'react-redux';
import { generatePath, useNavigate, useParams, useLocation } from 'react-router-dom';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';

import { isBooleanFlagEnabled } from '../../../../utils/featureFlags';

import MobileHomeShipmentForm from 'components/Customer/MobileHomeShipment/MobileHomeShipmentForm/MobileHomeShipmentForm';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { customerRoutes, generalRoutes } from 'constants/routes';
import pageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import { createMTOShipment, patchMTOShipment } from 'services/internalApi';
import { SHIPMENT_OPTIONS, SHIPMENT_TYPES } from 'shared/constants';
import { updateMTOShipment } from 'store/entities/actions';
import { DutyLocationShape } from 'types';
import { MoveShape, ServiceMemberShape } from 'types/customerShapes';
import { ShipmentShape } from 'types/shipment';
import { validatePostalCode } from 'utils/validation';
import { toTotalInches } from 'utils/formatMtoShipment';

const MobileHomeShipmentCreate = ({
  mtoShipment,
  serviceMember,
  destinationDutyLocation,
  move,
  serviceMemberMoves,
}) => {
  const [errorMessage, setErrorMessage] = useState(null);
  const [multiMove, setMultiMove] = useState(false);

  const navigate = useNavigate();
  const { moveId } = useParams();
  const dispatch = useDispatch();
  const location = useLocation();
  const searchParams = new URLSearchParams(location.search);
  const shipmentNumber = searchParams.get('shipmentNumber');
  const isEditPage = location?.pathname?.includes('/edit');

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
    const baseMtoShipment = mtoShipment?.id ? mtoShipment : response;
    const data = {
      ...baseMtoShipment,
      mobileHomeShipment: response?.mobileHomeShipment,
      shipmentType: response?.shipmentType,
      customerRemarks: response?.customerRemarks,
      eTag: response?.eTag,
    };
    const currentMove = serviceMemberMoves?.currentMove[0];

    if (currentMove?.mtoShipments?.length) {
      currentMove?.mtoShipments?.forEach((element, idx) => {
        if (element.id === response.id) {
          currentMove.mtoShipments[idx] = data;
        }
      });
    }

    dispatch(updateMTOShipment(data));

    // navigate to the next page
    navigate(
      generatePath(customerRoutes.SHIPMENT_MOBILE_HOME_LOCATION_INFO, {
        moveId,
        mtoShipmentId: response.id,
      }),
    );
  };

  const handleSubmit = async (values, { setSubmitting }) => {
    setErrorMessage(null);
    const totalLengthInInches = toTotalInches(values.lengthFeet, values.lengthInches);
    const totalWidthInInches = toTotalInches(values.widthFeet, values.widthInches);
    const totalHeightInInches = toTotalInches(values.heightFeet, values.heightInches);

    const createOrUpdateShipment = {
      moveTaskOrderID: moveId,
      shipmentType: SHIPMENT_TYPES.MOBILE_HOME,
      mobileHomeShipment: {
        year: Number(values.year),
        make: values.make,
        model: values.model,
        lengthInInches: totalLengthInInches,
        widthInInches: totalWidthInInches,
        heightInInches: totalHeightInInches,
      },
      customerRemarks: values.customerRemarks,
    };

    if (isNewShipment) {
      createMTOShipment(createOrUpdateShipment)
        .then((shipmentResponse) => {
          onShipmentSaveSuccess(shipmentResponse, setSubmitting);
        })
        .catch((e) => {
          setSubmitting(false);
          const { response } = e;
          let errorMsg = 'There was an error attempting to create your shipment.';
          if (response?.body?.invalidFields) {
            const keys = Object.keys(response?.body?.invalidFields);
            const firstError = response?.body?.invalidFields[keys[0]][0];
            errorMsg = firstError;
          }
          setErrorMessage(errorMsg);
        });
    } else {
      createOrUpdateShipment.id = mtoShipment.id;
      createOrUpdateShipment.mobileHomeShipment.id = mtoShipment.mobileHomeShipment?.id;
      patchMTOShipment(mtoShipment.id, createOrUpdateShipment, mtoShipment.eTag)
        .then((shipmentResponse) => {
          onShipmentSaveSuccess(shipmentResponse, setSubmitting);
        })
        .catch((e) => {
          setSubmitting(false);
          const { response } = e;
          let errorMsg = 'There was an error attempting to update your shipment.';
          if (response?.body?.invalidFields) {
            const keys = Object.keys(response?.body?.invalidFields);
            const firstError = response?.body?.invalidFields[keys[0]][0];
            errorMsg = firstError;
          }
          setErrorMessage(errorMsg);
        });
    }
  };

  return (
    <div className={pageStyles.ppmPageStyle}>
      <NotificationScrollToTop dependency={errorMessage} />
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <ShipmentTag shipmentType={SHIPMENT_OPTIONS.MOBILE_HOME} shipmentNumber={shipmentNumber} />
            <h1>Mobile Home details and measurements</h1>
            {errorMessage && (
              <Alert headingLevel="h4" slim type="error">
                {errorMessage}
              </Alert>
            )}
            <MobileHomeShipmentForm
              mtoShipment={mtoShipment}
              serviceMember={serviceMember}
              destinationDutyLocation={destinationDutyLocation}
              move={move}
              onSubmit={handleSubmit}
              onBack={handleBack}
              postalCodeValidator={validatePostalCode}
              isEditPage={isEditPage}
            />
          </Grid>
        </Grid>
      </GridContainer>
    </div>
  );
};

MobileHomeShipmentCreate.propTypes = {
  mtoShipment: ShipmentShape,
  serviceMember: ServiceMemberShape.isRequired,
  destinationDutyLocation: DutyLocationShape.isRequired,
  move: MoveShape,
};

MobileHomeShipmentCreate.defaultProps = {
  move: {},
  mtoShipment: {},
};

export default MobileHomeShipmentCreate;
