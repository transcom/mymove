import React from 'react';

import styles from './MoveOrders.module.scss';

import DocumentViewer from 'components/DocumentViewer/DocumentViewer';

const MoveOrders = () => {
  return (
    <div className={styles.MoveOrders}>
      <div className={styles.embed}>
        <DocumentViewer
          filename="ThisIsAVeryLongDocumentTitle.pdf"
          fileType="pdf"
          filePath="http://officelocal:3000/storage/user/9bfa91d2-7a0c-4de0-ae02-b8cf8b4b858b/uploads/7193a2b7-d260-40eb-a401-7c179359b03d?contentType=application%2Fpdf"
        />
      </div>
      <div className={styles.sidebar}>View orders</div>
    </div>
  );
};

export default MoveOrders;
