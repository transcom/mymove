import React, { useState } from 'react';
import { useDispatch } from 'react-redux';
import { generatePath, useNavigate, useParams, useLocation } from 'react-router-dom';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';

import BoatShipmentForm from 'components/Customer/BoatShipment/BoatShipmentForm/BoatShipmentForm';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import scrollToTop from 'shared/scrollToTop';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { customerRoutes } from 'constants/routes';
import { boatShipmentTypes } from 'constants/shipments';
import pageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import { createMTOShipment, patchMTOShipment, deleteMTOShipment, getAllMoves } from 'services/internalApi';
import { SHIPMENT_OPTIONS, SHIPMENT_TYPES } from 'shared/constants';
import { updateMTOShipment, updateAllMoves } from 'store/entities/actions';
import { DutyLocationShape } from 'types';
import { MoveShape, ServiceMemberShape } from 'types/customerShapes';
import { ShipmentShape } from 'types/shipment';
import { validatePostalCode } from 'utils/validation';
import { toTotalInches } from 'utils/formatMtoShipment';
import BoatShipmentConfirmationModal from 'components/Customer/BoatShipment/BoatShipmentConfirmationModal/BoatShipmentConfirmationModal';

const BoatShipmentCreate = ({ mtoShipment, serviceMember, destinationDutyLocation, move, serviceMemberMoves }) => {
  const [errorMessage, setErrorMessage] = useState(null);
  const [showBoatConfirmationModal, setShowBoatConfirmationModal] = useState(false);
  const [isDimensionsMeetReq, setIsDimensionsMeetReq] = useState(false);
  const [boatShipmentObj, setBoatShipmentObj] = useState(null);
  const [submitValues, setSubmitValues] = useState(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);

  const navigate = useNavigate();
  const { moveId } = useParams();
  const dispatch = useDispatch();
  const location = useLocation();
  const searchParams = new URLSearchParams(location.search);
  const shipmentNumber = searchParams.get('shipmentNumber');
  const isEditPage = location?.pathname?.includes('/edit');

  const isNewShipment = !mtoShipment?.id;

  const handleBack = () => {
    if (isNewShipment) {
      navigate(generatePath(customerRoutes.SHIPMENT_SELECT_TYPE_PATH, { moveId }));
    } else {
      navigate(generatePath(customerRoutes.MOVE_HOME_PATH, { moveId }));
    }
  };

  const onShipmentSaveSuccess = (response) => {
    // Update submitting state
    setShowBoatConfirmationModal(false);
    setIsSubmitting(false);
    const baseMtoShipment = mtoShipment?.id ? mtoShipment : response;
    const data = {
      ...baseMtoShipment,
      boatShipment: response?.boatShipment,
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
      generatePath(customerRoutes.SHIPMENT_BOAT_LOCATION_INFO, {
        moveId,
        mtoShipmentId: response.id,
      }),
    );
  };

  const closeBoatConfirmationModal = () => {
    setShowBoatConfirmationModal(false);
  };

  const redirectShipment = () => {
    setTimeout(() => {
      scrollToTop();
      const createShipmentPath = generatePath(customerRoutes.SHIPMENT_CREATE_PATH, { moveId });
      navigate(`${createShipmentPath}?type=${SHIPMENT_TYPES.HHG}`, {
        state: {
          mtoShipment,
        },
      });
    }, 100);
  };

  const handleConfirmationDeleteAndRedirect = () => {
    if (isDeleting || isSubmitting) return;
    setIsDeleting(true);

    deleteMTOShipment(mtoShipment?.id)
      .then(() => {
        getAllMoves(serviceMember.id).then((res) => {
          updateAllMoves(res);
        });
        redirectShipment();
      })
      .catch(() => {
        const errorMsg = 'There was an error attempting to delete your shipment.';
        setErrorMessage(errorMsg);
      })
      .finally(() => {
        setIsDeleting(false);
        setShowBoatConfirmationModal(false);
      });
  };

  const handleConfirmationRedirect = () => {
    setShowBoatConfirmationModal(false);
    setIsSubmitting(false);
    redirectShipment();
  };

  // Submit as a Boat shipment
  const handleConfirmationContinue = async () => {
    const values = submitValues;
    setIsSubmitting(true);
    setErrorMessage(null);

    const mtoShipmentType =
      boatShipmentObj?.type === boatShipmentTypes.TOW_AWAY
        ? SHIPMENT_TYPES.BOAT_TOW_AWAY
        : SHIPMENT_TYPES.BOAT_HAUL_AWAY;

    const createOrUpdateShipment = {
      moveTaskOrderID: moveId,
      shipmentType: mtoShipmentType,
      boatShipment: { ...boatShipmentObj },
      customerRemarks: values.customerRemarks,
    };

    if (isNewShipment) {
      createMTOShipment(createOrUpdateShipment)
        .then((shipmentResponse) => {
          onShipmentSaveSuccess(shipmentResponse);
        })
        .catch((e) => {
          const { response } = e;
          let errorMsg = 'There was an error attempting to create your shipment.';
          if (response?.body?.invalidFields) {
            const keys = Object.keys(response?.body?.invalidFields);
            const firstError = response?.body?.invalidFields[keys[0]][0];
            errorMsg = firstError;
          }
          setShowBoatConfirmationModal(false);
          setIsSubmitting(false);
          setErrorMessage(errorMsg);
        });
    } else {
      createOrUpdateShipment.id = mtoShipment.id;
      createOrUpdateShipment.boatShipment.id = mtoShipment.boatShipment?.id;

      patchMTOShipment(mtoShipment.id, createOrUpdateShipment, mtoShipment.eTag)
        .then((shipmentResponse) => {
          onShipmentSaveSuccess(shipmentResponse);
        })
        .catch((e) => {
          const { response } = e;
          let errorMsg = 'There was an error attempting to update your shipment.';
          if (response?.body?.invalidFields) {
            const keys = Object.keys(response?.body?.invalidFields);
            const firstError = response?.body?.invalidFields[keys[0]][0];
            errorMsg = firstError;
          }
          setErrorMessage(errorMsg);
          setShowBoatConfirmationModal(false);
          setIsSubmitting(false);
        });
    }
  };

  // open confirmation modal to validate boat shipment
  const handleSubmit = async (values) => {
    const totalLengthInInches = toTotalInches(values.lengthFeet, values.lengthInches);
    const totalWidthInInches = toTotalInches(values.widthFeet, values.widthInches);
    const totalHeightInInches = toTotalInches(values.heightFeet, values.heightInches);
    const hasTrailerBool = values.hasTrailer === 'true';
    const isRoadworthyBool = values.isRoadworthy && hasTrailerBool ? values.isRoadworthy === 'true' : null;
    const boatShipmentType = isRoadworthyBool ? boatShipmentTypes.TOW_AWAY : boatShipmentTypes.HAUL_AWAY;

    const boatShipment = {
      type: boatShipmentType,
      year: Number(values.year),
      make: values.make,
      model: values.model,
      lengthInInches: totalLengthInInches,
      widthInInches: totalWidthInInches,
      heightInInches: totalHeightInInches,
      hasTrailer: values.hasTrailer === 'true',
      isRoadworthy: values.hasTrailer === 'true' ? isRoadworthyBool : null,
    };
    setBoatShipmentObj(boatShipment);
    if (totalLengthInInches <= 168 && totalWidthInInches <= 82 && totalHeightInInches <= 77) {
      setIsDimensionsMeetReq(false);
    } else {
      setIsDimensionsMeetReq(true);
    }
    setSubmitValues(values);
    setShowBoatConfirmationModal(true);
  };

  return (
    <>
      <div className={pageStyles.ppmPageStyle}>
        <NotificationScrollToTop dependency={errorMessage} />
        <GridContainer>
          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }}>
              <ShipmentTag shipmentType={SHIPMENT_OPTIONS.BOAT} shipmentNumber={shipmentNumber} />
              <h1>Boat details and measurements</h1>
              {errorMessage && (
                <Alert headingLevel="h4" slim type="error">
                  {errorMessage}
                </Alert>
              )}
              <BoatShipmentForm
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
      <BoatShipmentConfirmationModal
        isDimensionsMeetReq={isDimensionsMeetReq}
        boatShipmentType={boatShipmentObj?.type}
        isOpen={showBoatConfirmationModal}
        closeModal={closeBoatConfirmationModal}
        handleConfirmationContinue={handleConfirmationContinue}
        handleConfirmationRedirect={handleConfirmationRedirect}
        handleConfirmationDeleteAndRedirect={handleConfirmationDeleteAndRedirect}
        isSubmitting={isSubmitting}
        isEditPage={isEditPage}
      />
    </>
  );
};

BoatShipmentCreate.propTypes = {
  mtoShipment: ShipmentShape,
  serviceMember: ServiceMemberShape.isRequired,
  destinationDutyLocation: DutyLocationShape.isRequired,
  move: MoveShape,
};

BoatShipmentCreate.defaultProps = {
  move: {},
  mtoShipment: {},
};

export default BoatShipmentCreate;
