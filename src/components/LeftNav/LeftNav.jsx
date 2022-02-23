import React, { useState, cloneElement } from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';

import styles from './LeftNav.module.scss';

import LeftNavSection from 'components/LeftNavSection/LeftNavSection';

const sectionLabels = {
  'requested-shipments': 'Requested shipments',
  'approved-shipments': 'Approved shipments',
  orders: 'Orders',
  allowances: 'Allowances',
  'customer-info': 'Customer info',
  'billable-weights': 'Billable weights',
  'payment-requests': 'Payment requests',
};

const LeftNav = ({ className, children, sections }) => {
  const arrayChildren = React.Children.toArray(children);
  const [activeSection, setActiveSection] = useState('');

  return (
    <nav className={classnames(styles.LeftNav, className)}>
      {sections.map((s) => {
        return (
          <LeftNavSection
            key={`sidenav_${s}`}
            sectionName={s}
            isActive={s === activeSection}
            onClickHandler={() => setActiveSection(s)}
          >
            {sectionLabels[`${s}`]}
            {React.Children.map(arrayChildren, (child) => {
              return cloneElement(child, { showTag: s === child.props.associatedSectionName && child.props.showTag });
            })}
          </LeftNavSection>
        );
      })}
    </nav>
  );
};

LeftNav.propTypes = {
  className: PropTypes.string,
  children: PropTypes.node,
  sections: PropTypes.arrayOf(PropTypes.string).isRequired,
};

LeftNav.defaultProps = {
  className: '',
  children: null,
};

export default LeftNav;
