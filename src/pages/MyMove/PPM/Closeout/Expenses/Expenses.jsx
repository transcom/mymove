import React, { useEffect, useState } from 'react';
import { generatePath, useNavigate, useParams } from 'react-router-dom';
import { useSelector, useDispatch } from 'react-redux';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';
import classnames from 'classnames';

import { isBooleanFlagEnabled } from '../../../../../utils/featureFlags';

import styles from './Expenses.module.scss';

import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { shipmentTypes } from 'constants/shipments';
import ExpenseForm from 'components/Customer/PPM/Closeout/ExpenseForm/ExpenseForm';
import { selectExpenseAndIndexById, selectMTOShipmentById } from 'store/entities/selectors';
import { customerRoutes } from 'constants/routes';
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

const Expenses = () => {
  const [errorMessage, setErrorMessage] = useState(null);
  const [multiMove, setMultiMove] = useState(false);

  const dispatch = useDispatch();
  const navigate = useNavigate();
  const { moveId, mtoShipmentId, expenseId } = useParams();

  const mtoShipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));
  const { expense: currentExpense, index: currentIndex } = useSelector((state) =>
    selectExpenseAndIndexById(state, mtoShipmentId, expenseId),
  );

  useEffect(() => {
    isBooleanFlagEnabled('multi_move').then((enabled) => {
      setMultiMove(enabled);
    });
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
    if (multiMove) {
      navigate(generatePath(customerRoutes.MOVE_HOME_PATH, { moveId }));
    } else {
      navigate(customerRoutes.MOVE_HOME_PAGE);
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
    };

    patchMovingExpense(mtoShipment?.ppmShipment?.id, currentExpense.id, payload, currentExpense.eTag)
      .then((resp) => {
        setSubmitting(false);
        mtoShipment.ppmShipment.movingExpenses[currentIndex] = resp;
        navigate(generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, { moveId, mtoShipmentId }));
        dispatch(updateMTOShipment(mtoShipment));
      })
      .catch(() => {
        setSubmitting(false);
        setErrorMessage('Failed to save updated trip record');
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
            <div className={styles.introSection}>
              <p>
                Document your qualified expenses by uploading receipts. They should include a description of the item,
                the price you paid, the date of purchase, and the business name. All documents must be legible and
                unaltered.
              </p>
              <p>
                Your finance office will make the final decision about which expenses are deductible or reimbursable.
              </p>
              <p>Upload one receipt at a time. Please do not put multiple receipts in one image.</p>
            </div>
            <ExpenseForm
              expense={currentExpense}
              receiptNumber={currentIndex + 1}
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

export default Expenses;
