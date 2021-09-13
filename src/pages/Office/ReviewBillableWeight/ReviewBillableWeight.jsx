import React from 'react';
import { Button } from '@trussworks/react-uswds';
import { useHistory, useParams } from 'react-router-dom';
import { generatePath } from 'react-router';

import styles from 'styles/documentViewerWithSidebar.module.scss';
import { tioRoutes } from 'constants/routes';
import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import WeightSummary from 'components/Office/WeightSummary/WeightSummary';
import EditBillableWeight from 'components/Office/BillableWeight/EditBillableWeight/EditBillableWeight';
import DocumentViewerSidebar from 'pages/Office/DocumentViewerSidebar/DocumentViewerSidebar';
import { calcWeightRequested, calcTotalBillableWeight, calcTotalEstimatedWeight } from 'utils/shipmentWeights';
import { useOrdersDocumentQueries, useMovePaymentRequestsQueries } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

export default function ReviewBillableWeight() {
  const { moveCode } = useParams();
  const history = useHistory();
  const [sidebarType, setSidebarType] = React.useState('MAX');

  const { upload, isLoading, isError } = useOrdersDocumentQueries(moveCode);
  const { order, mtoShipments } = useMovePaymentRequestsQueries(moveCode);
  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  // weights
  const maxBillableWeight = order.entitlement.authorizedWeight;
  const weightAllowance = order.entitlement.totalWeight;
  const weightRequested = calcWeightRequested(mtoShipments);
  const totalBillableWeight = calcTotalBillableWeight(mtoShipments);
  const totalEstimatedWeight = calcTotalEstimatedWeight(mtoShipments);
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
              <WeightSummary
                maxBillableWeight={maxBillableWeight}
                totalBillableWeight={totalBillableWeight}
                weightRequested={weightRequested}
                weightAllowance={weightAllowance}
                shipments={mtoShipments}
              />
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
          <DocumentViewerSidebar title="Review weights" subtitle="Shipment weights" onClose={handleClose}>
            <DocumentViewerSidebar.Content>
              Review shipment weight content should go in here
            </DocumentViewerSidebar.Content>
          </DocumentViewerSidebar>
        )}
      </div>
    </div>
  );
}
