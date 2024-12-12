import React, { useEffect, useMemo, useRef, useState } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { Alert, Button, Grid } from '@trussworks/react-uswds';
import { generatePath, useNavigate, useParams } from 'react-router-dom';

import styles from './ReviewDocuments.module.scss';

import ReviewDocumentsSidePanel from 'components/Office/PPM/ReviewDocumentsSidePanel/ReviewDocumentsSidePanel';
import { ErrorMessage } from 'components/form';
import { servicesCounselingRoutes, tooRoutes } from 'constants/routes';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import DocumentViewerSidebar from 'pages/Office/DocumentViewerSidebar/DocumentViewerSidebar';
import { useReviewShipmentWeightsQuery, usePPMShipmentDocsQueries } from 'hooks/queries';
import ReviewWeightTicket from 'components/Office/PPM/ReviewWeightTicket/ReviewWeightTicket';
import ReviewExpense from 'components/Office/PPM/ReviewExpense/ReviewExpense';
import { DOCUMENTS } from 'constants/queryKeys';
import ReviewProGear from 'components/Office/PPM/ReviewProGear/ReviewProGear';
import { roleTypes } from 'constants/userRoles';
import { calculateWeightRequested } from 'hooks/custom';

// TODO: This should be in src/constants/ppms.js, but it's causing a lot of errors in unrelated tests, so I'll leave
//  this here for now.
const DOCUMENT_TYPES = {
  WEIGHT_TICKET: 'WEIGHT_TICKET',
  PROGEAR_WEIGHT_TICKET: 'PROGEAR_WEIGHT_TICKET',
  MOVING_EXPENSE: 'MOVING_EXPENSE',
};

export const ReviewDocuments = ({ readOnly }) => {
  const { shipmentId, moveCode } = useParams();
  const { orders, mtoShipments } = useReviewShipmentWeightsQuery(moveCode);
  const { mtoShipment, documents, ppmActualWeight, isLoading, isError } = usePPMShipmentDocsQueries(shipmentId);

  const order = Object.values(orders)?.[0];
  const [currentTotalWeight, setCurrentTotalWeight] = useState(0);
  const [currentAllowableWeight, setCurrentAllowableWeight] = useState(0);
  const [currentMtoShipments, setCurrentMtoShipments] = useState([]);

  const [documentSetIndex, setDocumentSetIndex] = useState(0);
  const [moveHasExcessWeight, setMoveHasExcessWeight] = useState(false);

  const [ppmShipmentInfo, setPpmShipmentInfo] = useState({});

  let documentSets = useMemo(() => [], []);
  const weightTickets = documents?.WeightTickets ?? [];
  const proGearWeightTickets = documents?.ProGearWeightTickets ?? [];
  const movingExpenses = documents?.MovingExpenses ?? [];
  const updateTotalWeight = (newWeight) => {
    setCurrentTotalWeight(newWeight);
  };
  useEffect(() => {
    if (currentTotalWeight === 0 && documentSets[documentSetIndex]?.documentSet.status !== 'REJECTED') {
      updateTotalWeight(ppmActualWeight?.actualWeight || 0);
    }
  }, [currentMtoShipments, ppmActualWeight?.actualWeight, currentTotalWeight, documentSets, documentSetIndex]);
  useEffect(() => {
    const totalMoveWeight = calculateWeightRequested(currentMtoShipments);
    setMoveHasExcessWeight(totalMoveWeight > order.entitlement.totalWeight);
  }, [currentMtoShipments, order.entitlement.totalWeight, currentTotalWeight]);
  useEffect(() => {
    setCurrentAllowableWeight(currentAllowableWeight);
  }, [currentAllowableWeight]);
  useEffect(() => {
    setCurrentMtoShipments(mtoShipments);
  }, [mtoShipments]);

  useEffect(() => {
    if (mtoShipment) {
      const updatedPpmShipmentInfo = {
        ...mtoShipment.ppmShipment,
        miles: mtoShipment.distance,
        actualWeight: currentTotalWeight,
      };
      setPpmShipmentInfo(updatedPpmShipmentInfo);
    }
  }, [mtoShipment, currentTotalWeight]);

  const chronologicalComparatorProperty = (input) => input.createdAt;
  const compareChronologically = (itemA, itemB) =>
    chronologicalComparatorProperty(itemA) < chronologicalComparatorProperty(itemB) ? -1 : 1;

  const constructWeightTicket = (weightTicket, tripNumber) => ({
    documentSetType: DOCUMENT_TYPES.WEIGHT_TICKET,
    documentSet: weightTicket,
    uploads: [
      ...weightTicket.emptyDocument.uploads,
      ...weightTicket.fullDocument.uploads,
      ...weightTicket.proofOfTrailerOwnershipDocument.uploads,
    ],
    tripNumber,
  });

  if (weightTickets.length > 0) {
    weightTickets.sort(compareChronologically);

    documentSets = documentSets.concat(weightTickets.map(constructWeightTicket));
  }

  const constructProGearWeightTicket = (weightTicket, tripNumber) => ({
    documentSetType: DOCUMENT_TYPES.PROGEAR_WEIGHT_TICKET,
    documentSet: weightTicket,
    uploads: weightTicket.document.uploads,
    tripNumber,
  });

  if (proGearWeightTickets.length > 0) {
    proGearWeightTickets.sort(compareChronologically);

    documentSets = documentSets.concat(proGearWeightTickets.map(constructProGearWeightTicket));
  }

  if (movingExpenses.length > 0) {
    // index individual input set elements by categorical type and chronological index.
    const accumulateMovingExpensesCategoricallyIndexed = (input) => {
      const constructExpenseCategoricallyIndexed = (movingExpense, categoryIndex) => ({
        documentSetType: DOCUMENT_TYPES.MOVING_EXPENSE,
        documentSet: movingExpense,
        uploads: movingExpense.document.uploads,
        categoryIndex,
      });

      const addFlattenedIndexToExpense = (expenseView, index) => ({ ...expenseView, tripNumber: index });
      // safari's dev team hasn't caught up to the chromium javascript ecma version, so there is no cross-browser availability for Object.groupBy
      const groupByFix = (iterable, key) => {
        const groupByResult = iterable.reduce((accumulator, item) => {
          (accumulator[key(item)] ??= []).push(item);
          return accumulator;
        }, {});
        return groupByResult;
      };
      const groupResult = groupByFix(input, ({ movingExpenseType }) => movingExpenseType);
      const assignDiscreetIndexesPerGroupElements = Object.values(groupResult).map((grp) =>
        grp.map(constructExpenseCategoricallyIndexed),
      );
      const flattenedGroupsWithUnifiedIndex = assignDiscreetIndexesPerGroupElements
        .flat()
        // even though the initial set was ordered, we have to adjust the order again. (Maintaining the index of chronological existence)
        .sort((itemA, itemB) => compareChronologically(itemA.documentSet, itemB.documentSet))
        .map(addFlattenedIndexToExpense);
      return flattenedGroupsWithUnifiedIndex;
    };

    // sort expenses by occurrence
    const sortedExpenses = [...movingExpenses].sort(compareChronologically);
    const resultSet = accumulateMovingExpensesCategoricallyIndexed(sortedExpenses);

    documentSets = documentSets.concat(resultSet);
  }

  const navigate = useNavigate();

  const formRef = useRef();
  const mainRef = useRef();

  const [serverError, setServerError] = useState(null);
  const [showOverview, setShowOverview] = useState(false);

  const queryClient = useQueryClient();

  const onClose = () => {
    navigate(
      generatePath(
        roleTypes.SERVICES_COUNSELOR ? servicesCounselingRoutes.BASE_MOVE_VIEW_PATH : tooRoutes.BASE_MOVE_VIEW_PATH,
        { moveCode },
      ),
    );
  };

  const onBack = () => {
    setServerError(null);
    if (showOverview) {
      setShowOverview(false);
    } else if (documentSetIndex > 0) {
      setDocumentSetIndex(documentSetIndex - 1);
    }
  };

  const onSuccess = () => {
    queryClient.invalidateQueries([DOCUMENTS, shipmentId]);

    if (documentSetIndex < documentSets.length - 1) {
      const newDocumentSetIndex = documentSetIndex + 1;
      setDocumentSetIndex(newDocumentSetIndex);
    } else {
      setShowOverview(true);
    }
  };

  const getAllUploads = () => {
    return documentSets.reduce((acc, documentSet) => {
      return acc.concat(documentSet.uploads);
    }, []);
  };

  const onError = () => {
    setServerError('There was an error submitting the form. Please try again later.');
  };

  const onErrorMessage = (errorMessage) => {
    setServerError(errorMessage);
  };

  const onConfirmSuccess = () => {
    if (roleTypes.SERVICES_COUNSELOR)
      navigate(generatePath(servicesCounselingRoutes.BASE_MOVE_VIEW_PATH, { moveCode }));
    else if (roleTypes.TOO) navigate(generatePath(tooRoutes.BASE_MOVE_VIEW_PATH, { moveCode }));
  };

  const onContinue = () => {
    setServerError(null);
    if (formRef.current) {
      formRef.current.handleSubmit();
    }
  };

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const currentDocumentSet = documentSets[documentSetIndex];
  const updateDocumentSetAllowableWeight = (newWeight) => {
    currentDocumentSet.documentSet.allowableWeight = newWeight;
  };
  const disableBackButton = documentSetIndex === 0 && !showOverview;

  const reviewShipmentWeightsURL = generatePath(servicesCounselingRoutes.BASE_REVIEW_SHIPMENT_WEIGHTS_PATH, {
    moveCode,
    shipmentId,
  });

  const reviewShipmentWeightsLink = <a href={reviewShipmentWeightsURL}>Review shipment weights</a>;

  const currentTripNumber = currentDocumentSet.tripNumber + 1;
  const currentDocumentCategoryIndex = currentDocumentSet.categoryIndex + 1;

  const formatDocumentSetDisplay = documentSetIndex + 1;

  let nextButton = 'Continue';
  if (showOverview) {
    nextButton = readOnly ? 'Close' : 'Confirm';
  }

  return (
    <div data-testid="ReviewDocuments test" className={styles.ReviewDocuments}>
      <div className={styles.embed}>
        <DocumentViewer files={showOverview ? getAllUploads() : currentDocumentSet.uploads} allowDownload />
      </div>
      <DocumentViewerSidebar
        title={readOnly ? 'View documents' : 'Review documents'}
        onClose={onClose}
        className={styles.sidebar}
        supertitle={
          showOverview ? 'All Document Sets' : `${formatDocumentSetDisplay} of ${documentSets.length} Document Sets`
        }
        defaultH3
        hyperlink={readOnly ? '' : reviewShipmentWeightsLink}
        readOnly={readOnly}
      >
        <DocumentViewerSidebar.Content mainRef={mainRef}>
          <NotificationScrollToTop dependency={documentSetIndex || serverError} target={mainRef.current} />
          {moveHasExcessWeight && (
            <Grid className={styles.alertContainer}>
              <Alert headingLevel="h4" slim type="warning">
                <span>This move has excess weight. Edit the PPM net weight to resolve.</span>
              </Alert>
            </Grid>
          )}
          <ErrorMessage className={styles.errorMessage} display={!!serverError}>
            {serverError}
          </ErrorMessage>
          {documentSets &&
            (showOverview ? (
              <ReviewDocumentsSidePanel
                ppmShipment={mtoShipment.ppmShipment}
                ppmShipmentInfo={ppmShipmentInfo}
                weightTickets={weightTickets}
                proGearTickets={proGearWeightTickets}
                expenseTickets={movingExpenses}
                onError={onError}
                onSuccess={onConfirmSuccess}
                formRef={formRef}
                allowableWeight={currentAllowableWeight}
                readOnly={readOnly}
                order={order}
              />
            ) : (
              <>
                {currentDocumentSet.documentSetType === DOCUMENT_TYPES.WEIGHT_TICKET && (
                  <ReviewWeightTicket
                    key={documentSetIndex}
                    weightTicket={currentDocumentSet.documentSet}
                    ppmShipmentInfo={ppmShipmentInfo}
                    ppmNumber="1"
                    tripNumber={currentTripNumber}
                    mtoShipment={mtoShipment}
                    order={order}
                    currentMtoShipments={currentMtoShipments}
                    setCurrentMtoShipments={setCurrentMtoShipments}
                    onError={onError}
                    onSuccess={onSuccess}
                    formRef={formRef}
                    allowableWeight={currentAllowableWeight}
                    updateTotalWeight={updateTotalWeight}
                    updateDocumentSetAllowableWeight={updateDocumentSetAllowableWeight}
                    readOnly={readOnly}
                  />
                )}
                {currentDocumentSet.documentSetType === DOCUMENT_TYPES.PROGEAR_WEIGHT_TICKET && (
                  <ReviewProGear
                    proGear={currentDocumentSet.documentSet}
                    ppmShipmentInfo={ppmShipmentInfo}
                    ppmNumber="1"
                    tripNumber={currentTripNumber}
                    mtoShipment={mtoShipment}
                    onError={onError}
                    onSuccess={onSuccess}
                    formRef={formRef}
                    readOnly={readOnly}
                    order={order}
                  />
                )}
                {currentDocumentSet.documentSetType === DOCUMENT_TYPES.MOVING_EXPENSE && (
                  <ReviewExpense
                    key={documentSetIndex}
                    expense={currentDocumentSet.documentSet}
                    ppmShipmentInfo={ppmShipmentInfo}
                    documentSets={documentSets}
                    documentSetIndex={documentSetIndex}
                    categoryIndex={currentDocumentCategoryIndex}
                    ppmNumber="1"
                    tripNumber={currentTripNumber}
                    mtoShipment={mtoShipment}
                    onError={onErrorMessage}
                    onSuccess={onSuccess}
                    formRef={formRef}
                    readOnly={readOnly}
                    order={order}
                  />
                )}
              </>
            ))}
        </DocumentViewerSidebar.Content>
        <DocumentViewerSidebar.Footer>
          <Button className="usa-button--secondary" onClick={onBack} disabled={disableBackButton}>
            Back
          </Button>
          <Button type="submit" onClick={onContinue} data-testid="reviewDocumentsContinueButton">
            {nextButton}
          </Button>
        </DocumentViewerSidebar.Footer>
      </DocumentViewerSidebar>
    </div>
  );
};

export default ReviewDocuments;
