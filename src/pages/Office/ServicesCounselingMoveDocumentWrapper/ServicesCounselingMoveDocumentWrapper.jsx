import React from 'react';
import { matchPath, useLocation, useParams } from 'react-router-dom';

import styles from 'styles/documentViewerWithSidebar.module.scss';
import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import { servicesCounselingRoutes } from 'constants/routes';
import { useOrdersDocumentQueries } from 'hooks/queries';
import ServicesCounselingMoveAllowances from 'pages/Office/ServicesCounselingMoveAllowances/ServicesCounselingMoveAllowances';
import ServicesCounselingOrders from 'pages/Office/ServicesCounselingOrders/ServicesCounselingOrders';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const ServicesCounselingMoveDocumentWrapper = () => {
  const { moveCode } = useParams();
  const { pathname } = useLocation();
  const { upload, amendedUpload, isLoading, isError } = useOrdersDocumentQueries(moveCode);
  const documentsForViewer = Object.values(upload || {}).concat(Object.values(amendedUpload || {}));
  const hasDocuments = documentsForViewer?.length > 0;

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const showOrders = matchPath(
    {
      path: servicesCounselingRoutes.BASE_ORDERS_EDIT_PATH,
      exact: true,
    },
    pathname,
  );

  return (
    <div className={styles.DocumentWrapper}>
      {documentsForViewer && (
        <div className={styles.embed}>
          <DocumentViewer files={documentsForViewer} allowDownload />
        </div>
      )}
      {showOrders ? (
        <ServicesCounselingOrders moveCode={moveCode} hasDocuments={hasDocuments} />
      ) : (
        <ServicesCounselingMoveAllowances moveCode={moveCode} />
      )}
    </div>
  );
};

export default ServicesCounselingMoveDocumentWrapper;
