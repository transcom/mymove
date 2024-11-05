import React from 'react';

import styles from './MaintenancePage.module.scss';

import CUIHeader from 'components/CUIHeader/CUIHeader';
import MilMoveHeader from 'components/MilMoveHeader';

const MaintenancePage = () => {
  return (
    <>
      <CUIHeader />
      <MilMoveHeader />
      <div className={styles.maintenance_wrapper}>
        <h1>System Maintenance</h1>
        <p>This system is currently undergoing maintenance. Please check back later.</p>
      </div>
    </>
  );
};

export default MaintenancePage;
