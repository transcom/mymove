import React, { useState } from 'react';
import { GridContainer, Grid, Button } from '@trussworks/react-uswds';
import { Link, useParams, generatePath } from 'react-router-dom';
import { useSelector } from 'react-redux';
import classnames from 'classnames';

import styles from './Review.module.scss';

import ppmPageStyles from 'pages/MyMove/PPM/PPM.module.scss';
import ScrollToTop from 'components/ScrollToTop';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { shipmentTypes } from 'constants/shipments';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { customerRoutes, generalRoutes } from 'constants/routes';
import { selectMTOShipmentById } from 'store/entities/selectors';
import ReviewItems from 'components/Customer/PPM/Closeout/ReviewItems/ReviewItems';
import {
  formatAboutYourPPMItem,
  formatExpenseItems,
  formatProGearItems,
  formatWeightTicketItems,
} from 'utils/ppmCloseout';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { formatCents, formatWeight } from 'utils/formatters';
import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';
import { deleteWeightTicket } from 'services/internalApi';
import ppmStyles from 'components/Customer/PPM/PPM.module.scss';
import { hasCompleteWeightTickets } from 'utils/shipments';

const ReviewDeleteCloseoutItemModal = ({ onClose, onSubmit, itemToDelete }) => (
  <div>
    <Overlay />
    <ModalContainer>
      <Modal>
        <ModalClose handleClick={() => onClose(false)} />
        <ModalTitle>
          <h3>Delete this?</h3>
        </ModalTitle>
        <p>Your information will be gone. You’ll need to start over if you want it back.</p>
        <ModalActions>
          <Button
            className="usa-button--destructive"
            type="submit"
            onClick={() => onSubmit(itemToDelete.itemType, itemToDelete.itemId, itemToDelete.itemETag)}
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

const Review = () => {
  const [isDeleteModalVisible, setIsDeleteModalVisible] = useState(false);
  const [itemToDelete, setItemToDelete] = useState();
  const { moveId, mtoShipmentId } = useParams();
  const mtoShipment = useSelector((state) => selectMTOShipmentById(state, mtoShipmentId));

  const weightTickets = mtoShipment?.ppmShipment?.weightTickets;
  const proGear = mtoShipment?.ppmShipment?.proGear;
  const expenses = mtoShipment?.ppmShipment?.expenses;

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

  const handleDelete = (itemType, itemId, itemETag) => {
    setItemToDelete(() => ({
      itemType,
      itemId,
      itemETag,
    }));
    setIsDeleteModalVisible(true);
  };

  const onDeleteSubmit = (itemType, itemId, itemETag) => {
    if (itemType === 'weightTicket') {
      deleteWeightTicket(itemId, itemETag)
        .then(() => {
          setIsDeleteModalVisible(false);
        })
        .catch(() => {});
    }
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

  const weightTicketsTotal = weightTickets?.reduce((prev, curr) => prev + (curr.fullWeight - curr.emptyWeight), 0);

  const canAdvance = hasCompleteWeightTickets(weightTickets);

  const proGearContents = formatProGearItems(
    proGear,
    customerRoutes.SHIPMENT_PPM_PRO_GEAR_EDIT_PATH,
    { moveId, mtoShipmentId },
    handleDelete,
  );

  const proGearTotal = proGear?.reduce((prev, curr) => {
    if (curr.constructedWeight) {
      return prev + curr.constructedWeight;
    }
    return prev + (curr.fullWeight - curr.emptyWeight);
  }, 0);

  const expenseContents = formatExpenseItems(
    expenses,
    customerRoutes.SHIPMENT_PPM_EXPENSES_EDIT_PATH,
    {
      moveId,
      mtoShipmentId,
    },
    handleDelete,
  );

  const expensesTotal = expenses?.reduce((prev, curr) => prev + curr.amount, 0);

  return (
    <div className={classnames(ppmPageStyles.ppmPageStyle, styles.PPMReview)}>
      <ScrollToTop />
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
            <ShipmentTag shipmentType={shipmentTypes.PPM} />
            <h1>Review</h1>
            <SectionWrapper className={styles.aboutSection}>
              <ReviewItems heading={<h2>About Your PPM</h2>} contents={aboutYourPPM} />
            </SectionWrapper>
            <SectionWrapper>
              <h2>Documents</h2>
              <ReviewItems
                className={styles.reviewItems}
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
                emptyMessage="No weight tickets uploaded. Add at least one set of weight tickets to request payment."
              />
              <ReviewItems
                className={styles.reviewItems}
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
                className={styles.reviewItems}
                heading={
                  <>
                    <h3>Expenses</h3>
                    <span>(${formatCents(expensesTotal)})</span>
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
              >
                Finish Later
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
