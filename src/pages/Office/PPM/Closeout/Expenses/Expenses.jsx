import React, { useEffect, useState } from 'react';
import { generatePath, useNavigate, useParams } from 'react-router-dom';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';
import { useQueryClient, useMutation } from '@tanstack/react-query';

import ppmPageStyles from 'pages/Office/PPM/PPM.module.scss';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { shipmentTypes } from 'constants/shipments';
import ExpenseForm from 'components/Shared/PPM/Closeout/ExpenseForm/ExpenseForm';
import { servicesCounselingRoutes } from 'constants/routes';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import {
  createMovingExpense,
  patchExpense,
  createUploadForPPMDocument,
  deleteUploadForDocument,
} from 'services/ghcApi';
import { formatDateForSwagger } from 'shared/dates';
import { convertDollarsToCents } from 'shared/utils';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { usePPMShipmentAndDocsOnlyQueries } from 'hooks/queries';
import { DOCUMENTS } from 'constants/queryKeys';
import { APP_NAME } from 'constants/apps';

const Expenses = () => {
  const [errorMessage, setErrorMessage] = useState(null);
  const [isSubmitted, setIsSubmitted] = useState(false);
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { moveCode, shipmentId, expenseId } = useParams();

  const { mtoShipment, documents, isError } = usePPMShipmentAndDocsOnlyQueries(shipmentId);
  const appName = APP_NAME.OFFICE;
  const ppmShipment = mtoShipment?.ppmShipment ?? [];
  const expenses = documents?.MovingExpenses ?? [];
  const { ppmType } = ppmShipment;

  const currentExpense = expenses?.find((item) => item.id === expenseId) ?? null;
  const currentIndex = Array.isArray(expenses) ? expenses.findIndex((ele) => ele.id === expenseId) : -1;

  const reviewPath = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_REVIEW_PATH, { moveCode, shipmentId });

  const { mutate: mutateCreateMovingExpense } = useMutation(createMovingExpense, {
    onSuccess: (createdMovingExpense) => {
      queryClient.invalidateQueries([DOCUMENTS, shipmentId]);
      const path = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_EXPENSES_EDIT_PATH, {
        moveCode,
        shipmentId,
        expenseId: createdMovingExpense?.id,
      });
      navigate(path, { replace: true });
    },
    onError: () => {
      setErrorMessage(`Failed to create trip record`);
    },
  });

  const { mutate: mutatePatchMovingExpense } = useMutation(patchExpense, {
    onSuccess: async () => {
      await queryClient.invalidateQueries([DOCUMENTS, shipmentId]);
      navigate(reviewPath);
    },
    onError: () => {
      setIsSubmitted(false);
      setErrorMessage('Failed to save updated trip record');
    },
  });

  useEffect(() => {
    if (!expenseId) {
      mutateCreateMovingExpense(ppmShipment?.id);
    }
  }, [mutateCreateMovingExpense, ppmShipment?.id, expenseId]);

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

    createUploadForPPMDocument(ppmShipment?.id, documentId, newFile, false)
      .then((upload) => {
        documents?.MovingExpenses[currentIndex][fieldName]?.uploads.push(upload);
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
    deleteUploadForDocument(uploadId, null, ppmShipment?.id)
      .then(() => {
        const filteredUploads = documents?.MovingExpenses[currentIndex][fieldName]?.uploads.filter(
          (upload) => upload.id !== uploadId,
        );
        documents.MovingExpenses[currentIndex][fieldName].uploads = filteredUploads;

        setFieldValue(fieldName, filteredUploads, true);
        setFieldTouched(fieldName, true, true);
      })
      .catch(() => {
        setErrorMessage('Failed to delete the file upload');
      });
  };

  const handleBack = () => {
    navigate(reviewPath);
  };

  const handleSubmit = async (values) => {
    if (isSubmitted) return;

    setIsSubmitted(true);
    setErrorMessage(null);
    const payload = {
      ppmShipmentId: ppmShipment?.id,
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

    mutatePatchMovingExpense({
      ppmShipmentId: currentExpense.ppmShipmentId,
      movingExpenseId: currentExpense.id,
      payload,
      eTag: currentExpense.eTag,
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

  if (isError) {
    return <SomethingWentWrong />;
  }

  if (!mtoShipment || !currentExpense) {
    return renderError() || <LoadingPlaceholder />;
  }

  return (
    <div className={ppmPageStyles.tabContent}>
      <div className={ppmPageStyles.container}>
        <NotificationScrollToTop dependency={errorMessage} />
        <GridContainer className={ppmPageStyles.gridContainer}>
          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }}>
              <div className={ppmPageStyles.closeoutPageWrapper}>
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
              </div>
            </Grid>
          </Grid>
        </GridContainer>
      </div>
    </div>
  );
};

export default Expenses;
