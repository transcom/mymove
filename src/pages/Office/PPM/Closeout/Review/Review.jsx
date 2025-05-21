import React, { useState } from 'react';
import { GridContainer, Grid, Button } from '@trussworks/react-uswds';
import { Link, useParams, generatePath, useNavigate } from 'react-router-dom';
import classnames from 'classnames';
import { useQueryClient, useMutation } from '@tanstack/react-query';

import styles from './Review.module.scss';

import ppmStyles from 'components/Shared/PPM/PPM.module.scss';
import formStyles from 'styles/form.module.scss';
import Alert from 'shared/Alert';
import ppmPageStyles from 'pages/Office/PPM/PPM.module.scss';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { shipmentTypes } from 'constants/shipments';
import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import { servicesCounselingRoutes } from 'constants/routes';
import ReviewItems from 'components/Shared/PPM/Closeout/ReviewItems/ReviewItems';
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
  hasCompletedAllWeightTickets,
  hasCompletedAllExpenses,
  hasCompletedAllProGear,
  hasIncompleteWeightTicket,
} from 'utils/shipments';
import { usePPMShipmentAndDocsOnlyQueries } from 'hooks/queries';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { deleteMovingExpense, deleteWeightTicket, deleteProGearWeightTicket } from 'services/ghcApi';
import { DOCUMENTS } from 'constants/queryKeys';
import { PPM_TYPES } from 'shared/constants';

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

const Review = () => {
  const [isDeleteModalVisible, setIsDeleteModalVisible] = useState(false);
  const [itemToDelete, setItemToDelete] = useState();
  const [isDeleting, setIsDeleting] = useState();
  const [alert, setAlert] = useState(null);
  const { moveCode, shipmentId } = useParams();
  const { mtoShipment, documents, isLoading, isError } = usePPMShipmentAndDocsOnlyQueries(shipmentId);
  const ppmShipment = mtoShipment?.ppmShipment || {};
  const { ppmType } = ppmShipment;

  const weightTickets = documents?.WeightTickets ?? [];
  const proGear = documents?.ProGearWeightTickets ?? [];
  const expenses = documents?.MovingExpenses ?? [];
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const { mutate: mutateWeightTicket } = useMutation(deleteWeightTicket, {
    onSuccess: () => {
      setIsDeleteModalVisible(false);
      queryClient.invalidateQueries([DOCUMENTS, shipmentId]);
      setIsDeleting(false);
    },
  });

  const { mutate: mutateDeleteMovingExpense } = useMutation(deleteMovingExpense, {
    onSuccess: () => {
      setIsDeleteModalVisible(false);
      queryClient.invalidateQueries([DOCUMENTS, shipmentId]);
      setIsDeleting(false);
    },
  });

  const { mutate: mutateProGearWeightTicket } = useMutation(deleteProGearWeightTicket, {
    onSuccess: () => {
      setIsDeleteModalVisible(false);
      queryClient.invalidateQueries([DOCUMENTS, shipmentId]);
      setIsDeleting(false);
    },
  });

  const weightTicketCreatePath = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_WEIGHT_TICKETS_PATH, {
    moveCode,
    shipmentId,
  });
  const proGearCreatePath = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_PRO_GEAR_PATH, {
    moveCode,
    shipmentId,
  });
  const expensesCreatePath = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_EXPENSES_PATH, {
    moveCode,
    shipmentId,
  });

  if (isError) return <SomethingWentWrong />;

  if (!mtoShipment || isLoading) {
    return <LoadingPlaceholder />;
  }

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
    if (isDeleting) return;
    const ppmShipmentId = mtoShipment.ppmShipment?.id;

    if (itemType === 'weightTicket') {
      setIsDeleting(true);
      mutateWeightTicket(
        { ppmShipmentId, weightTicketId: itemId },
        {
          onSuccess: () => {
            setAlert({ type: 'success', message: `${itemNumber} successfully deleted.` });
          },
          onError: () => {
            setIsDeleting(false);
            setAlert({ type: 'error', message: `Something went wrong deleting ${itemNumber}. Please try again.` });
          },
        },
      );
    }
    if (itemType === 'proGear') {
      setIsDeleting(true);
      mutateProGearWeightTicket(
        { ppmShipmentId, proGearWeightTicketId: itemId },
        {
          onSuccess: () => {
            setAlert({ type: 'success', message: `${itemNumber} successfully deleted.` });
          },
          onError: (error) => {
            setIsDeleting(false);
            setAlert({
              type: 'error',
              message: `${error} Something went wrong deleting ${itemNumber}. Please try again.`,
            });
          },
        },
      );
    }
    if (itemType === 'expense') {
      setIsDeleting(true);
      mutateDeleteMovingExpense(
        { ppmShipmentId, movingExpenseId: itemId },
        {
          onSuccess: () => {
            setAlert({ type: 'success', message: `${itemNumber} successfully deleted.` });
          },
          onError: () => {
            setIsDeleting(false);
            setAlert({ type: 'error', message: `Something went wrong deleting ${itemNumber}. Please try again.` });
          },
        },
      );
    }
  };

  const aboutYourPPM = formatAboutYourPPMItem(
    mtoShipment?.ppmShipment,
    servicesCounselingRoutes.BASE_SHIPMENT_PPM_ABOUT_PATH,
    {
      moveCode,
      shipmentId,
    },
  );

  const weightTicketContents = formatWeightTicketItems(
    weightTickets,
    servicesCounselingRoutes.BASE_SHIPMENT_PPM_WEIGHT_TICKETS_EDIT_PATH,
    { moveCode, shipmentId },
    handleDelete,
  );

  const weightTicketsTotal = getTotalNetWeightForWeightTickets(weightTickets);

  const canAdvance =
    hasCompletedAllWeightTickets(weightTickets, ppmType) &&
    hasCompletedAllExpenses(expenses) &&
    hasCompletedAllProGear(proGear);

  // PPM-SPRs must have at least one moving expense to advance
  const ppmSmalLPackageCanAdvance = ppmType === PPM_TYPES.SMALL_PACKAGE && expenses && expenses.length < 1;

  const showIncompleteError =
    hasIncompleteWeightTicket(weightTickets) || !hasCompletedAllExpenses(expenses) || !hasCompletedAllProGear(proGear);

  const proGearContents = formatProGearItems(
    proGear,
    servicesCounselingRoutes.BASE_SHIPMENT_PPM_PRO_GEAR_EDIT_PATH,
    { moveCode, shipmentId },
    handleDelete,
  );

  const proGearTotal = calculateTotalNetWeightForProGearWeightTickets(proGear);

  const expenseContents = formatExpenseItems(
    expenses,
    servicesCounselingRoutes.BASE_SHIPMENT_PPM_EXPENSES_EDIT_PATH,
    {
      moveCode,
      shipmentId,
    },
    handleDelete,
  );

  const expensesTotal = calculateTotalMovingExpensesAmount(expenses);

  const handleBack = () => {
    const path = generatePath(servicesCounselingRoutes.BASE_MOVE_VIEW_PATH, { moveCode });
    navigate(path);
  };

  const handleSubmit = () => {
    const path = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_COMPLETE_PATH, { moveCode, shipmentId });
    navigate(path);
  };

  return (
    <div className={ppmPageStyles.tabContent}>
      <div className={classnames(ppmPageStyles.container, styles.PPMReview)}>
        <GridContainer className={ppmPageStyles.gridContainer}>
          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }}>
              <div className={ppmPageStyles.closeoutPageWrapper}>
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
                {showIncompleteError && (
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
                  <h2>{ppmType === PPM_TYPES.SMALL_PACKAGE ? 'Small Package Expenses' : 'Documents'}</h2>
                  {ppmType !== PPM_TYPES.SMALL_PACKAGE && (
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
                  )}
                  {ppmType !== PPM_TYPES.SMALL_PACKAGE && (
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
                  )}
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
                <div className={`${formStyles.formActions} ${ppmStyles.buttonGroup}`}>
                  <Button
                    data-testid="formBackButton"
                    className={ppmStyles.backButton}
                    type="button"
                    onClick={handleBack}
                    secondary
                    outline
                  >
                    Back
                  </Button>
                  <Button
                    data-testid="saveAndContinueButton"
                    className={ppmStyles.saveButton}
                    type="button"
                    onClick={handleSubmit}
                    disabled={!canAdvance || ppmSmalLPackageCanAdvance}
                  >
                    Save & Continue
                  </Button>
                </div>
              </div>
            </Grid>
          </Grid>
        </GridContainer>
      </div>
    </div>
  );
};

export default Review;
