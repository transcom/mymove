import React from 'react';
import { matchPath, useLocation, useParams } from 'react-router-dom';

import styles from './ServicesCounselingMoveDocumentWrapper.module.scss';

import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { useOrdersDocumentQueries } from 'hooks/queries';
import ServicesCounselingMoveAllowances from 'pages/Office/ServicesCounselingMoveAllowances/ServicesCounselingMoveAllowances';
import ServicesCounselingOrders from 'pages/Office/ServicesCounselingOrders/ServicesCounselingOrders';

const ServicesCounselingMoveDocumentWrapper = () => {
  const { moveCode } = useParams();
  const { pathname } = useLocation();

  const { upload, isLoading, isError } = useOrdersDocumentQueries(moveCode);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const documentsForViewer = Object.values(upload);

  const showOrders = matchPath(pathname, {
    path: '/counseling/moves/:moveCode/orders',
    exact: true,
  });

  return (
    <div className={styles.DocumentWrapper}>
      {documentsForViewer && (
        <div className={styles.embed}>
          <DocumentViewer files={documentsForViewer} />
        </div>
      )}
      {showOrders ? (
        <ServicesCounselingOrders moveCode={moveCode} />
      ) : (
        <ServicesCounselingMoveAllowances moveCode={moveCode} />
      )}
    </div>
  );
};

export default ServicesCounselingMoveDocumentWrapper;
