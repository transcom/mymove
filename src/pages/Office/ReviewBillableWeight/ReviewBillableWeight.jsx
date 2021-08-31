import React from 'react';
import { Button } from '@trussworks/react-uswds';
import { useParams } from 'react-router-dom';

import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import DocumentViewerSidebar from 'pages/Office/DocumentViewerSidebar/DocumentViewerSidebar';
import { useOrdersDocumentQueries } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

export default function ReviewBillableWeight() {
  const { moveCode } = useParams();

  const { upload, isLoading, isError } = useOrdersDocumentQueries(moveCode);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const documentsForViewer = Object.values(upload);

  return (
    <>
      <DocumentViewer files={documentsForViewer} />
      <DocumentViewerSidebar title="Review weights" subtitle="Edit max billableweight" onClose={() => {}}>
        <DocumentViewerSidebar.Content>Hello</DocumentViewerSidebar.Content>
        <DocumentViewerSidebar.Footer>
          <Button>Button</Button>
        </DocumentViewerSidebar.Footer>
      </DocumentViewerSidebar>
    </>
  );
}
