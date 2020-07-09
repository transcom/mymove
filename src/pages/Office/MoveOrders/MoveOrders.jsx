import React from 'react';

import styles from './MoveOrders.module.scss';

import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import samplePDF from 'components/DocumentViewer/sample.pdf';
import samplePDF2 from 'components/DocumentViewer/sample2.pdf';
import samplePDF3 from 'components/DocumentViewer/sample3.pdf';

const MoveOrders = () => {
  const testFiles = [
    {
      filename: 'Test File.pdf',
      fileType: 'pdf',
      filePath: samplePDF,
    },
    {
      filename: 'Test File 2.pdf',
      fileType: 'pdf',
      filePath: samplePDF2,
    },
    {
      filename: 'Test File 3.pdf',
      fileType: 'pdf',
      filePath: samplePDF3,
    },
  ];

  return (
    <div className={styles.MoveOrders}>
      <div className={styles.embed}>
        <DocumentViewer files={testFiles} />
      </div>
      <div className={styles.sidebar}>View orders</div>
    </div>
  );
};

export default MoveOrders;
