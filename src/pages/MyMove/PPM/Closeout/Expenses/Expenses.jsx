import React, { useEffect, useState } from 'react';
import { generatePath, useNavigate, useParams } from 'react-router-dom';
import { useSelector, useDispatch } from 'react-redux';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './Expenses.module.scss';

import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { shipmentTypes } from 'constants/shipments';
import ExpenseForm from 'components/Shared/PPM/Closeout/ExpenseForm/ExpenseForm';
import { selectExpenseAndIndexById, selectMTOShipmentById } from 'store/entities/selectors';
import { customerRoutes } from 'constants/routes';
import { CUSTOMER_ERROR_MESSAGES } from 'constants/errorMessages';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import {
  createUploadForPPMDocument,
  createMovingExpense,
  deleteUpload,
  patchMovingExpense,
} from 'services/internalApi';
import { updateMTOShipment } from 'store/entities/actions';
import { formatDateForSwagger } from 'shared/dates';
import { convertDollarsToCents } from 'shared/utils';
import appendTimestampToFilename from 'utils/fileUpload';
import { APP_NAME } from 'constants/apps';

const Expenses = () => {
  const [errorMessage, setErrorMessage] = useState(null);

  const dispatch = useDispatch();
  const navigate = useNavigate();
  const { moveId, mtoShipmentId, expenseId } = useParams();

  const mtoShipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));
  const { expense: currentExpense, index: currentIndex } = useSelector((state) =>
    selectExpenseAndIndexById(state, mtoShipmentId, expenseId),
  );

  const ppmShipment = mtoShipment?.ppmShipment || {};
  const { ppmType } = ppmShipment;

  const appName = APP_NAME.MYMOVE;

  useEffect(() => {
    if (!expenseId) {
      createMovingExpense(mtoShipment?.ppmShipment?.id)
        .then((resp) => {
          if (mtoShipment?.ppmShipment?.movingExpenses) {
            mtoShipment.ppmShipment.movingExpenses.push(resp);
          } else {
            mtoShipment.ppmShipment.movingExpenses = [resp];
          }
          const path = generatePath(customerRoutes.SHIPMENT_PPM_EXPENSES_EDIT_PATH, {
            moveId,
            mtoShipmentId,
            expenseId: resp.id,
          });
          navigate(path, { replace: true });
          dispatch(updateMTOShipment(mtoShipment));
        })
        .catch(() => {
          setErrorMessage('Failed to create trip record');
        });
    }
  }, [expenseId, moveId, mtoShipmentId, navigate, dispatch, mtoShipment]);

  const handleCreateUpload = async (fieldName, file, setFieldTouched) => {
    const documentId = currentExpense[`${fieldName}Id`];
    createUploadForPPMDocument(mtoShipment.ppmShipment.id, documentId, appendTimestampToFilename(file), false)
      .then((upload) => {
        mtoShipment.ppmShipment.movingExpenses[currentIndex][fieldName].uploads.push(upload);
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
        const filteredUploads = mtoShipment.ppmShipment.movingExpenses[currentIndex][fieldName].uploads.filter(
          (upload) => upload.id !== uploadId,
        );
        mtoShipment.ppmShipment.movingExpenses[currentIndex][fieldName].uploads = filteredUploads;

        setFieldValue(fieldName, filteredUploads, true);
        setFieldTouched(fieldName, true, true);
        dispatch(updateMTOShipment(mtoShipment));
      })
      .catch(() => {
        setErrorMessage('Failed to delete the file upload');
      });
  };

  const handleBack = () => {
    const path = generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, {
      moveId,
      mtoShipmentId,
    });
    navigate(path);
  };

  const handleErrorMessage = (error) => {
    if (error?.response?.status === 412) {
      setErrorMessage(CUSTOMER_ERROR_MESSAGES.PRECONDITION_FAILED);
    } else {
      setErrorMessage('Failed to save updated trip record');
    }
  };

  const handleSubmit = async (values, { setSubmitting }) => {
    setErrorMessage(null);
    const payload = {
      ppmShipmentId: mtoShipment.ppmShipment.id,
      movingExpenseType: values.expenseType,
      amount: convertDollarsToCents(values.amount),
      description: values.description,
      missingReceipt: values.missingReceipt,
      paidWithGTCC: values.paidWithGTCC === 'true',
      SITEndDate: formatDateForSwagger(values.sitEndDate),
      SITStartDate: formatDateForSwagger(values.sitStartDate),
      WeightStored: parseInt(values.sitWeight, 10),
      SITLocation: values.sitLocation,
      weightShipped: parseInt(values.weightShipped, 10),
      trackingNumber: values.trackingNumber,
      isProGear: values.isProGear === 'true',
      ...(values.isProGear === 'true' && {
        proGearBelongsToSelf: values.proGearBelongsToSelf === 'true',
        proGearDescription: values.proGearDescription,
      }),
    };

    patchMovingExpense(mtoShipment?.ppmShipment?.id, currentExpense.id, payload, currentExpense.eTag)
      .then((resp) => {
        setSubmitting(false);
        mtoShipment.ppmShipment.movingExpenses[currentIndex] = resp;
        navigate(generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, { moveId, mtoShipmentId }));
        dispatch(updateMTOShipment(mtoShipment));
      })
      .catch((error) => {
        setSubmitting(false);
        handleErrorMessage(error);
      });
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

  if (!mtoShipment || !currentExpense) {
    return renderError() || <LoadingPlaceholder />;
  }

  return (
    <div className={classnames(styles.Expenses, ppmPageStyles.ppmPageStyle)}>
      <NotificationScrollToTop dependency={errorMessage} />
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <ShipmentTag shipmentType={shipmentTypes.PPM} />
            <h1>Expenses</h1>
            {renderError()}
            <ExpenseForm
              ppmType={ppmType}
              expense={currentExpense}
              receiptNumber={currentIndex + 1}
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

export default Expenses;
