import React from 'react';
import { NavLink } from 'react-router-dom';
import { Tag } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';
import classnames from 'classnames';

import styles from './ServicesCounselingTabNav.module.scss';

import 'styles/office.scss';
import TabNav from 'components/TabNav';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

const ServicesCounselingTabNav = ({
  unapprovedShipmentCount = 0,
  shipmentWarnConcernCount = 0,
  missingOrdersInfoCount,
  moveCode,
}) => {
  const [supportingDocsFF, setSupportingDocsFF] = React.useState(false);
  React.useEffect(() => {
    const fetchData = async () => {
      setSupportingDocsFF(await isBooleanFlagEnabled('manage_supporting_docs'));
    };
    fetchData();
  }, []);

  let moveDetailsTagCount = 0;
  if (unapprovedShipmentCount > 0) {
    moveDetailsTagCount += unapprovedShipmentCount;
  }
  if (shipmentWarnConcernCount > 0) {
    moveDetailsTagCount += shipmentWarnConcernCount;
  }
  if (missingOrdersInfoCount > 0) {
    moveDetailsTagCount += missingOrdersInfoCount;
  }

  const items = [
    <NavLink
      end
      className={({ isActive }) => (isActive ? 'usa-current' : '')}
      to={`/counseling/moves/${moveCode}/details`}
      data-testid="MoveDetails-Tab"
    >
      <span className="tab-title">Move details</span>
      {moveDetailsTagCount > 0 && <Tag>{moveDetailsTagCount}</Tag>}
    </NavLink>,
    <NavLink
      end
      className={({ isActive }) => (isActive ? 'usa-current' : '')}
      to={`/counseling/moves/${moveCode}/mto`}
      data-testid="MoveTaskOrder-Tab"
    >
      <span className="tab-title">Move Task Order</span>
    </NavLink>,
    <NavLink
      end
      className={({ isActive }) => (isActive ? 'usa-current' : '')}
      to={`/counseling/moves/${moveCode}/customer-support-remarks`}
    >
      <span className="tab-title">Customer support remarks</span>
    </NavLink>,
    <NavLink
      end
      className={({ isActive }) => (isActive ? 'usa-current' : '')}
      to={`/counseling/moves/${moveCode}/history`}
      data-testid="MoveHistory-Tab"
    >
      <span className="tab-title">Move history</span>
    </NavLink>,
  ];

  if (supportingDocsFF)
    items.push(
      <NavLink
        end
        className={({ isActive }) => (isActive ? 'usa-current' : '')}
        to="supporting-documents"
        data-testid="SupportingDocuments-Tab"
      >
        <span className="tab-title">Supporting Documents</span>
      </NavLink>,
    );

  return (
    <header className="nav-header">
      <div
        className={
          supportingDocsFF ? classnames('grid-container-desktop-lg', styles.TabNav) : 'grid-container-desktop-lg'
        }
      >
        <TabNav items={items} />
      </div>
    </header>
  );
};

ServicesCounselingTabNav.defaultProps = {
  unapprovedShipmentCount: 0,
};

ServicesCounselingTabNav.propTypes = {
  unapprovedShipmentCount: PropTypes.number,
  moveCode: PropTypes.string.isRequired,
};

export default ServicesCounselingTabNav;
