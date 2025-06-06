import React, { useEffect, useState } from 'react';
import { generatePath, useNavigate, useParams } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';

import NotificationScrollToTop from 'components/NotificationScrollToTop';
import {
  selectMTOShipmentById,
  selectGunSafeWeightTicketAndIndexById,
  selectServiceMemberFromLoggedInUser,
  selectProGearEntitlements,
} from 'store/entities/selectors';
import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { shipmentTypes } from 'constants/shipments';
import { customerRoutes } from 'constants/routes';
import {
  createUploadForPPMDocument,
  createGunSafeWeightTicket,
  deleteUpload,
  patchGunSafeWeightTicket,
  getMTOShipmentsForMove,
  getAllMoves,
} from 'services/internalApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import GunSafeForm from 'components/Shared/PPM/Closeout/GunSafeForm/GunSafeForm';
import { updateAllMoves, updateMTOShipment } from 'store/entities/actions';
import { CUSTOMER_ERROR_MESSAGES } from 'constants/errorMessages';
import { APP_NAME } from 'constants/apps';
import appendTimestampToFilename from 'utils/fileUpload';

const GunSafe = () => {
  const dispatch = useDispatch();
  const navigate = useNavigate();

  const serviceMember = useSelector((state) => selectServiceMemberFromLoggedInUser(state));
  const serviceMemberId = serviceMember.id;

  const entitlements = useSelector((state) => selectProGearEntitlements(state));

  const appName = APP_NAME.MYMOVE;
  const { moveId, mtoShipmentId, gunSafeId } = useParams();

  const handleBack = () => {
    const path = generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, {
      moveId,
      mtoShipmentId,
    });
    navigate(path);
  };
  const [errorMessage, setErrorMessage] = useState(null);

  const mtoShipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));
  const { gunSafeWeightTicket: currentGunSafeWeightTicket, index: currentIndex } = useSelector((state) =>
    selectGunSafeWeightTicketAndIndexById(state, mtoShipmentId, gunSafeId),
  );

  useEffect(() => {
    if (!gunSafeId) {
      createGunSafeWeightTicket(mtoShipment?.ppmShipment?.id)
        .then((resp) => {
          if (mtoShipment?.ppmShipment?.gunSafeWeightTickets) {
            mtoShipment.ppmShipment.gunSafeWeightTickets.push(resp);
          } else {
            mtoShipment.ppmShipment.gunSafeWeightTickets = [resp];
          }
          // Update the URL so the back button would work and not create a new weight ticket or on
          // refresh either.
          navigate(
            generatePath(customerRoutes.SHIPMENT_PPM_GUN_SAFE_EDIT_PATH, {
              moveId,
              mtoShipmentId,
              gunSafeId: resp.id,
            }),
            { replace: true },
          );
          dispatch(updateMTOShipment(mtoShipment));
        })
        .catch(() => {
          setErrorMessage('Failed to create trip record');
        });
    }
  }, [gunSafeId, moveId, mtoShipmentId, navigate, dispatch, mtoShipment]);

  useEffect(() => {
    const moves = getAllMoves(serviceMemberId);
    dispatch(updateAllMoves(moves));
  }, [gunSafeId, moveId, mtoShipmentId, navigate, dispatch, mtoShipment, serviceMemberId]);

  const handleErrorMessage = (error) => {
    if (error?.response?.status === 412) {
      setErrorMessage(CUSTOMER_ERROR_MESSAGES.PRECONDITION_FAILED);
    } else {
      setErrorMessage('Failed to fetch shipment information');
    }
  };

  const handleCreateUpload = async (fieldName, file, setFieldTouched) => {
    const { documentId } = currentGunSafeWeightTicket;

    createUploadForPPMDocument(mtoShipment.ppmShipment.id, documentId, appendTimestampToFilename(file), false)
      .then((upload) => {
        mtoShipment.ppmShipment.gunSafeWeightTickets[currentIndex].document.uploads.push(upload);
        dispatch(updateMTOShipment(mtoShipment));
        setFieldTouched(fieldName, true);
        return upload;
      })
      .catch(() => {
        setErrorMessage('Failed to save the file upload');
      });
  };

  const handleUploadComplete = (err) => {
    if (err) {
      setErrorMessage('Encountered error when completing file upload');
    }
  };

  const handleUploadDelete = (uploadId, fieldName, setFieldTouched, setFieldValue) => {
    deleteUpload(uploadId, null, mtoShipment?.ppmShipment?.id)
      .then(() => {
        const filteredUploads = mtoShipment.ppmShipment.gunSafeWeightTickets[currentIndex].document.uploads.filter(
          (upload) => upload.id !== uploadId,
        );
        mtoShipment.ppmShipment.gunSafeWeightTickets[currentIndex].document.uploads = filteredUploads;
        setFieldValue(fieldName, filteredUploads, true);
        setFieldTouched(fieldName, true, true);
        dispatch(updateMTOShipment(mtoShipment));
      })
      .catch(() => {
        setErrorMessage('Failed to delete the file upload');
      });
  };

  const updateGunSafeWeightTicket = (values) => {
    const hasWeightTickets = !values.missingWeightTicket;
    const payload = {
      ppmShipmentId: mtoShipment.ppmShipment.id,
      gunSafeWeightTicketId: currentGunSafeWeightTicket.id,
      description: values.description,
      weight: parseInt(values.weight, 10),
      hasWeightTickets,
    };

    patchGunSafeWeightTicket(
      mtoShipment?.ppmShipment?.id,
      currentGunSafeWeightTicket.id,
      payload,
      currentGunSafeWeightTicket.eTag,
    )
      .then((resp) => {
        mtoShipment.ppmShipment.gunSafeWeightTickets[currentIndex] = resp;
        getMTOShipmentsForMove(moveId)
          .then((response) => {
            dispatch(updateMTOShipment(response.mtoShipments[mtoShipmentId]));
            mtoShipment.eTag = response.mtoShipments[mtoShipmentId].eTag;
            navigate(generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, { moveId, mtoShipmentId }));
          })
          .catch(() => {
            setErrorMessage('Failed to fetch shipment information');
          });
      })
      .catch((error) => {
        handleErrorMessage(error);
      });
  };

  const handleSubmit = async (values, { setSubmitting, setErrors }) => {
    setErrorMessage(null);
    setErrors({});
    setSubmitting(false);
    updateGunSafeWeightTicket(values);
  };

  const renderError = () => {
    if (!errorMessage) {
      return null;
    }

    return (
      <Alert slim type="error">
        {errorMessage}
      </Alert>
    );
  };

  if (!mtoShipment || !currentGunSafeWeightTicket) {
    return renderError() || <LoadingPlaceholder />;
  }
  return (
    <div className={ppmPageStyles.ppmPageStyle}>
      <NotificationScrollToTop dependency={errorMessage} />
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <ShipmentTag shipmentType={shipmentTypes.PPM} />
            <h1>Gun Safe</h1>
            {renderError()}
            <GunSafeForm
              gunSafe={currentGunSafeWeightTicket}
              setNumber={currentIndex + 1}
              entitlements={entitlements}
              onBack={handleBack}
              onSubmit={handleSubmit}
              onCreateUpload={handleCreateUpload}
              onUploadComplete={handleUploadComplete}
              onUploadDelete={handleUploadDelete}
              appName={appName}
            />
          </Grid>
        </Grid>
      </GridContainer>
    </div>
  );
};

export default GunSafe;
