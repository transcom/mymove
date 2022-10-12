import React, { useEffect, useState } from 'react';
import { generatePath, useHistory, useParams } from 'react-router-dom';
import { useSelector, useDispatch } from 'react-redux';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './Expenses.module.scss';

import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { shipmentTypes } from 'constants/shipments';
import ExpenseForm from 'components/Customer/PPM/Closeout/ExpenseForm/ExpenseForm';
import { selectExpenseAndIndexById, selectMTOShipmentById } from 'store/entities/selectors';
import { customerRoutes, generalRoutes } from 'constants/routes';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { createUploadForDocument, createMovingExpense, deleteUpload, patchMovingExpense } from 'services/internalApi';
import { updateMTOShipment } from 'store/entities/actions';
import { formatDateForSwagger } from 'shared/dates';

const Expenses = () => {
  const [errorMessage, setErrorMessage] = useState(null);

  const dispatch = useDispatch();
  const history = useHistory();
  const { moveId, mtoShipmentId, expenseId } = useParams();

  const mtoShipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));
  const { expense: currentExpense, index: currentIndex } = useSelector((state) =>
    selectExpenseAndIndexById(state, mtoShipmentId, expenseId),
  );

  useEffect(() => {
    if (!expenseId) {
      createMovingExpense(mtoShipment?.ppmShipment?.id)
        .then((resp) => {
          if (mtoShipment?.ppmShipment?.movingExpenses) {
            mtoShipment.ppmShipment.movingExpenses.push(resp);
          } else {
            mtoShipment.ppmShipment.movingExpenses = [resp];
          }
          history.replace(
            generatePath(customerRoutes.SHIPMENT_PPM_EXPENSES_EDIT_PATH, {
              moveId,
              mtoShipmentId,
              expenseId: resp.id,
            }),
          );
          dispatch(updateMTOShipment(mtoShipment));
        })
        .catch(() => {
          setErrorMessage('Failed to create trip record');
        });
    }
  }, [expenseId, moveId, mtoShipmentId, history, dispatch, mtoShipment]);

  const handleCreateUpload = async (fieldName, file, setFieldTouched) => {
    const documentId = currentExpense[`${fieldName}Id`];

    createUploadForDocument(file, documentId)
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
    deleteUpload(uploadId)
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
    history.push(generalRoutes.HOME_PATH);
  };

  const handleSubmit = async (values, { setSubmitting }) => {
    setErrorMessage(null);
    const payload = {
      ppmShipmentId: mtoShipment.ppmShipment.id,
      movingExpenseType: values.expenseType,
      amount: values.amount * 100,
      description: values.description,
      missingReceipt: values.missingReceipt,
      paidWithGTCC: values.paidWithGTCC === 'true',
      SITEndDate: formatDateForSwagger(values.sitEndDate),
      SITStartDate: formatDateForSwagger(values.sitStartDate),
    };

    patchMovingExpense(mtoShipment?.ppmShipment?.id, currentExpense.id, payload, currentExpense.eTag)
      .then((resp) => {
        setSubmitting(false);
        mtoShipment.ppmShipment.movingExpenses[currentIndex] = resp;
        history.push(generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, { moveId, mtoShipmentId }));
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
