import React, { useState } from 'react';
import { GridContainer, Grid, Button } from '@trussworks/react-uswds';
import { Link, useParams, generatePath } from 'react-router-dom';
import { useSelector, useDispatch } from 'react-redux';
import classnames from 'classnames';

import styles from './Review.module.scss';

import Alert from 'shared/Alert';
import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { shipmentTypes } from 'constants/shipments';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { customerRoutes, generalRoutes } from 'constants/routes';
import { selectMTOShipmentById } from 'store/entities/selectors';
import ReviewItems from 'components/Customer/PPM/Closeout/ReviewItems/ReviewItems';
import {
  calculateTotalMovingExpensesAmount,
  formatAboutYourPPMItem,
  formatExpenseItems,
  formatProGearItems,
  formatWeightTicketItems,
} from 'utils/ppmCloseout';
import {
  calculateTotalNetWeightForProGearWeightTickets,
  getTotalNetWeightForWeightTickets,
} from 'utils/shipmentWeights';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { formatCents, formatWeight } from 'utils/formatters';
import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';
import {
  deleteWeightTicket,
  deleteProGearWeightTicket,
  deleteMovingExpense,
  getMTOShipmentsForMove,
} from 'services/internalApi';
import ppmStyles from 'components/Customer/PPM/PPM.module.scss';
import { hasCompletedAllWeightTickets, hasCompletedAllExpenses, hasCompletedAllProGear } from 'utils/shipments';
import { updateMTOShipment } from 'store/entities/actions';

const ReviewDeleteCloseoutItemModal = ({ onClose, onSubmit, itemToDelete }) => {
  const deleteDetailMessage = <p>You are about to delete {itemToDelete.itemNumber}. This cannot be undone.</p>;
  return (
    <div>
      <Overlay />
      <ModalContainer>
        <Modal>
          <ModalClose handleClick={() => onClose(false)} />
          <ModalTitle>
            <h3>Delete this?</h3>
          </ModalTitle>
          {deleteDetailMessage}
          <ModalActions>
            <Button
              className="usa-button--destructive"
              type="submit"
              onClick={() => onSubmit(itemToDelete.itemType, itemToDelete.itemId, itemToDelete.itemNumber)}
            >
              Yes, Delete
            </Button>
            <Button type="button" onClick={() => onClose(false)} data-testid="modalBackButton" secondary>
              No, Keep It
            </Button>
          </ModalActions>
        </Modal>
      </ModalContainer>
    </div>
  );
};

function deleteLineItem(ppmShipmentId, itemType, itemId) {
  if (itemType === 'weightTicket') {
    return deleteWeightTicket(ppmShipmentId, itemId);
  }
  if (itemType === 'proGear') {
    return deleteProGearWeightTicket(ppmShipmentId, itemId);
  }
  if (itemType === 'expense') {
    return deleteMovingExpense(ppmShipmentId, itemId);
  }
  return Promise.reject();
}

const Review = () => {
  const [isDeleteModalVisible, setIsDeleteModalVisible] = useState(false);
  const [itemToDelete, setItemToDelete] = useState();
  const [alert, setAlert] = useState(null);
  const { moveId, mtoShipmentId } = useParams();
  const mtoShipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));

  const weightTickets = mtoShipment?.ppmShipment?.weightTickets;
  const proGear = mtoShipment?.ppmShipment?.proGearWeightTickets;
  const expenses = mtoShipment?.ppmShipment?.movingExpenses;
  const dispatch = useDispatch();

  if (!mtoShipment) {
    return <LoadingPlaceholder />;
  }

  const weightTicketCreatePath = generatePath(customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_PATH, {
    moveId,
    mtoShipmentId,
  });
  const proGearCreatePath = generatePath(customerRoutes.SHIPMENT_PPM_PRO_GEAR_PATH, { moveId, mtoShipmentId });
  const expensesCreatePath = generatePath(customerRoutes.SHIPMENT_PPM_EXPENSES_PATH, { moveId, mtoShipmentId });
  const completePath = generatePath(customerRoutes.SHIPMENT_PPM_COMPLETE_PATH, { moveId, mtoShipmentId });

  const handleDelete = (itemType, itemId, itemETag, itemNumber) => {
    setItemToDelete(() => ({
      itemType,
      itemId,
      itemETag,
      itemNumber,
    }));
    setIsDeleteModalVisible(true);
  };

  const onDeleteSubmit = (itemType, itemId, itemNumber) => {
    deleteLineItem(mtoShipment.ppmShipment.id, itemType, itemId)
      .then(() => {
        setIsDeleteModalVisible(false);
        getMTOShipmentsForMove(mtoShipment.moveTaskOrderID).then((moveResponse) =>
          dispatch(updateMTOShipment(moveResponse.mtoShipments[mtoShipment.id])),
        );
      })
      .then(() => setAlert({ type: 'success', message: `${itemNumber} successfully deleted.` }))
      .catch(() => {
        setIsDeleteModalVisible(false);
        setAlert({ type: 'error', message: `Something went wrong deleting ${itemNumber}. Please try again.` });
      });
  };

  const aboutYourPPM = formatAboutYourPPMItem(mtoShipment?.ppmShipment, customerRoutes.SHIPMENT_PPM_ABOUT_PATH, {
    moveId,
    mtoShipmentId,
  });

  const weightTicketContents = formatWeightTicketItems(
    weightTickets,
    customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_EDIT_PATH,
    { moveId, mtoShipmentId },
    handleDelete,
  );

  const weightTicketsTotal = getTotalNetWeightForWeightTickets(weightTickets);

  const canAdvance =
    hasCompletedAllWeightTickets(weightTickets) && hasCompletedAllExpenses(expenses) && hasCompletedAllProGear(proGear);

  const proGearContents = formatProGearItems(
    proGear,
    customerRoutes.SHIPMENT_PPM_PRO_GEAR_EDIT_PATH,
    { moveId, mtoShipmentId },
    handleDelete,
  );

  const proGearTotal = calculateTotalNetWeightForProGearWeightTickets(proGear);

  const expenseContents = formatExpenseItems(
    expenses,
    customerRoutes.SHIPMENT_PPM_EXPENSES_EDIT_PATH,
    {
      moveId,
      mtoShipmentId,
    },
    handleDelete,
  );

  const expensesTotal = calculateTotalMovingExpensesAmount(expenses);

  return (
    <div className={classnames(ppmPageStyles.ppmPageStyle, styles.PPMReview)}>
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            {isDeleteModalVisible && (
              <ReviewDeleteCloseoutItemModal
                onSubmit={onDeleteSubmit}
                onClose={setIsDeleteModalVisible}
                itemToDelete={itemToDelete}
              />
            )}
            {alert && (
              <>
                <Alert type={alert.type}>{alert.message}</Alert>
                <br />
              </>
            )}
            {!canAdvance && (
              <>
                <Alert type="error">
                  There are items below that are missing required information. Please select “Edit” to enter all
                  required information or “Delete” to remove the item.
                </Alert>
                <br />
              </>
            )}
            <ShipmentTag shipmentType={shipmentTypes.PPM} />
            <h1>Review</h1>
            <SectionWrapper className={styles.aboutSection} data-testid="aboutYourPPM">
              <ReviewItems heading={<h2>About Your PPM</h2>} contents={aboutYourPPM} />
            </SectionWrapper>
            <SectionWrapper>
              <h2>Documents</h2>
              <ReviewItems
                className={classnames(styles.reviewItems, 'reviewWeightTickets')}
                heading={
                  <>
                    <h3>Weight moved</h3>
                    <span>({formatWeight(weightTicketsTotal)})</span>
                  </>
                }
                contents={weightTicketContents}
                renderAddButton={() => (
                  <Link className="usa-button usa-button--secondary" to={weightTicketCreatePath}>
                    Add More Weight
                  </Link>
                )}
                emptyMessage="No weight moved documented. At least one trip is required to continue."
              />
              <ReviewItems
                className={classnames(styles.reviewItems, 'progearSection')}
                heading={
                  <>
                    <h3>Pro-gear</h3>
                    <span>({formatWeight(proGearTotal)})</span>
                  </>
                }
                contents={proGearContents}
                renderAddButton={() => (
                  <Link className="usa-button usa-button--secondary" to={proGearCreatePath}>
                    Add Pro-gear Weight
                  </Link>
                )}
                emptyMessage="No pro-gear weight documented."
              />
              <ReviewItems
                className={classnames(styles.reviewItems, 'reviewExpenses')}
                heading={
                  <>
                    <h3>Expenses</h3>
                    <span>(${expensesTotal ? formatCents(expensesTotal) : 0})</span>
                  </>
                }
                contents={expenseContents}
                renderAddButton={() => (
                  <Link className="usa-button usa-button--secondary" to={expensesCreatePath}>
                    Add Expenses
                  </Link>
                )}
                emptyMessage="No receipts uploaded."
              />
            </SectionWrapper>
            <div className={classnames(ppmStyles.buttonContainer, styles.navigationButtons)}>
              <Link
                className={classnames(ppmStyles.backButton, 'usa-button', 'usa-button--secondary')}
                to={generalRoutes.HOME_PATH}
                data-testid="reviewReturnToHomepageLink"
              >
                Return To Homepage
              </Link>
              <Link
                className={classnames(ppmStyles.saveButton, 'usa-button', {
                  'usa-button--disabled': !canAdvance,
                })}
                aria-disabled={!canAdvance}
                to={completePath}
              >
                Save & Continue
              </Link>
            </div>
          </Grid>
        </Grid>
      </GridContainer>
    </div>
  );
};

export default Review;
