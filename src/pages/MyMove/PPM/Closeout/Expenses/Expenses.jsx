import React from 'react';
import { generatePath, useHistory, useParams } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './Expenses.module.scss';

import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import ScrollToTop from 'components/ScrollToTop';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { shipmentTypes } from 'constants/shipments';
import ExpenseForm from 'components/Customer/PPM/Closeout/ExpenseForm/ExpenseForm';
import { selectExpenseAndIndexById, selectMTOShipmentById } from 'store/entities/selectors';
import { customerRoutes, generalRoutes } from 'constants/routes';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';

const Expenses = () => {
  const { moveId, mtoShipmentId, expenseId } = useParams();
  const history = useHistory();

  const mtoShipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));
  const { expense: currentExpense, index: currentIndex } = useSelector((state) =>
    selectExpenseAndIndexById(state, mtoShipmentId, expenseId),
  );

  const handleBack = () => {
    history.push(generatePath(generalRoutes.HOME_PATH));
  };

  const handleSubmit = () => {
    // TODO: Calls update expense API endpoint
    history.push(generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, { moveId, mtoShipmentId }));
  };

  const handleCreateUpload = () => {};
  const handleUploadComplete = () => {};
  const handleUploadDelete = () => {};

  if (!mtoShipment) {
    return <LoadingPlaceholder />;
  }

  return (
    <div className={classnames(styles.Expenses, ppmPageStyles.ppmPageStyle)}>
      <ScrollToTop />
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <ShipmentTag shipmentType={shipmentTypes.PPM} />
            <h1>Expenses</h1>
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
              receiptNumber={currentIndex >= 0 ? currentIndex + 1 : undefined}
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
