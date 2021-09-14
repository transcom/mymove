import React, { useState } from 'react';
// import { connect } from 'react-redux';
import { Alert, Button } from '@trussworks/react-uswds';
import { useHistory, useParams } from 'react-router-dom';
import { generatePath } from 'react-router';

import DocumentViewerSidebar from '../DocumentViewerSidebar/DocumentViewerSidebar';

import styles from 'styles/documentViewerWithSidebar.module.scss';
import { tioRoutes } from 'constants/routes';
import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import ShipmentCard from 'components/Office/BillableWeight/ShipmentCard/ShipmentCard';
import WeightSummary from 'components/Office/WeightSummary/WeightSummary';
import { useMoveTaskOrderQueries, useOrdersDocumentQueries } from 'hooks/queries';
import shipmentIsOverweight from 'utils/shipmentIsOverweight';

export default function ReviewBillableWeight() {
  const [selectedShipmentIndex, setSelectedShipmentIndex] = useState(0);

  const { moveCode } = useParams();
  const { move, mtoShipments } = useMoveTaskOrderQueries(moveCode);
  const handleClickNextButton = () => {
    const newSelectedShipmentIdx = selectedShipmentIndex + 1;
    setSelectedShipmentIndex(newSelectedShipmentIdx);
  };

  const handleClickBackButton = () => {
    const newSelectedShipmentIdx = selectedShipmentIndex - 1;
    setSelectedShipmentIndex(newSelectedShipmentIdx);
  };

  const isLastShipment = selectedShipmentIndex === mtoShipments?.length - 1;
  const history = useHistory();
  const [sidebarType, setSidebarType] = React.useState('MAX');

  const { upload, isLoading, isError } = useOrdersDocumentQueries(moveCode);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const documentsForViewer = Object.values(upload);

  const handleClose = () => {
    history.push(generatePath(tioRoutes.PAYMENT_REQUESTS_PATH, { moveCode }));
  };

  return (
    <div className={styles.DocumentWrapper}>
      <div className={styles.embed}>
        <DocumentViewer files={documentsForViewer} />
      </div>
      <div className={styles.sidebar}>
        {sidebarType === 'MAX' ? (
          <DocumentViewerSidebar title="Review weights" subtitle="Edit max billable weight" onClose={handleClose}>
            <DocumentViewerSidebar.Content>
              Review max billable weight content should go in here
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
            description={`Shipment ${selectedShipmentIndex + 1} of ${mtoShipments.length}`}
            onClose={() => {}}
          >
            <DocumentViewerSidebar.Content>
              <div style={{ width: '350px', backgroundColor: 'white', marginBottom: '16px' }}>
                <Alert slim type="error">
                  {`Max billable weight exceeded. \nPlease resolve.`}
                </Alert>
                {(!mtoShipments[0].reweighWeight || !mtoShipments[0].estimatedWeight) && (
                  <Alert slim type="warning">
                    Shipment missing information
                  </Alert>
                )}
                {shipmentIsOverweight(mtoShipments[0].estimatedWeight, mtoShipments[0].billableWeight) && (
                  <Alert slim type="warning">
                    Shipment exceeds 110% of estimated weight.
                  </Alert>
                )}
                <WeightSummary
                  maxBillableWeight={move.maxBillableWeight}
                  totalBillableWeight={move.totalBillableWeight}
                  weightRequested={move.weightRequested}
                  weightAllowance={move.weightAllowance}
                  totalBillableWeightFlag
                  shipments={mtoShipments}
                />
              </div>
              <div style={{ height: '100%', width: '350px' }}>
                <ShipmentCard
                  billableWeight={mtoShipments[0].billableWeight}
                  dateReweighRequested={mtoShipments[0].dateReweighRequested}
                  departedDate={mtoShipments[0].departedDate}
                  pickupAddress={mtoShipments[selectedShipmentIndex].pickupAddress}
                  destinationAddress={mtoShipments[0].destinationAddress}
                  estimatedWeight={mtoShipments[0].estimatedWeight}
                  originalWeight={mtoShipments[0].originalWeight}
                  reweighRemarks={mtoShipments[0].reweighRemarks}
                  reweighWeight={mtoShipments[0].reweighWeight}
                />
              </div>
            </DocumentViewerSidebar.Content>
            <DocumentViewerSidebar.Footer
              style={{ position: 'fixed', bottom: '0', width: '100%', backgroundColor: 'white' }}
            >
              <div style={{ display: 'flex' }}>
                <Button type="button" onClick={handleClickBackButton} secondary>
                  Back
                </Button>
                {!isLastShipment && (
                  <Button type="button" onClick={handleClickNextButton}>
                    Next Shipment
                  </Button>
                )}
              </div>
            </DocumentViewerSidebar.Footer>
          </DocumentViewerSidebar>
        )}
      </div>
    </div>
  );
}
