import React, { useEffect, useMemo, useRef, useState } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { Alert, Button, Grid } from '@trussworks/react-uswds';
import { generatePath, useNavigate, useParams } from 'react-router-dom';
import classNames from 'classnames';
import { connect } from 'react-redux';

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
import DocumentViewerFileManager from 'components/DocumentViewerFileManager/DocumentViewerFileManager';
import { PPM_TYPES, PPM_DOCUMENT_TYPES } from 'shared/constants';
import { PPM_DOCUMENT_STATUS } from 'constants/ppms';
import { fetchPaymentPacketBlob } from 'services/ghcApi';
import CompletePPMCloseoutConfirmationModal from 'components/Office/PPM/CompletePPMCloseoutConfirmationModal/CompletePPMCloseoutConfirmationModal';
import { setShowLoadingSpinner as setShowLoadingSpinnerAction } from 'store/general/actions';

export const ReviewDocuments = ({ readOnly, setShowLoadingSpinner }) => {
  const { shipmentId, moveCode } = useParams();
  const { orders, mtoShipments } = useReviewShipmentWeightsQuery(moveCode);
  const { mtoShipment, documents, ppmActualWeight, isLoading, isError } = usePPMShipmentDocsQueries(shipmentId);

  const [serverError, setServerError] = useState(null);
  const [showOverview, setShowOverview] = useState(false);

  const [isFileUploading, setFileUploading] = useState(false);

  const order = Object.values(orders)?.[0];
  const [currentTotalWeight, setCurrentTotalWeight] = useState(0);
  const [currentAllowableWeight, setCurrentAllowableWeight] = useState(0);
  const [currentMtoShipments, setCurrentMtoShipments] = useState([]);

  const [documentSetIndex, setDocumentSetIndex] = useState(0);
  const [moveHasExcessWeight, setMoveHasExcessWeight] = useState(false);

  const [paymentPacketFile, setPaymentPacketFile] = useState(null);
  const [packetLoading, setPacketLoading] = useState(false);
  const [isConfirmModalVisible, setIsConfirmModalVisible] = useState(false);

  const [allWeightTicketsRejected, setAllWeightTicketsRejected] = useState(false);
  const [allMovingExpensesRejected, setAllMovingExpensesRejected] = useState(false);
  let documentSets = useMemo(() => [], []);
  const weightTickets = useMemo(() => documents?.WeightTickets ?? [], [documents?.WeightTickets]);
  const movingExpenses = useMemo(() => documents?.MovingExpenses ?? [], [documents?.MovingExpenses]);
  const proGearWeightTickets = documents?.ProGearWeightTickets ?? [];
  const updateTotalWeight = (newWeight) => {
    setCurrentTotalWeight(newWeight);
  };

  const ppmShipmentInfo = useMemo(() => {
    if (mtoShipment && mtoShipment.ppmShipment) {
      return {
        ...mtoShipment.ppmShipment,
        miles: mtoShipment.distance,
        actualWeight: ppmActualWeight?.actualWeight ?? currentTotalWeight,
      };
    }
    return {};
  }, [mtoShipment, ppmActualWeight, currentTotalWeight]);
  const isPPMSPR = ppmShipmentInfo?.ppmType === PPM_TYPES.SMALL_PACKAGE;

  useEffect(() => {
    if (weightTickets.length > 0) {
      const allRejected = weightTickets.every((ticket) => ticket.status === PPM_DOCUMENT_STATUS.REJECTED);
      setAllWeightTicketsRejected(allRejected);
    } else {
      setAllWeightTicketsRejected(false);
    }
  }, [weightTickets]);

  useEffect(() => {
    if (movingExpenses.length > 0) {
      const allRejected = movingExpenses.every((ticket) => ticket.status === PPM_DOCUMENT_STATUS.REJECTED);
      setAllMovingExpensesRejected(allRejected);
    } else {
      setAllMovingExpensesRejected(false);
    }
  }, [movingExpenses]);

  useEffect(() => {
    if (showOverview) {
      if (allWeightTicketsRejected && weightTickets.length > 0) {
        setServerError(
          'Cannot closeout PPM. All weight tickets have been rejected. At least one approved weight ticket is required.',
        );
      } else if (allMovingExpensesRejected && movingExpenses.length > 0 && isPPMSPR) {
        setServerError(
          'Cannot closeout PPM. All moving expenses have been rejected. At least one approved moving expense is required for a PPM-SPR.',
        );
      } else {
        setServerError(null);
      }
    }
  }, [showOverview, allWeightTicketsRejected, weightTickets, allMovingExpensesRejected, movingExpenses, isPPMSPR]);

  const chronologicalComparatorProperty = (input) => input.createdAt;
  const compareChronologically = (itemA, itemB) =>
    chronologicalComparatorProperty(itemA) < chronologicalComparatorProperty(itemB) ? -1 : 1;

  const constructWeightTicket = (weightTicket, tripNumber) => ({
    documentSetType: PPM_DOCUMENT_TYPES.WEIGHT_TICKET,
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
    documentSetType: PPM_DOCUMENT_TYPES.PROGEAR_WEIGHT_TICKET,
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
        documentSetType: PPM_DOCUMENT_TYPES.MOVING_EXPENSE,
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

  // when a payment packet is previewed in the doc viewer, we can use browser memory to store and view the file using JS "blobs"
  // this stores the file in the browser memory and then we can point to the blob URL when previewing the file
  // using React State, we can just load the PDF via the temp URL
  const handleDownloadPaymentPacket = async () => {
    try {
      setPacketLoading(true);
      setShowLoadingSpinner(true, null);
      const blob = await fetchPaymentPacketBlob(mtoShipment.ppmShipment.id);
      const fileUrl = window.URL.createObjectURL(blob);

      setPaymentPacketFile({
        id: `payment-packet-${shipmentId}`,
        filename: `payment-packet-${shipmentId}.pdf`,
        url: fileUrl,
        createdAt: new Date().toISOString(),
        rotation: 0,
        contentType: 'application/pdf',
      });
    } catch (error) {
      const msg = error instanceof Error ? error.message : String(error);
      setServerError(msg);
    } finally {
      setPacketLoading(false);
      setShowLoadingSpinner(false, null);
    }
  };

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

  const formatDocumentSetDisplay = documentSetIndex + 1;
  const currentTripNumber = currentDocumentSet?.tripNumber != null ? currentDocumentSet.tripNumber + 1 : 0;
  const currentDocumentCategoryIndex =
    currentDocumentSet?.categoryIndex != null ? currentDocumentSet.categoryIndex + 1 : 0;

  useEffect(() => {
    if (
      currentTotalWeight === 0 &&
      documentSets[documentSetIndex]?.documentSet.status !== PPM_DOCUMENT_STATUS.REJECTED
    ) {
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

  const getAllUploads = () => {
    return documentSets.reduce((acc, documentSet) => {
      return acc.concat(documentSet.uploads);
    }, []);
  };

  const uploads = showOverview ? getAllUploads() : currentDocumentSet?.uploads;

  const handleSubmitPPMShipmentModal = () => {
    onContinue();
  };

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  return (
    <div data-testid="ReviewDocuments test" className={styles.ReviewDocuments}>
      <CompletePPMCloseoutConfirmationModal
        isOpen={isConfirmModalVisible}
        onClose={setIsConfirmModalVisible}
        onSubmit={handleSubmitPPMShipmentModal}
      />
      <div className={styles.embed}>
        {paymentPacketFile ? (
          // View the payment packet preview, allowing full unmount of the uploads
          <DocumentViewer key="packet" files={[paymentPacketFile]} allowDownload isFileUploading={false} />
        ) : (
          // View the uploads
          <DocumentViewer key="docs" files={uploads} allowDownload isFileUploading={isFileUploading} />
        )}
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

          {allWeightTicketsRejected && weightTickets.length > 0 && (
            <Grid className={styles.alertContainer}>
              <Alert headingLevel="h4" slim type="error">
                <span>
                  All weight tickets have been rejected. At least one approved weight ticket is required to proceed.
                </span>
              </Alert>
            </Grid>
          )}
          <div className={classNames(styles.top, styles.noBottomBorder)}>
            {!readOnly && !showOverview && currentDocumentSet.documentSetType === PPM_DOCUMENT_TYPES.WEIGHT_TICKET && (
              <>
                <DocumentViewerFileManager
                  title="Full Weight Documents"
                  orderId={order.orderId}
                  documentId={currentDocumentSet.documentSet.emptyDocumentId}
                  files={currentDocumentSet.documentSet.emptyDocument.uploads}
                  documentType={PPM_DOCUMENT_TYPES.WEIGHT_TICKET}
                  onAddFile={() => {
                    setFileUploading(true);
                  }}
                  mtoShipment={mtoShipment}
                  useChevron
                />
                &nbsp;
                <DocumentViewerFileManager
                  title="Empty Weight Documents"
                  orderId={order.orderId}
                  documentId={currentDocumentSet.documentSet.fullDocumentId}
                  files={currentDocumentSet.documentSet.fullDocument.uploads}
                  documentType={PPM_DOCUMENT_TYPES.WEIGHT_TICKET}
                  onAddFile={() => {
                    setFileUploading(true);
                  }}
                  mtoShipment={mtoShipment}
                  useChevron
                />
              </>
            )}
            {!readOnly && !showOverview && currentDocumentSet.documentSetType === PPM_DOCUMENT_TYPES.MOVING_EXPENSE && (
              <DocumentViewerFileManager
                title="Moving Expense Documents"
                orderId={order.orderId}
                documentId={currentDocumentSet.documentSet.documentId}
                files={currentDocumentSet.documentSet.document.uploads}
                documentType={PPM_DOCUMENT_TYPES.MOVING_EXPENSE}
                onAddFile={() => {
                  setFileUploading(true);
                }}
                mtoShipment={mtoShipment}
                useChevron
              />
            )}
            {!readOnly &&
              !showOverview &&
              currentDocumentSet.documentSetType === PPM_DOCUMENT_TYPES.PROGEAR_WEIGHT_TICKET && (
                <DocumentViewerFileManager
                  title="Pro Gear Documents"
                  orderId={order.orderId}
                  documentId={currentDocumentSet.documentSet.documentId}
                  files={currentDocumentSet.documentSet.document.uploads}
                  documentType={PPM_DOCUMENT_TYPES.PROGEAR_WEIGHT_TICKET}
                  onAddFile={() => {
                    setFileUploading(true);
                  }}
                  mtoShipment={mtoShipment}
                  useChevron
                />
              )}
            {!readOnly &&
              showOverview &&
              documentSets.map((documentSet) => {
                if (documentSet.documentSetType === PPM_DOCUMENT_TYPES.WEIGHT_TICKET) {
                  return (
                    <>
                      <DocumentViewerFileManager
                        title="Full Weight Documents"
                        orderId={order.orderId}
                        documentId={documentSet.documentSet.emptyDocumentId}
                        files={documentSet.documentSet.emptyDocument.uploads}
                        documentType={PPM_DOCUMENT_TYPES.WEIGHT_TICKET}
                        onAddFile={() => {
                          setFileUploading(true);
                        }}
                        mtoShipment={mtoShipment}
                        useChevron
                      />
                      <DocumentViewerFileManager
                        title="Empty Weight Documents"
                        orderId={order.orderId}
                        documentId={documentSet.documentSet.fullDocumentId}
                        files={documentSet.documentSet.fullDocument.uploads}
                        documentType={PPM_DOCUMENT_TYPES.WEIGHT_TICKET}
                        onAddFile={() => {
                          setFileUploading(true);
                        }}
                        mtoShipment={mtoShipment}
                        useChevron
                      />
                    </>
                  );
                }
                if (documentSet.documentSetType === PPM_DOCUMENT_TYPES.PROGEAR_WEIGHT_TICKET) {
                  return (
                    <DocumentViewerFileManager
                      title="Pro Gear Documents"
                      orderId={order.orderId}
                      documentId={documentSet.documentSet.documentId}
                      files={documentSet.documentSet.document.uploads}
                      documentType={PPM_DOCUMENT_TYPES.PROGEAR_WEIGHT_TICKET}
                      onAddFile={() => {
                        setFileUploading(true);
                      }}
                      mtoShipment={mtoShipment}
                      useChevron
                    />
                  );
                }
                return (
                  <DocumentViewerFileManager
                    title="Moving Expense Documents"
                    orderId={order.orderId}
                    documentId={documentSet.documentSet.documentId}
                    files={documentSet.documentSet.document.uploads}
                    documentType={PPM_DOCUMENT_TYPES.MOVING_EXPENSE}
                    onAddFile={() => {
                      setFileUploading(true);
                    }}
                    mtoShipment={mtoShipment}
                    useChevron
                  />
                );
              })}
          </div>
          <br />
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
                {currentDocumentSet.documentSetType === PPM_DOCUMENT_TYPES.WEIGHT_TICKET && (
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
                {currentDocumentSet.documentSetType === PPM_DOCUMENT_TYPES.PROGEAR_WEIGHT_TICKET && (
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
                {currentDocumentSet.documentSetType === PPM_DOCUMENT_TYPES.MOVING_EXPENSE && (
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
          {!paymentPacketFile && (
            <Button className="usa-button--secondary" onClick={onBack} disabled={disableBackButton}>
              Back
            </Button>
          )}

          {showOverview && !paymentPacketFile && !readOnly && (
            <Button
              onClick={handleDownloadPaymentPacket}
              disabled={
                packetLoading ||
                (showOverview && allWeightTicketsRejected && weightTickets.length > 0) ||
                (showOverview && allMovingExpensesRejected && movingExpenses.length > 0 && isPPMSPR)
              }
            >
              Preview PPM Payment Packet
            </Button>
          )}

          {showOverview && paymentPacketFile && (
            <>
              <Button
                className="usa-button--secondary"
                onClick={() => {
                  // reset back to document review
                  setPaymentPacketFile(null);
                  setDocumentSetIndex(0);
                  setShowOverview(false);
                }}
              >
                Edit PPM
              </Button>
              <Button
                onClick={() => {
                  setIsConfirmModalVisible(true);
                }}
              >
                Complete PPM Review
              </Button>
            </>
          )}

          {!showOverview && (
            <Button
              type="submit"
              onClick={onContinue}
              data-testid="reviewDocumentsContinueButton"
              disabled={
                (showOverview && allWeightTicketsRejected && weightTickets.length > 0) ||
                (showOverview && allMovingExpensesRejected && movingExpenses.length > 0 && isPPMSPR)
              }
            >
              Continue
            </Button>
          )}

          {readOnly && showOverview && (
            <Button type="submit" onClick={onClose} data-testid="closeBtn">
              Close
            </Button>
          )}
        </DocumentViewerSidebar.Footer>
      </DocumentViewerSidebar>
    </div>
  );
};

ReviewDocuments.defaultProps = {
  loadingMessage: null,
};

const mapDispatchToProps = {
  setShowLoadingSpinner: setShowLoadingSpinnerAction,
};

export default connect(() => ({}), mapDispatchToProps)(ReviewDocuments);
