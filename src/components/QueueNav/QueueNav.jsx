import React from 'react';
import { NavLink } from 'react-router-dom';

import styles from './QueueNav.module.scss';

import { tioRoutes, generalRoutes, servicesCounselingRoutes } from 'constants/routes';
import TabNav from 'components/TabNav';

const QueueNav = () => {
  return (
    <TabNav
      className={styles.tableTabs}
      items={[
        <NavLink
          end
          className={({ isActive }) => (isActive ? 'usa-current' : '')}
          to={servicesCounselingRoutes.BASE_QUEUE_COUNSELING_PATH}
        >
          <span data-testid="counseling-tab-link" className="tab-title">
            Counseling Queue
          </span>
        </NavLink>,
        <NavLink
          end
          className={({ isActive }) => (isActive ? 'usa-current' : '')}
          to={servicesCounselingRoutes.BASE_QUEUE_CLOSEOUT_PATH}
        >
          <span data-testid="closeout-tab-link" className="tab-title">
            PPM Closeout Queue
          </span>
        </NavLink>,
        <NavLink
          end
          className={({ isActive }) => (isActive ? 'usa-current' : '')}
          to={tioRoutes.BASE_QUEUE_COUNSELING_PATH}
        >
          <span data-testid="counseling-tab-link" className="tab-title">
            Payment Request Queue
          </span>
        </NavLink>,
        <NavLink
          end
          className={({ isActive }) => (isActive ? 'usa-current' : '')}
          to={generalRoutes.BASE_QUEUE_SEARCH_PATH}
        >
          <span data-testid="search-tab-link" className="tab-title">
            Search
          </span>
        </NavLink>,
      ]}
    />
  );
};

export default QueueNav;
