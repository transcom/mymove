import React from 'react';
import { useParams } from 'react-router-dom';

import styles from './ServicesCounselingMoveDocumentWrapper.module.scss';

import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { useOrdersDocumentQueries } from 'hooks/queries';
import ServicesCounselingMoveAllowances from 'pages/Office/ServicesCounselingMoveAllowances/ServicesCounselingMoveAllowances';

const ServicesCounselingMoveDocumentWrapper = () => {
  const { moveCode } = useParams();

  const { upload, isLoading, isError } = useOrdersDocumentQueries(moveCode);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const documentsForViewer = Object.values(upload);

  return (
    <div className={styles.DocumentWrapper}>
      {documentsForViewer && (
        <div className={styles.embed}>
          <DocumentViewer files={documentsForViewer} />
        </div>
      )}
      <ServicesCounselingMoveAllowances moveCode={moveCode} />
    </div>
  );
};

export default ServicesCounselingMoveDocumentWrapper;
