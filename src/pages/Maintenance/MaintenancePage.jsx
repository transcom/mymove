import React from 'react';

import styles from './MaintenancePage.module.scss';

const MaintenancePage = () => {
  return (
    <div className={styles.maintenance_wrapper}>
      <h1>System Maintenance</h1>
      <p>This system is currently undergoing maintenance. Please check back later.</p>
    </div>
  );
};

export default MaintenancePage;
