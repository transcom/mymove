import React from 'react';

import styles from './MoveOrders.module.scss';

import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import samplePDF from 'components/DocumentViewer/sample.pdf';

const MoveOrders = () => {
  return (
    <div className={styles.MoveOrders}>
      <div className={styles.embed}>
        <DocumentViewer filename="ThisIsAVeryLongDocumentTitle.pdf" fileType="pdf" filePath={samplePDF} />
      </div>
      <div className={styles.sidebar}>View orders</div>
    </div>
  );
};

export default MoveOrders;
