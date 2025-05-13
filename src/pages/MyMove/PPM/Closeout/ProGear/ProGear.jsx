import React, { useEffect, useState } from 'react';
import { generatePath, useNavigate, useParams } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';

import NotificationScrollToTop from 'components/NotificationScrollToTop';
import {
  selectMTOShipmentById,
  selectProGearWeightTicketAndIndexById,
  selectServiceMemberFromLoggedInUser,
  selectProGearEntitlements,
} from 'store/entities/selectors';
import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { shipmentTypes } from 'constants/shipments';
import { customerRoutes } from 'constants/routes';
import {
  createUploadForPPMDocument,
  createProGearWeightTicket,
  deleteUpload,
  patchProGearWeightTicket,
  getMTOShipmentsForMove,
  getAllMoves,
} from 'services/internalApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import ProGearForm from 'components/Shared/PPM/Closeout/ProGearForm/ProGearForm';
import { updateAllMoves, updateMTOShipment } from 'store/entities/actions';
import { CUSTOMER_ERROR_MESSAGES } from 'constants/errorMessages';
import { APP_NAME } from 'constants/apps';

const ProGear = () => {
  const dispatch = useDispatch();
  const navigate = useNavigate();

  const serviceMember = useSelector((state) => selectServiceMemberFromLoggedInUser(state));
  const serviceMemberId = serviceMember.id;

  const proGearEntitlements = useSelector((state) => selectProGearEntitlements(state));

  const appName = APP_NAME.MYMOVE;
  const { moveId, mtoShipmentId, proGearId } = useParams();

  const handleBack = () => {
    const path = generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, {
      moveId,
      mtoShipmentId,
    });
    navigate(path);
  };
  const [errorMessage, setErrorMessage] = useState(null);

  const mtoShipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));
  const { proGearWeightTicket: currentProGearWeightTicket, index: currentIndex } = useSelector((state) =>
    selectProGearWeightTicketAndIndexById(state, mtoShipmentId, proGearId),
  );

  useEffect(() => {
    if (!proGearId) {
      createProGearWeightTicket(mtoShipment?.ppmShipment?.id)
        .then((resp) => {
          if (mtoShipment?.ppmShipment?.proGearWeightTickets) {
            mtoShipment.ppmShipment.proGearWeightTickets.push(resp);
          } else {
            mtoShipment.ppmShipment.proGearWeightTickets = [resp];
          }
          // Update the URL so the back button would work and not create a new weight ticket or on
          // refresh either.
          navigate(
            generatePath(customerRoutes.SHIPMENT_PPM_PRO_GEAR_EDIT_PATH, {
              moveId,
              mtoShipmentId,
              proGearId: resp.id,
            }),
            { replace: true },
          );
          dispatch(updateMTOShipment(mtoShipment));
        })
        .catch(() => {
          setErrorMessage('Failed to create trip record');
        });
    }
  }, [proGearId, moveId, mtoShipmentId, navigate, dispatch, mtoShipment]);

  useEffect(() => {
    const moves = getAllMoves(serviceMemberId);
    dispatch(updateAllMoves(moves));
  }, [proGearId, moveId, mtoShipmentId, navigate, dispatch, mtoShipment, serviceMemberId]);

  const handleErrorMessage = (error) => {
    if (error?.response?.status === 412) {
      setErrorMessage(CUSTOMER_ERROR_MESSAGES.PRECONDITION_FAILED);
    } else {
      setErrorMessage('Failed to fetch shipment information');
    }
  };

  const handleCreateUpload = async (fieldName, file, setFieldTouched) => {
    const documentId = currentProGearWeightTicket[`${fieldName}Id`];

    // Create a date-time stamp in the format "yyyymmddhh24miss"
    const now = new Date();
    const timestamp =
      now.getFullYear().toString() +
      (now.getMonth() + 1).toString().padStart(2, '0') +
      now.getDate().toString().padStart(2, '0') +
      now.getHours().toString().padStart(2, '0') +
      now.getMinutes().toString().padStart(2, '0') +
      now.getSeconds().toString().padStart(2, '0');

    // Create a new filename with the timestamp prepended
    const newFileName = `${file.name}-${timestamp}`;

    // Create and return a new File object with the new filename
    const newFile = new File([file], newFileName, { type: file.type });

    createUploadForPPMDocument(mtoShipment.ppmShipment.id, documentId, newFile, false)
      .then((upload) => {
        mtoShipment.ppmShipment.proGearWeightTickets[currentIndex][fieldName].uploads.push(upload);
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
        const filteredUploads = mtoShipment.ppmShipment.proGearWeightTickets[currentIndex][fieldName].uploads.filter(
          (upload) => upload.id !== uploadId,
        );
        mtoShipment.ppmShipment.proGearWeightTickets[currentIndex][fieldName].uploads = filteredUploads;
        setFieldValue(fieldName, filteredUploads, true);
        setFieldTouched(fieldName, true, true);
        dispatch(updateMTOShipment(mtoShipment));
      })
      .catch(() => {
        setErrorMessage('Failed to delete the file upload');
      });
  };

  const updateProGearWeightTicket = (values) => {
    const hasWeightTickets = !values.missingWeightTicket;
    const belongsToSelf = values.belongsToSelf === 'true';
    const payload = {
      ppmShipmentId: mtoShipment.ppmShipment.id,
      proGearWeightTicketId: currentProGearWeightTicket.id,
      description: values.description,
      weight: parseInt(values.weight, 10),
      belongsToSelf,
      hasWeightTickets,
    };

    patchProGearWeightTicket(
      mtoShipment?.ppmShipment?.id,
      currentProGearWeightTicket.id,
      payload,
      currentProGearWeightTicket.eTag,
    )
      .then((resp) => {
        mtoShipment.ppmShipment.proGearWeightTickets[currentIndex] = resp;
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
    updateProGearWeightTicket(values);
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

  if (!mtoShipment || !currentProGearWeightTicket) {
    return renderError() || <LoadingPlaceholder />;
  }
  return (
    <div className={ppmPageStyles.ppmPageStyle}>
      <NotificationScrollToTop dependency={errorMessage} />
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <ShipmentTag shipmentType={shipmentTypes.PPM} />
            <h1>Pro-gear</h1>
            {renderError()}
            <ProGearForm
              entitlements={proGearEntitlements}
              proGear={currentProGearWeightTicket}
              setNumber={currentIndex + 1}
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

export default ProGear;
