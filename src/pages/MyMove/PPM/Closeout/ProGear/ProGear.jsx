import React, { useEffect, useState } from 'react';
import { generatePath, useNavigate, useParams } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';

import { isBooleanFlagEnabled } from '../../../../../utils/featureFlags';

import {
  selectMTOShipmentById,
  selectProGearWeightTicketAndIndexById,
  selectServiceMemberFromLoggedInUser,
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
  patchMTOShipment,
  getMTOShipmentsForMove,
  getAllMoves,
} from 'services/internalApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import closingPageStyles from 'pages/MyMove/PPM/Closeout/Closeout.module.scss';
import ProGearForm from 'components/Customer/PPM/Closeout/ProGearForm/ProGearForm';
import { updateAllMoves, updateMTOShipment } from 'store/entities/actions';

const ProGear = () => {
  const dispatch = useDispatch();
  const navigate = useNavigate();

  const serviceMember = useSelector((state) => selectServiceMemberFromLoggedInUser(state));
  const serviceMemberId = serviceMember.id;

  const { moveId, mtoShipmentId, proGearId } = useParams();

  const [multiMove, setMultiMove] = useState(false);
  const handleBack = () => {
    if (multiMove) {
      navigate(generatePath(customerRoutes.MOVE_HOME_PATH, { moveId }));
    } else {
      navigate(customerRoutes.MOVE_HOME_PAGE);
    }
  };
  const [errorMessage, setErrorMessage] = useState(null);

  const mtoShipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));
  const { proGearWeightTicket: currentProGearWeightTicket, index: currentIndex } = useSelector((state) =>
    selectProGearWeightTicketAndIndexById(state, mtoShipmentId, proGearId),
  );

  useEffect(() => {
    isBooleanFlagEnabled('multi_move').then((enabled) => {
      setMultiMove(enabled);
    });
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

  const updateMtoShipment = (values) => {
    const belongsToSelf = values.belongsToSelf === 'true';
    let proGear;
    let spouseProGear;
    if (belongsToSelf) {
      proGear = values.weight;
    }
    if (!belongsToSelf) {
      spouseProGear = values.weight;
    }
    const payload = {
      belongsToSelf,
      ppmShipment: {
        id: mtoShipment.ppmShipment.id,
      },
      shipmentType: mtoShipment.shipmentType,
      actualSpouseProGearWeight: parseInt(spouseProGear, 10),
      actualProGearWeight: parseInt(proGear, 10),
      shipmentLocator: values.shipmentLocator,
      eTag: mtoShipment.eTag,
    };

    patchMTOShipment(mtoShipment.id, payload, payload.eTag)
      .then((response) => {
        navigate(generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, { moveId, mtoShipmentId }));
        dispatch(updateMTOShipment(response));
      })
      .catch(() => {
        setErrorMessage('Failed to update MTO shipment due to server error.');
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
            updateMtoShipment(values);
            navigate(generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, { moveId, mtoShipmentId }));
          })
          .catch(() => {
            setErrorMessage('Failed to fetch shipment information');
          });
      })
      .catch(() => {
        setErrorMessage('Failed to save updated trip record');
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
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <ShipmentTag shipmentType={shipmentTypes.PPM} />
            <h1>Pro-gear</h1>
            {renderError()}
            <div className={closingPageStyles['closing-section']}>
              <p>
                If you moved pro-gear for yourself or your spouse as part of this PPM, document the total weight here.
                Reminder: This pro-gear should be included in your total weight moved.
              </p>
            </div>
            <ProGearForm
              proGear={currentProGearWeightTicket}
              setNumber={currentIndex + 1}
              onBack={handleBack}
              onSubmit={handleSubmit}
              onCreateUpload={handleCreateUpload}
              onUploadComplete={handleUploadComplete}
              onUploadDelete={handleUploadDelete}
            />
          </Grid>
        </Grid>
      </GridContainer>
    </div>
  );
};

export default ProGear;
