import React from 'react';
import { Button, Alert } from '@trussworks/react-uswds';
import { useHistory, useParams } from 'react-router-dom';
import { generatePath } from 'react-router';

import DocumentViewerSidebar from '../DocumentViewerSidebar/DocumentViewerSidebar';

import reviewBillableWeightStyles from './ReviewBillableWeight.module.scss';

import styles from 'styles/documentViewerWithSidebar.module.scss';
import { tioRoutes } from 'constants/routes';
import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import ShipmentCard from 'components/Office/BillableWeight/ShipmentCard/ShipmentCard';
import WeightSummary from 'components/Office/WeightSummary/WeightSummary';
import EditBillableWeight from 'components/Office/BillableWeight/EditBillableWeight/EditBillableWeight';
import { useOrdersDocumentQueries, useMovePaymentRequestsQueries } from 'hooks/queries';
import {
  includedStatusesForCalculatingWeights,
  useCalculatedTotalBillableWeight,
  useCalculatedWeightRequested,
  useCalculatedEstimatedWeight,
} from 'hooks/custom';
import { shipmentIsOverweight } from 'utils/shipmentWeights';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

export default function ReviewBillableWeight() {
  const [selectedShipmentIndex, setSelectedShipmentIndex] = React.useState(0);
  const [sidebarType, setSidebarType] = React.useState('MAX');

  const { moveCode } = useParams();
  const handleClickNextButton = () => {
    const newSelectedShipmentIdx = selectedShipmentIndex + 1;
    setSelectedShipmentIndex(newSelectedShipmentIdx);
  };

  const handleClickBackButton = () => {
    const newSelectedShipmentIdx = selectedShipmentIndex - 1;
    if (newSelectedShipmentIdx >= 0) {
      setSelectedShipmentIndex(newSelectedShipmentIdx);
    } else {
      setSidebarType('MAX');
    }
  };

  const history = useHistory();

  const { upload, isLoading, isError } = useOrdersDocumentQueries(moveCode);
  const { order, mtoShipments } = useMovePaymentRequestsQueries(moveCode);
  /* Only show shipments in statuses of approved, diversion requested, or cancellation requested */
  const filteredShipments = mtoShipments.filter((shipment) => includedStatusesForCalculatingWeights(shipment.status));
  const isLastShipment = selectedShipmentIndex === filteredShipments?.length - 1;

  const totalBillableWeight = useCalculatedTotalBillableWeight(filteredShipments);
  const weightRequested = useCalculatedWeightRequested(filteredShipments);
  const totalEstimatedWeight = useCalculatedEstimatedWeight(filteredShipments);
  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const maxBillableWeight = order.entitlement.authorizedWeight;
  const weightAllowance = order.entitlement.totalWeight;

  const documentsForViewer = Object.values(upload);

  const handleClose = () => {
    history.push(generatePath(tioRoutes.PAYMENT_REQUESTS_PATH, { moveCode }));
  };

  const selectedShipment = filteredShipments[selectedShipmentIndex];

  return (
    <div className={styles.DocumentWrapper}>
      <div className={styles.embed}>
        <DocumentViewer files={documentsForViewer} />
      </div>
      <div className={styles.sidebar}>
        {sidebarType === 'MAX' ? (
          <DocumentViewerSidebar title="Review weights" subtitle="Edit max billable weight" onClose={handleClose}>
            <DocumentViewerSidebar.Content>
              {totalBillableWeight > maxBillableWeight && (
                <Alert slim type="error" data-testid="maxBillableWeightAlert">
                  {`Max billable weight exceeded. \nPlease resolve.`}
                </Alert>
              )}
              <div className={reviewBillableWeightStyles.weightSummary}>
                <WeightSummary
                  maxBillableWeight={maxBillableWeight}
                  totalBillableWeight={totalBillableWeight}
                  weightRequested={weightRequested}
                  weightAllowance={weightAllowance}
                  totalBillableWeightFlag={totalBillableWeight > maxBillableWeight}
                  shipments={filteredShipments}
                />
              </div>
              <EditBillableWeight
                title="Max billable weight"
                estimatedWeight={totalEstimatedWeight}
                maxBillableWeight={maxBillableWeight}
                weightAllowance={weightAllowance}
              />
            </DocumentViewerSidebar.Content>
            <DocumentViewerSidebar.Footer>
              <Button
                onClick={() => {
                  setSidebarType('SHIPMENT');
                }}
              >
                Review shipment weights
              </Button>
            </DocumentViewerSidebar.Footer>
          </DocumentViewerSidebar>
        ) : (
          <DocumentViewerSidebar
            title="Review weights"
            subtitle="Shipment weights"
            description={`Shipment ${selectedShipmentIndex + 1} of ${filteredShipments?.length}`}
            onClose={handleClose}
          >
            <DocumentViewerSidebar.Content>
              <div className={reviewBillableWeightStyles.contentContainer}>
                {totalBillableWeight > maxBillableWeight && (
                  <Alert slim type="error" data-testid="maxBillableWeightAlert">
                    {`Max billable weight exceeded. \nPlease resolve.`}
                  </Alert>
                )}
                {((!selectedShipment.reweigh?.weight && selectedShipment.reweigh?.requestedAt) ||
                  !selectedShipment.primeEstimatedWeight) && (
                  <Alert slim type="warning">
                    Shipment missing information
                  </Alert>
                )}
                {shipmentIsOverweight(selectedShipment.primeEstimatedWeight, selectedShipment.primeActualWeight) && (
                  <Alert slim type="warning">
                    Shipment exceeds 110% of estimated weight.
                  </Alert>
                )}
                <div className={reviewBillableWeightStyles.weightSummary}>
                  <WeightSummary
                    maxBillableWeight={maxBillableWeight}
                    totalBillableWeight={totalBillableWeight}
                    weightRequested={weightRequested}
                    weightAllowance={weightAllowance}
                    totalBillableWeightFlag={totalBillableWeight > maxBillableWeight}
                    shipments={filteredShipments}
                  />
                </div>
              </div>
              <div className={reviewBillableWeightStyles.contentContainer}>
                <ShipmentCard
                  billableWeight={selectedShipment.billableWeightCap}
                  dateReweighRequested={selectedShipment.reweigh?.requestedAt}
                  departedDate={selectedShipment.actualPickupDate}
                  pickupAddress={selectedShipment.pickupAddress}
                  destinationAddress={selectedShipment.destinationAddress}
                  estimatedWeight={selectedShipment.primeEstimatedWeight}
                  originalWeight={selectedShipment.primeActualWeight}
                  reweighRemarks={selectedShipment.reweigh?.verificationReason}
                  reweighWeight={selectedShipment.reweigh?.weight}
                />
              </div>
            </DocumentViewerSidebar.Content>
            <DocumentViewerSidebar.Footer className={reviewBillableWeightStyles.footer}>
              <div className={reviewBillableWeightStyles.flex}>
                <Button onClick={handleClickBackButton} secondary>
                  Back
                </Button>
                {!isLastShipment && <Button onClick={handleClickNextButton}>Next Shipment</Button>}
              </div>
            </DocumentViewerSidebar.Footer>
          </DocumentViewerSidebar>
        )}
      </div>
    </div>
  );
}
