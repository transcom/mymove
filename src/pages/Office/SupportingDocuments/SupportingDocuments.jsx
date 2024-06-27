import React from 'react';

import DocumentViewer from 'components/DocumentViewer/DocumentViewer';

const SupportingDocuments = ({ uploads }) => {
  if (!uploads || uploads.constructor !== Array || uploads?.length <= 0) {
    return <h2>No supporting documents have been uploaded.</h2>;
  }
  return <DocumentViewer files={uploads} allowDownload />;
};

export default SupportingDocuments;
